#include <stdio.h>
#include <stdlib.h>
#include <inttypes.h>
#include <stdint.h>
#include <string.h>
#include <libavformat/avformat.h>
#include <libavcodec/avcodec.h>
#include <libavutil/avutil.h>

#include "codec.h"

struct DecodedFrame *convert_format(char *pcma_payload,int plen)
{ /* FILE *fp; */
    /* size_t plen; */
    /* uint8_t *payload = av_malloc(PSIZE); */

    /* fp = fopen("/home/yh/develop/xmedia/media_dir/100203.wav","rb"); */
    /* plen = fread(payload,1,PSIZE,fp); */
    /* fclose(fp); */

    /* uint8_t *orig_payload = payload; */
    /* payload = (uint8_t*)memmem(payload,plen,"data",4); */
    /* unsigned int diff = (payload - orig_payload); */
    /* plen -= diff; */
    /* payload = orig_payload + diff; */


    AVCodec *codec = avcodec_find_decoder(AV_CODEC_ID_PCM_ALAW);
    AVCodecContext *context = avcodec_alloc_context3(codec);
    AVPacket *packet = av_packet_alloc();
    AVFrame *frame = av_frame_alloc();
    av_packet_from_data(packet, (uint8_t*)pcma_payload, plen);

    context->sample_rate = 8000;
    //context->sample_fmt = AV_SAMPLE_FMT_S32;
    context->channels = 1;
    //context->channel_layout = AV_CH_LAYOUT_MONO;

    if (avcodec_open2(context, codec, NULL) < 0) {
        printf("avcodec_open2 error");
        exit(1);
    }


    int ret;
    ret = avcodec_send_packet(context, packet);
    //printf("avcodec_send_packet ret %d\n",ret);
    ret = avcodec_receive_frame(context, frame);
    //printf("avcodec_receive_frame ret %d\n",ret);
    //printf("frame nb_samples is %d\n",frame->nb_samples);

    int sample_size = av_get_bytes_per_sample(context->sample_fmt);
    //printf("sample size is %d\n",sample_size);
    uint8_t *decoded_data = frame->data[0];
    int decoded_size = sample_size * frame->nb_samples;
    void *decoded_copy = malloc(decoded_size);
    memcpy(decoded_copy,decoded_data,decoded_size);
    struct DecodedFrame *df = (struct DecodedFrame*)malloc(sizeof(struct DecodedFrame));
    df->data = decoded_copy;
    df->size = decoded_size;


    /* fp = fopen("pcm-s16.output","wb"); */
    /* fwrite(decoded_data,sample_size,frame->nb_samples,fp); */
    /* fclose(fp); */

    avcodec_free_context(&context);
    av_packet_free(&packet);
    av_frame_free(&frame);

    return df;
}
