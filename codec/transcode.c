#include <stdio.h>
#include <stdlib.h>
#include <inttypes.h>
#include <stdint.h>
#include <string.h>
#include <sys/stat.h>
#include <unistd.h>

#include "codec.h"


int need_resample(enum AVSampleFormat src_fmt,int src_sample_rate,enum AVSampleFormat dst_fmt,int dst_sample_rate)
{
    return src_fmt != dst_fmt || src_sample_rate != dst_sample_rate;
}


int decode_resample(struct TranscodeContext *trans_ctx,AVPacket *packet)
{
    AVAudioFifo *fifo = trans_ctx->fifo_queue;
    AVCodecContext *enc_ctx = trans_ctx->encode_ctx;
    AVCodecContext *dec_ctx = trans_ctx->decode_ctx;
    SwrContext *resample_ctx = trans_ctx->resample_ctx;
    AVFrame *frame = trans_ctx->frame;
    int src_rate = dec_ctx->sample_rate;
    int dst_rate = enc_ctx->sample_rate;
    int dst_sample_number;
    int ret;

    ret = avcodec_send_packet(trans_ctx->decode_ctx, packet);
    //while (av_audio_fifo_size(fifo) < enc_ctx->frame_size) {
    while(1) {
        ret = avcodec_receive_frame(dec_ctx, frame);
        if (AVERROR(EAGAIN) == ret) {
            break;
        } else if (AVERROR_EOF == ret) {
            goto eof;
        } else if (ret < 0) {
            goto error;
        } else {
            /* decoded one frame, resample it if needed, then put the resampled data to fifo queue */
            if (resample_ctx) {
                int linesize;
                dst_sample_number = av_rescale_rnd(swr_get_delay(trans_ctx->resample_ctx, src_rate) + frame->nb_samples,
                                                   dst_rate,src_rate,AV_ROUND_UP);
                printf("resample: dst_sample_number is %d\n",dst_sample_number);
                if (dst_sample_number > trans_ctx->resample_max_sample_number) {
                    if (trans_ctx->resample_buff[0]) {
                        free(trans_ctx->resample_buff[0]);
                    }
                    ret = av_samples_alloc(trans_ctx->resample_buff, &linesize, 1, dst_sample_number, enc_ctx->sample_fmt, 1);
                    if (ret < 0) {
                        PERR("av_samples_alloc failed");
                        goto error;
                    }
                    trans_ctx->resample_max_sample_number = dst_sample_number;
                }
                ret = swr_convert(resample_ctx, trans_ctx->resample_buff, dst_sample_number,(const uint8_t**)frame->extended_data, frame->nb_samples);
                if (ret < 0) {
                    PERR("swr_convert failed");
                    goto error;
                }

                if ((ret = av_audio_fifo_realloc(fifo, av_audio_fifo_size(fifo) + dst_sample_number)) < 0) {
                    PERR("av_audio_fifo_realloc failed");
                    goto error;
                }

                /* put converted samples to fifo queue */
                if (av_audio_fifo_write(fifo, (void **)trans_ctx->resample_buff,dst_sample_number) < dst_sample_number) {
                    PERR("av_audio_fifo_write failed");
                    goto error;
                }
            } else {
                // the decoded frame contains just the samples required by encoder, no resampling needed
                uint8_t *datap[1] = {0};
                datap[0] = (uint8_t*)frame->extended_data;
                dst_sample_number = frame->nb_samples;
                //printf("dst_sample_number is %d\n",dst_sample_number);
                if (av_audio_fifo_write(fifo,(void**)datap,dst_sample_number) < dst_sample_number) {
                    PERR("av_audio_fifo_write failed");
                    goto error;
                }
            }

        }

    }

    return 0;
eof:

error:
    return -1;
}

