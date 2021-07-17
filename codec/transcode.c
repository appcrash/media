#include <stdio.h>
#include <stdlib.h>
#include <inttypes.h>
#include <stdint.h>
#include <string.h>
#include <sys/stat.h>
#include <unistd.h>

#include "codec.h"

int decode_filter(struct TranscodeContext *trans_ctx,AVPacket *packet)
{
    AVAudioFifo *fifo = trans_ctx->fifo_queue;
    AVCodecContext *dec_ctx = trans_ctx->decode_ctx;
    AVFrame *frame = av_frame_alloc();
    int ret;

    ret = avcodec_send_packet(trans_ctx->decode_ctx, packet);
    while(1) {
        ret = avcodec_receive_frame(dec_ctx, frame);
        if (AVERROR(EAGAIN) == ret) {
            break;
        } else if (AVERROR_EOF == ret) {
            goto eof;
        } else if (ret < 0) {
            goto error;
        } else {
            if (av_buffersrc_add_frame(trans_ctx->bufsrc_ctx, frame) < 0) {
                PERR("add frame to src buffer failed");
                goto error;
            }
            while(1) {
                ret = av_buffersink_get_frame(trans_ctx->bufsink_ctx,frame);
                if (AVERROR(EAGAIN) == ret || AVERROR_EOF == ret) {
                    break;
                }
                if (av_audio_fifo_write(fifo,(void**)frame->extended_data,frame->nb_samples) < frame->nb_samples) {
                    PERR("av_audio_fifo_write failed");
                    goto error;
                }
            }

        }

    }

    av_frame_free(&frame);
    return 0;
eof:

error:
    av_frame_unref(frame);
    return -1;
}

int encode(struct TranscodeContext *trans_ctx)
{
    AVAudioFifo *fifo = trans_ctx->fifo_queue;
    AVCodecContext *enc_ctx = trans_ctx->encode_ctx;
    AVFrame *frame = av_frame_alloc(); // TODO: reuse frame
    struct DataBuffer *outbuff = trans_ctx->out_buffer;
    AVPacket pkt;
    int ret,sample_number;

    outbuff->size = 0;             /* reset the buffer, append encoded data into it */
    /*
     * prepare frame, initialize and alloc buffer large enough to receive samples
     */
    sample_number = av_audio_fifo_size(fifo);
    if (0 == sample_number) {
        return 0;
    }

    if (enc_ctx->frame_size != 0) {
        /* frame_size == 0 means frame size is not restricted */
        sample_number = FFMIN(sample_number,enc_ctx->frame_size);
    }

    //printf("encoder sample number is %d, format is %d\n",sample_number,enc_ctx->sample_fmt);
    frame->nb_samples = sample_number;
    frame->channel_layout = AV_CH_LAYOUT_MONO;
    frame->format = enc_ctx->sample_fmt;
    frame->sample_rate = enc_ctx->sample_rate;
    frame->pts = AV_NOPTS_VALUE;
    if (av_frame_get_buffer(frame, 0) < 0) {
        PERR("av_frame_get_buffer failed");
        av_frame_free(&frame);
        frame = NULL;
        goto error;
    }
    /*
     * read samples from fifo queue
     */
    if (av_audio_fifo_read(fifo,(void**)&frame->data,sample_number) < sample_number) {
        PERR("av_audio_fifo_read does not get enough samples");
        goto error;
    }
    /*
     * got all samples, start encoding until EOF or EAGAIN is met,
     * collect all encoded data
     */
    ret = avcodec_send_frame(enc_ctx, frame);
    if (AVERROR_EOF == ret) {
        goto cleanup;
    } else if (ret < 0) {
        PERR("avcodec_send_frame failed");
        goto error;
    }
    av_init_packet(&pkt);
    pkt.buf = NULL;
    pkt.data = NULL;
    pkt.size = 0;

    while(1) {
        ret = avcodec_receive_packet(enc_ctx, &pkt);
        if (AVERROR(EAGAIN) == ret) {
            break;
        } else if (AVERROR_EOF == ret) {
            goto cleanup;
        } else if (ret < 0) {
            PERR("avcodec_receive_packet failed");
            goto error;
        } else {
            //printf("encoded packet size is %d\n",pkt.size);

            if (pkt.size > outbuff->capacity - outbuff->size) {
                /* enlarge buffer to receive encoded data */
                int newcap = FFMAX(outbuff->capacity * 2,pkt.size);
                uint8_t *newbuff = malloc(newcap);
                memcpy(newbuff,outbuff->data,outbuff->size);
                free(outbuff->data);
                outbuff->data = newbuff;
                outbuff->capacity = newcap;
            }
            memcpy(&outbuff->data[outbuff->size],pkt.data,pkt.size);
            outbuff->size += pkt.size;
            av_packet_unref(&pkt);
        }
    }

    av_frame_free(&frame);
    return 0;
cleanup:
error:
    if (pkt.buf) {
        av_packet_unref(&pkt);
    }
    if (frame) {
        av_frame_unref(frame);
    }
    return -1;
}

