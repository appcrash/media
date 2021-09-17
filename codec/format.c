#include <stdio.h>
#include <stdlib.h>
#include <inttypes.h>
#include <stdint.h>
#include <string.h>
#include <sys/stat.h>
#include <unistd.h>

#include "codec.h"


struct Payload* read_media_file(const char* file_path)
{
    AVFormatContext *ctx = NULL;
    int ret;
    AVPacket pkt;
    struct Payload *payload = NULL;
    struct stat fstat;

    if (!file_path) {
        PERR("read media file with file_path null");
        goto error;
    }
    if (stat(file_path,&fstat) < 0) {
        PERR("read media file %s with stat error",file_path);
        goto error;
    }
    if (fstat.st_size == 0) {
        PERR("read media file %s with size 0",file_path);
        goto error;
    }
    // ensure payload buffer is enough to hold all data
    // as it will be freed in go soon, wasting some bytes is ok
    payload = malloc(sizeof(struct Payload));
    payload->data = (char*)malloc(fstat.st_size);
    payload->size = 0;

    ret = avformat_open_input(&ctx, file_path, NULL, NULL);
    if (ret < 0) {
        PERR("read media file error");
        goto error;
    }
    ret = avformat_find_stream_info(ctx, NULL);
    if (ret < 0) {
        PERR("find stream info failed");
        goto error;
    }
    payload->bitrate = ctx->bit_rate;

    printf("nbstream is %d\nbitrate is %ld\npacket size is %d\nduration is %ld\n",
           ctx->nb_streams,ctx->bit_rate,ctx->packet_size,ctx->duration);
    printf("codec_id is %d\n",ctx->streams[0]->codecpar->codec_id);
    while(1) {
        ret = av_read_frame(ctx,&pkt);
        if (ret == AVERROR_EOF) {
            printf("eof is met, return\n");
            break;
        }
        //printf("pkt size is %d\n",pkt.size);
        memcpy(&payload->data[payload->size],pkt.data,pkt.size);
        payload->size += pkt.size;
    }

    avformat_close_input(&ctx);
    return payload;

error:
    if (ctx) {
        avformat_close_input(&ctx);
    }
    if (payload) {
        if (payload->data) {
            free(payload->data);
        }
        free(payload);
    }
    return NULL;
}

// support one channel only
int write_media_file(char *payload,int length,const char *file_path,int codec_id,int duration)
{
    AVFormatContext *ctx = NULL;
    AVStream *ostream = NULL;
    AVPacket pkt;
    int ret;

    av_init_packet(&pkt);
    ret = avformat_alloc_output_context2(&ctx,NULL,NULL,file_path);
    if (ret < 0) {
        PERR("avformat_alloc_output_context2 failed");
        goto error;
    }
    ostream = avformat_new_stream(ctx, NULL);
    if (!ostream) {
        PERR("avformat_new_stream failed");
        goto error;
    }

    AVCodecParameters *cp = ostream->codecpar;
    cp->channels = 1;
    cp->sample_rate = 8000;
    cp->codec_id = codec_id;
    cp->codec_type = AVMEDIA_TYPE_AUDIO;

    printf("file path is %s\n",file_path);
    printf("format is %s\n",ctx->oformat->name);
    printf("length is %d, duration is %d\n",length,duration);

    if (!(ctx->oformat->flags & AVFMT_NOFILE)) {
        printf("oformat flags is %x\n",ctx->oformat->flags);
        ret = avio_open(&ctx->pb, file_path, AVIO_FLAG_WRITE);
        if (ret < 0) {
            PERR("Could not open output file '%s'", file_path);
            goto error;
        }
    }


    ret = avformat_write_header(ctx, NULL);
    if (ret < 0) {
        PERR("avformat_write_header failed");
        goto error;
    }

    pkt.data = (uint8_t*)payload;
    pkt.size = length;
    pkt.stream_index = 0;
    pkt.duration = duration;

    ret = av_write_frame(ctx, &pkt);
    if (ret < 0) {
        PERR("av_write_frame failed");
        goto error;
    }
    av_write_trailer(ctx);
    avformat_free_context(ctx);

    return 0;

error:
    if (ctx) {
        avformat_free_context(ctx);
    }
    return -1;
}

