#include <stdio.h>
#include <stdlib.h>
#include <inttypes.h>
#include <stdint.h>
#include <string.h>
#include <sys/stat.h>
#include <unistd.h>

#include "codec.h"

struct TranscodeContext *transcode_init_context(const char *from_codec_name,int from_sample_rate,int to_codec_id)
{
    AVCodecContext *codec_ctx = NULL;
    struct TranscodeContext *trans_ctx = NULL;
    AVPacket *packet = NULL;
    AVFrame *frame = NULL;
    AVCodec *fcodec = avcodec_find_decoder_by_name(from_codec_name);

    if (!fcodec) {
        PERR("decoder not available %s",from_codec_name);
        goto error;
    }

    codec_ctx = avcodec_alloc_context3(fcodec);
    //codec_ctx->sample_rate = from_sample_rate;
    codec_ctx->channels = 1;
    if (avcodec_open2(codec_ctx, fcodec, NULL) < 0) {
        PERR("avcodec_open2 failed");
        goto error;
    }
    packet = av_packet_alloc();
    if (!packet) {
        PERR("packet alloc failed");
        goto error;
    }
    frame = av_frame_alloc();
    if (!frame) {
        PERR("frame alloc failed");
        goto error;
    }

    trans_ctx = (struct TranscodeContext*)malloc(sizeof(struct TranscodeContext));
    if (!trans_ctx) {
        PERR("memory out!");
        goto error;
    }
    bzero(trans_ctx,sizeof(struct TranscodeContext));
    trans_ctx->codec_ctx = codec_ctx;
    trans_ctx->codec = fcodec;
    trans_ctx->packet = packet;
    trans_ctx->frame = frame;

    return trans_ctx;
error:
    if (packet) {
        av_packet_free(&packet);
    }
    if (frame) {
        av_frame_free(&frame);
    }
    if (codec_ctx) {
        avcodec_free_context(&codec_ctx);
    }
    return NULL;
}

struct DecodedFrame *transcode_iterate(struct TranscodeContext *trans_ctx,const char *compressed_data,int compressed_size,int *reason)
{
    AVPacket *packet = trans_ctx->packet;
    AVFrame *frame = trans_ctx->frame;
    int ret;

    *reason = 0;
    if (av_new_packet(packet, compressed_size) != 0) {
        PERR("av_new_packet failed");
        goto error;
    }
    memcpy(packet->data,compressed_data,compressed_size);

    ret = avcodec_send_packet(trans_ctx->codec_ctx, packet);
    if (ret != 0) {
        PERR("avcodec_send_packekt failed");
        *reason = AVERROR(EINVAL);
        goto error;
    }
    ret = avcodec_receive_frame(trans_ctx->codec_ctx, frame);
    if (ret < 0) {
        if (ret == AVERROR(EAGAIN)) {
            *reason = AVERROR(EAGAIN);
            goto again;
        } else {
            *reason = AVERROR(EINVAL);
            PERR("avcodec_receive_frame failed");
        }
        goto error;
    }

    int decoded_size = frame->linesize[0];
    struct DecodedFrame *dframe = malloc(sizeof(struct DecodedFrame));
    dframe->data = malloc(decoded_size);
    dframe->size = decoded_size;
    memcpy(dframe->data,frame->data[0],decoded_size);

    av_packet_unref(packet);
    av_frame_unref(frame);
    return dframe;
again:
    return NULL;
error:
    return NULL;
}

void transcode_free(struct TranscodeContext *trans_ctx)
{
    if (!trans_ctx) {
        PERR("free NULL transcode context");
        return;
    }

    if (trans_ctx->packet) {
        av_packet_free(&trans_ctx->packet);
    }
    if (trans_ctx->frame) {
        av_frame_free(&trans_ctx->frame);
    }
    if (trans_ctx->codec_ctx) {
        avcodec_free_context(&trans_ctx->codec_ctx);
    }

    free(trans_ctx);
}