static int init_codec_context(AVCodecContext **ctx,const char *name,const char *opts,int is_encoder)
{
    AVCodec *codec;
    AVCodecContext *codec_ctx = NULL;

    if (is_encoder) {
        codec = avcodec_find_encoder_by_name(name);
    } else {
        codec = avcodec_find_decoder_by_name(name);
    }
    if (!codec) {
        PERR("no such %s (%s)",is_encoder?"encoder":"decoder",name);
        goto error;
    }
    codec_ctx = avcodec_alloc_context3(codec);
    if (!codec_ctx) {
        goto error;
    }
    /* set options to codec context here */
    if (av_set_options_string(codec_ctx, opts, "=", ",") < 0) {
        goto error;
    }
    /* ensure encoder context has a valid sample_fmt/channel_layout */
    if (is_encoder && AV_SAMPLE_FMT_NONE == codec_ctx->sample_fmt) {
        /* client forgot to set the sample_fmt option, use default one */
        if (codec->sample_fmts && codec->sample_fmts[0]) {
            /* use first supported format of encoder */
            codec_ctx->sample_fmt = codec->sample_fmts[0];
        } else {
            goto error;
        }
        if (!codec_ctx->channel_layout) {
            codec_ctx->channel_layout = av_get_default_channel_layout(codec_ctx->channels);
        }
    }

    /* to suppress experimental warnings for encoder/decoder, set this field */
    codec_ctx->strict_std_compliance = -2;
    if (avcodec_open2(codec_ctx,codec,NULL) < 0) {
        PERR("avcodec_open2 failed");
    }

    *ctx = codec_ctx;
    return 0;
error:
    if (codec_ctx) {
        avcodec_free_context(&codec_ctx);
    }
    return -1;
}


static void init_context_from_param(AVCodecContext **encode_ctx,AVCodecContext **decode_ctx,char **fg_desc,const char *param_string,int length)
{
    const char *start = param_string;
    const char *end = &param_string[length];
    const char *next;
    char *type_name,*codec_name,*opts;
    int exit = 0;
    *encode_ctx = *decode_ctx = NULL;
    *fg_desc = NULL;

    /*
     *  parse line by line
     *  1. extract [encoder/decoder type,name,options] tuple, then allocate codec context, set it with options
     *  2. in case of filter graph description, assign it to fg_desc for later use
     */
    while(start < end) {
        type_name = codec_name = opts = NULL;
        start += strspn(start," \t");
        next = av_strnstr(start, ":", end - start);
        if (!next) {
            break;
        }
        /* extract type name */
        type_name = av_strndup(start,next - start);
        if (av_strncasecmp(type_name, "encoder", 7) &&
            av_strncasecmp(type_name, "decoder", 7) &&
            av_strncasecmp(type_name, "filter_graph",12)) {
            /* type is not any one of them, stop parsing */
            exit = 1;
            goto cleanup;
        }
        start = next + 1; /* skip ":" */
        if (start >= end) {
            goto cleanup;
        }
        if (!av_strncasecmp(type_name,"filter_graph",12)) {
            goto search_line_end;
        }

        next = av_strnstr(start," ",end - start);
        if (!next) {
            break;
        }
        codec_name = av_strndup(start,next - start);
        start = next;
    search_line_end:
        start += strspn(start," \t");  /* forward to options */
        next = av_strnstr(start,"\n",end - start);
        if (!next) {
            /* can not find line end, maybe already in the last line */
            next = end;
        }
        opts = av_strndup(start,next - start);
        //printf("type_name: (%s), codec_name: (%s), options: (%s)\n",type_name,codec_name,opts);
        /* one line parsed, finish the actual work */
        if (!av_strncasecmp(type_name, "encoder",7)) {
            if (init_codec_context(encode_ctx, codec_name, opts, 1) < 0) {
                exit = 1;
                goto cleanup;
            }
        } else if (!av_strncasecmp(type_name,"decoder",7)) {
            if (init_codec_context(decode_ctx, codec_name, opts, 0) < 0) {
                exit = 1;
                goto cleanup;
            }
        } else {
            *fg_desc = av_strdup(opts);
        }

        start = next + 1; /* skip "\n" */
    cleanup:
        if (type_name) {
            av_free(type_name);
        }
        if (codec_name) {
            av_free(codec_name);
        }
        if (opts) {
            av_free(opts);
        }
        if (exit) {
            break;
        }
    }

}