struct RecordContext *record_init_context(const char *file_path,const char *params)
{
    struct RecordContext *record_ctx = NULL;
    AVFormatContext *ctx = NULL;
    AVDictionary *dict = NULL;
    AVDictionaryEntry *t = NULL;
    AVStream *ostream = NULL;
    int ret;

    ret = avformat_alloc_output_context2(&ctx,NULL, NULL, file_path);
    if (ret < 0) {
        PERR("avformat_alloc_output_context2 failed");
        goto cleanup;
    }
    ostream = avformat_new_stream(ctx, NULL);
    if (!ostream) {
        PERR("avformat_new_stream failed");
        goto cleanup;
    }
    AVCodecParameters *cp = ostream->codecpar;
    cp->codec_type = AVMEDIA_TYPE_AUDIO;

    /* AVCodecParameters does not support av_opt_* ... */
    if (av_dict_parse_string(&dict, params, "=", ",", 0) < 0) {
        goto cleanup;
    }
    while ((t = av_dict_get(dict, "", t, AV_DICT_IGNORE_SUFFIX))) {
        if (!av_strncasecmp(t->key, "channels",8)) {
            cp->channels = atoi(t->value);
        } else if (!av_strncasecmp(t->key, "sample_rate",11)) {
            cp->sample_rate = atoi(t->value);
        } else if (!av_strncasecmp(t->key,"codec_id",8)) {
            cp->codec_id = atoi(t->value);
        }
    }

    printf("record_ctx: channels=%d\n",cp->channels);
    printf("record_ctx: sample_rate=%d\n",cp->sample_rate);
    printf("record_ctx: codec_id=%d\n",cp->codec_id);
    if (!(ctx->oformat->flags & AVFMT_NOFILE)) {
        //printf("record_ctx:oformat flags is %x\n",ctx->oformat->flags);
        ret = avio_open(&ctx->pb, file_path, AVIO_FLAG_WRITE);
        if (ret < 0) {
            PERR("Could not open output file '%s'", file_path);
            goto cleanup;
        }
    }
    if (avformat_write_header(ctx,NULL) < 0) {
        PERR("record_ctx: avformat_write_header failed");
    }

    record_ctx = av_malloc(sizeof(struct RecordContext));
    record_ctx->ctx = ctx;
    goto done;
cleanup:
    if (ctx) {
        avformat_free_context(ctx);
    }
done:
    av_dict_free(&dict);
    return record_ctx;
}

void record_iterate(struct RecordContext *ctx,const char *buff,int32_t frame_delimits[],int nb_frame)
{
    int i,frame_start,frame_len;
    AVPacket pkt;

    if (nb_frame <= 0) {
        return;
    }
    av_init_packet(&pkt);
    pkt.stream_index = 0;
    frame_start = 0;
    for (i = 0; i < nb_frame; i++) {
        frame_len = frame_delimits[i] - frame_start;
        av_init_packet(&pkt);
        printf("[%d:%d  %d] ",frame_start,frame_len,frame_delimits[i]);
        pkt.data = (uint8_t*)&buff[frame_start];
        pkt.size = frame_len;
        frame_start = frame_delimits[i];
        if (av_write_frame(ctx->ctx, &pkt) < 0) {
            PERR("av_write_frame failed");
            goto error;
        }
    }

error:
    return;
}

void record_free(struct RecordContext *ctx)
{
    if (ctx) {
        if (ctx->ctx) {
            av_write_trailer(ctx->ctx);
            avformat_free_context(ctx->ctx);
        }
        av_free(ctx);
    }

}
