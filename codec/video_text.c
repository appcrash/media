#include "codec.h"


#include <ft2build.h>
#include FT_FREETYPE_H
#include FT_GLYPH_H
#include FT_STROKER_H

static FT_Library ft_library;

static void encode(AVCodecContext *enc_ctx, AVFrame *frame, AVPacket *pkt, FILE *outfile)
{
    int ret;

    /* send the frame to the encoder */
    if (frame)
        printf("Send frame %3"PRId64"\n", frame->pts);

    ret = avcodec_send_frame(enc_ctx, frame);
    if (ret < 0) {
        fprintf(stderr, "Error sending a frame for encoding\n");
        exit(1);
    }

    while (ret >= 0) {
        ret = avcodec_receive_packet(enc_ctx, pkt);
        if (ret == AVERROR(EAGAIN) || ret == AVERROR_EOF)
            return;
        else if (ret < 0) {
            fprintf(stderr, "Error during encoding\n");
            exit(1);
        }

        printf("Write packet %3"PRId64" (size=%5d)\n", pkt->pts, pkt->size);
        fwrite(pkt->data, 1, pkt->size, outfile);
        av_packet_unref(pkt);
    }
}

void video_render()
{
    AVCodec *codec;
    AVCodecContext *c;
    int x,y,ret;

    if (FT_Init_FreeType(&ft_library)) {
        fprintf(stderr,"failed to init ft library\n");
        return;
    }

    codec = avcodec_find_encoder_by_name("libx264rgb");
    if (!codec) {
        fprintf(stderr,"can not find x264 encoder\n");
        return;
    }
    c = avcodec_alloc_context3(codec);
    if (!c) {
        fprintf(stderr,"create codec context failed\n");
        return;
    }


    c->bit_rate = 400000;
    c->width = 400;
    c->height = 300;
    c->time_base = (AVRational){1,25};
    c->framerate = (AVRational){25,1};
    c->gop_size = 10;
    c->max_b_frames = 1;
    //c->pix_fmt = AV_PIX_FMT_YUV420P;
    c->pix_fmt = AV_PIX_FMT_RGB24;


    ret = avcodec_open2(c, codec, NULL);
    if (ret < 0) {
        fprintf(stderr,"failed to open codec\n");
        return;
    }

    FILE *f = fopen("output.mpg","wb");

    AVPacket *pkt = av_packet_alloc();
    AVFrame *frame = av_frame_alloc();
    frame->format = c->pix_fmt;
    frame->width = c->width;
    frame->height = c->height;
    ret = av_frame_get_buffer(frame, 8);
    if (ret < 0) {
        fprintf(stderr,"could not allocate video frame buffer\n");
        return;
    }

    FT_Face face;
    int error = FT_New_Face(ft_library, "/usr/share/fonts/truetype/noto/NotoSansMyanmar-Bold.ttf", 0, &face);
    if (error) {
        fprintf(stderr,"failed to create new face\n");
        return;
    }

    for (int i = 0; i < 250; i++) {
        av_frame_make_writable(frame);
        //fprintf(stderr,"linesize is %d,width is %d\n",frame->linesize[0],c->width);

        /* prepare a dummy image */
        /* Y */
        for (y = 0; y < c->height; y++) {

            for (x = 0; x < c->width; x++) {
                //frame->data[0][y * frame->linesize[0] + x] = x + y + i * 3;
                int pos = (y * frame->width + x) * 3;
                uint8_t *p1 = &frame->data[0][pos];
                //uint8_t *p2 = &frame->data[1][pos];

                if (x < c->width / 2) {
                    p1[0] = i * 2 % 255;
                } else {
                    p1[1] = i;
                }
            }
        }

        /* Cb and Cr */
        /* for (y = 0; y < c->height/2; y++) { */
        /*     for (x = 0; x < c->width/2; x++) { */
        /*         frame->data[1][y * frame->linesize[1] + x] = 128 + y + i * 2; */
        /*         frame->data[2][y * frame->linesize[2] + x] = 64 + x + i * 5; */
        /*     } */
        /* } */

        frame->pts = i;

        encode(c,frame,pkt,f);
    }

    encode(c,NULL,pkt,f);


    uint8_t endcode[] = { 0, 0, 1, 0xb7 };
    fwrite(endcode, 1, sizeof(endcode), f);
    fclose(f);
    avcodec_free_context(&c);
    av_frame_free(&frame);
    av_packet_free(&pkt);
}
