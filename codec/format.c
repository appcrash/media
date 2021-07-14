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
    printf("av_write_frame done...\n");

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