/* transcode context initialization, support audio only */
struct TranscodeContext *transcode_init_context(const char *params_string,int length)
{
    AVCodecContext *encode_ctx,*decode_ctx;
    char *filter_graph_desc;
    struct TranscodeContext *trans_ctx = NULL;
    struct DataBuffer *data_buff = NULL;

    init_context_from_param(&encode_ctx,&decode_ctx,&filter_graph_desc,params_string,length);
    if (!(encode_ctx && decode_ctx)) {
        PERR("encoder/decoder initialize failed");
        goto error;
    }

    //printf("encoder sample_rate: %d, decoder sample_rate: %d\n",encode_ctx->sample_rate,decode_ctx->sample_rate);
    //printf("encoder sample_fmt: %d, decoder sample_fmt: %d\n",encode_ctx->sample_fmt,decode_ctx->sample_fmt);

    data_buff = malloc(sizeof(struct DataBuffer));
    data_buff->data = malloc(1024);
    data_buff->size = 0;
    data_buff->capacity = 1024;

    trans_ctx = (struct TranscodeContext*)malloc(sizeof(struct TranscodeContext));
    if (!trans_ctx) {
        PERR("memory out!");
        goto error;
    }
    bzero(trans_ctx,sizeof(struct TranscodeContext));
    trans_ctx->encode_ctx = encode_ctx;
    trans_ctx->decode_ctx = decode_ctx;
    trans_ctx->fifo_queue = av_audio_fifo_alloc(encode_ctx->sample_fmt, 1, 1);
    trans_ctx->out_buffer = data_buff;

    //printf("filter desc is %s\n",filter_graph_desc);
    if (init_filter_graph(trans_ctx,filter_graph_desc) < 0) {
        goto error;
    }

    return trans_ctx;
error:
    if (encode_ctx) {
        avcodec_free_context(&encode_ctx);
    }
    if (decode_ctx) {
        avcodec_free_context(&decode_ctx);
    }
    if (trans_ctx) {
        if (trans_ctx->filter_graph) {
            avfilter_graph_free(&trans_ctx->filter_graph);
        }
        free(trans_ctx);
    }

    return NULL;
}

void transcode_iterate(struct TranscodeContext *trans_ctx,char *compressed_data,int compressed_size,int *reason)
{
    AVPacket *packet;

    *reason = 0;
    /* prepare packet, fill data into it */
    packet = av_packet_alloc();
    if (av_new_packet(packet, compressed_size) != 0) {
        PERR("av_new_packet failed");
        goto error;
    }
    memcpy(packet->data,compressed_data,compressed_size);
    packet->size = compressed_size;

    /*
     * decode the packet, get all decoded frames until EOF or EAGAIN is met,
     * filter them to the samples of correct format/rate/layout required by encoder,
     * then put samples into fifo queue
     */
    if (decode_filter(trans_ctx,packet) < 0) {
        *reason = -1;
        goto error;
    }

    /*
     * pull samples out from fifo queue, encode them and append encoded
     * data to transcode context's out_buffer
     */
    if (encode(trans_ctx) < 0) {
        *reason = -1;
        goto error;
    }

error:
    av_packet_unref(packet);
}

void transcode_free(struct TranscodeContext *trans_ctx)
{
    if (!trans_ctx) {
        PERR("free NULL transcode context");
        return;
    }
    if (trans_ctx->decode_ctx) {
        avcodec_free_context(&trans_ctx->decode_ctx);
    }
    if (trans_ctx->encode_ctx) {
        avcodec_free_context(&trans_ctx->encode_ctx);
    }
    if (trans_ctx->filter_graph) {
        avfilter_graph_free(&trans_ctx->filter_graph);
    }
    if (trans_ctx->out_buffer) {
        if (trans_ctx->out_buffer->data) {
            free(trans_ctx->out_buffer->data);
        }
        free(trans_ctx->out_buffer);
    }

    free(trans_ctx);
}