int encode(struct TranscodeContext *trans_ctx)
{
    AVAudioFifo *fifo = trans_ctx->fifo_queue;
    AVCodecContext *enc_ctx = trans_ctx->encode_ctx;
    AVFrame *frame = av_frame_alloc(); // TODO: reuse frame
    struct DataBuffer *buff = trans_ctx->out_buffer;
    AVPacket pkt;
    int ret,sample_number;

    /*
     * prepare frame, initialize and alloc buffer large enough to receive samples
     */
    sample_number = av_audio_fifo_size(fifo);
    if (enc_ctx->frame_size != 0) {
        /* frame_size == 0 means frame size is not restricted */
        sample_number = FFMIN(sample_number,enc_ctx->frame_size);
    }

    //printf("sample number is %d\n",sample_number);
    frame->nb_samples = sample_number;
    frame->channel_layout = AV_CH_LAYOUT_MONO;
    frame->format = enc_ctx->sample_fmt;
    frame->sample_rate = enc_ctx->sample_rate;
    if (av_frame_get_buffer(frame, 0) < 0) {
        PERR("av_frame_get_buffer failed");
        av_frame_free(&frame);
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
    buff->size = 0;             /* reset the buffer, append encoded data into it */
    av_init_packet(&pkt);
    pkt.data = NULL;
    pkt.size = 0;
    while(1) {
        ret = avcodec_receive_packet(enc_ctx, &pkt);
        if (AVERROR(EAGAIN) == ret) {
            break;
        } else if (AVERROR_EOF == ret) {
            goto cleanup;
        } else if (ret < 0) {
            goto error;
        } else {
            if (pkt.size > buff->capacity - buff->size) {
                /* enlarge buffer to receive encoded data */
                int newcap = FFMAX(buff->capacity * 2,pkt.size);
                uint8_t *newbuff = malloc(newcap);
                memcpy(newbuff,buff->data,buff->size);
                free(buff->data);
                buff->data = newbuff;
                buff->capacity = newcap;
            }
            memcpy(&buff->data[buff->size],pkt.data,pkt.size);
            buff->size += pkt.size;
        }
        av_packet_unref(&pkt);
    }
    return 0;
cleanup:
error:
    av_packet_unref(&pkt);
    av_frame_unref(frame);
    return -1;
}

// transcode audio, only support 1 channel for each format
struct TranscodeContext *transcode_init_context(const char *from_codec_name,const char *to_codec_name)
{
    AVCodecContext *encode_ctx = NULL;
    AVCodecContext *decode_ctx = NULL;
    SwrContext *resample_ctx = NULL;
    uint8_t **resample_buff = NULL;
    struct TranscodeContext *trans_ctx = NULL;
    AVPacket *packet = NULL;
    AVFrame *frame = NULL;
    AVCodec *fcodec = NULL,*tcodec = NULL;
    struct DataBuffer *data_buff = NULL;

    fcodec = avcodec_find_decoder_by_name(from_codec_name);
    if (!fcodec) {
        PERR("decoder not available %s",from_codec_name);
        goto error;
    }
    tcodec = avcodec_find_encoder_by_name(to_codec_name);
    if (!tcodec) {
        PERR("encoder not available %s",to_codec_name);
        goto error;
    }


    decode_ctx = avcodec_alloc_context3(fcodec);
    //decode_ctx->channels = 1;
    if (avcodec_open2(decode_ctx, fcodec, NULL) < 0) {
        PERR("avcodec_open2 failed");
        goto error;
    }
    encode_ctx = avcodec_alloc_context3(tcodec);
    encode_ctx->sample_rate = 8000;
    encode_ctx->channels = 1;
    encode_ctx->sample_fmt = tcodec->sample_fmts[0];
    if (avcodec_open2(encode_ctx,tcodec,NULL) < 0) {
        PERR("avcodec_open2 failed");
        goto error;
    }
    printf("encoder sample_rate: %d, decoder sample_rate: %d\n",encode_ctx->sample_rate,decode_ctx->sample_rate);
    printf("encoder sample_fmt: %d, decoder sample_fmt: %d\n",encode_ctx->sample_fmt,decode_ctx->sample_fmt);

    // if decoded samples are not same as the one required by encoder, resample is needed
    if (need_resample(decode_ctx->sample_fmt, decode_ctx->sample_rate, encode_ctx->sample_fmt,encode_ctx->sample_rate)) {
        resample_ctx = swr_alloc_set_opts(NULL, AV_CH_LAYOUT_MONO, encode_ctx->sample_fmt, encode_ctx->sample_rate,
                                          AV_CH_LAYOUT_MONO, decode_ctx->sample_fmt, decode_ctx->sample_rate, 0, NULL);
        if (!resample_ctx) {
            PERR("swr_alloc_set_opts failed");
            goto error;
        }
        if (swr_init(resample_ctx) < 0) {
            PERR("swr_init failed");
            goto error;
        }
        resample_buff = malloc(sizeof(uint8_t*) * 1);  // only one channel
        resample_buff[0] = NULL;

        printf("resample from src(fmt:%d, rate:%d) to dst(fmt:%d, rate:%d)\n",decode_ctx->sample_fmt,decode_ctx->sample_rate,
               encode_ctx->sample_fmt,encode_ctx->sample_rate);
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
    trans_ctx->resample_ctx = resample_ctx;
    trans_ctx->resample_buff = resample_buff;
    trans_ctx->packet = packet;
    trans_ctx->frame = frame;

    trans_ctx->fifo_queue = av_audio_fifo_alloc(decode_ctx->sample_fmt, 1, 1);
    trans_ctx->out_buffer = data_buff;

    return trans_ctx;
error:
    if (packet) {
        av_packet_free(&packet);
    }
    if (frame) {
        av_frame_free(&frame);
    }
    if (encode_ctx) {
        avcodec_free_context(&encode_ctx);
    }
    if (decode_ctx) {
        avcodec_free_context(&decode_ctx);
    }
    if (resample_ctx) {
        swr_free(&resample_ctx);
    }
    return NULL;
}


void transcode_iterate(struct TranscodeContext *trans_ctx,const char *compressed_data,int compressed_size,int *reason)
{
    AVPacket packet;

    /* prepare packet, fill data into it */
    *reason = 0;
    if (av_new_packet(&packet, compressed_size) != 0) {
        PERR("av_new_packet failed");
        goto error;
    }
    memcpy(packet.data,compressed_data,compressed_size);

    /* decode the packet, get all decoded frames until EOF or EAGAIN is met,
     * resample them if decoded sample format or rate not the same as one
     * required by encoder, and put samples into fifo queue
     */
    decode_resample(trans_ctx,&packet);

    /*
     * pull samples out from fifo queue, encode them and append encoded
     * data to transcode context's out_buffer
     */
    encode(trans_ctx);
error:
    av_packet_unref(&packet);
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
    if (trans_ctx->decode_ctx) {
        avcodec_free_context(&trans_ctx->decode_ctx);
    }
    if (trans_ctx->out_buffer) {
        if (trans_ctx->out_buffer->data) {
            free(trans_ctx->out_buffer->data);
        }
        free(trans_ctx->out_buffer);
    }

    free(trans_ctx);
}
