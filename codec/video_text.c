#include "codec.h"



static FT_Library ft_library;
static int ft_inited = 0;

static void encode(struct VideoContext *ctx, AVFrame *frame)
{
    int ret;
    AVCodecContext *enc_ctx = ctx->encode_ctx;
    AVPacket *pkt = NULL;
    /* send the frame to the encoder */
    if (frame)
        printf("Send frame %3"PRId64"\n", frame->pts);

    ret = avcodec_send_frame(enc_ctx, frame);
    if (ret < 0) {
        fprintf(stderr, "Error sending a frame for encoding\n");
        exit(1);
    }

    ctx->nb_packet = 0;
    while (ret >= 0) {
        pkt = av_packet_alloc();
        ret = avcodec_receive_packet(enc_ctx, pkt);
        if (ret == AVERROR(EAGAIN) || ret == AVERROR_EOF) {
            goto done;
        } else if (ret < 0) {
            fprintf(stderr, "Error during encoding\n");
            goto done;
        }
        printf("packet data is %p\n",pkt->data);
        ctx->packet_data[ctx->nb_packet] = pkt;
        ctx->nb_packet++;
    }
done:
    av_packet_free(&pkt);
}

static void draw_glyph(char *surface_buffer,char *bitmap_buffer,int x,int y,
                       int surface_width,int width,int height)
{
    int xx,yy,p,q;

    printf("x:%d, y:%d, surface width:%d, bitmap_width:%d, bitmap_height:%d\n",x,y,surface_width,width,height);

    for (int i = 0; i < width;i++) {
        for(int j = 0;j < height;j++) {
            xx = x + i;
            yy = y + j;
            p = (yy * surface_width + xx) * 3;
            q = j * width + i;
            surface_buffer[p] = bitmap_buffer[q];
        }
    }
}

struct VideoContext *video_init()
{
    AVCodec *codec;
    AVCodecContext *c;
    int ret;
    FT_UInt32 code = 0x6211;

    codec = avcodec_find_encoder_by_name("libx264rgb");
    if (!codec) {
        fprintf(stderr,"can not find x264 encoder\n");
        return NULL;
    }
    c = avcodec_alloc_context3(codec);
    if (!c) {
        fprintf(stderr,"create codec context failed\n");
        return NULL;
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
        goto error;
    }
    if (!ft_inited) {
        if (FT_Init_FreeType(&ft_library)) {
            fprintf(stderr,"failed to init ft library\n");
            goto error;
        }
        ft_inited = 1;
    }

    struct VideoContext *ctx = malloc(sizeof(struct VideoContext));

    int error = FT_New_Face(ft_library, "/usr/share/fonts/truetype/wqy/wqy-microhei.ttc", 0, &ctx->face);

    if (error) {
        fprintf(stderr,"failed to create new face\n");
        goto error;
    }
    /* if (error = FT_Set_Char_Size(face, 50 * 64, 0, 100, 0)) { */
    /*     fprintf(stderr,"failed to set char size:0x%x\n",error); */
    /* } */
    if (error = FT_Set_Pixel_Sizes(ctx->face, 0, 200)) {
        fprintf(stderr,"failed to set pixel size: 0x%x\n",error);
        goto error;
    }
    if (error = FT_Load_Char(ctx->face, code, FT_LOAD_RENDER)) {
        fprintf(stderr,"failed to load char 0x%x\n",error);
        goto error;
    }
    if (FT_Get_Glyph(ctx->face->glyph, &ctx->glyph)) {
        fprintf(stderr,"failed to get glyph\n");
        goto error;
    }
    if (error = FT_Glyph_To_Bitmap(&ctx->glyph, FT_RENDER_MODE_NORMAL, 0, 1)) {
        fprintf(stderr,"to bitmap error:0x%x\n",error);
    }
    FT_BitmapGlyph bglyph = (FT_BitmapGlyph)ctx->glyph;
    ctx->bitmap = bglyph->bitmap;

    ctx->encode_ctx = c;
    ctx->iteration = 0;
    ctx->nb_packet = 0;
    return ctx;

error:
    return NULL;
}

void video_iterate(struct VideoContext *ctx)
{
    int x,y,ret;
    AVCodecContext *c = ctx->encode_ctx;
    AVFrame *frame = av_frame_alloc();
    FT_Bitmap bitmap = ctx->bitmap;
    int i = ctx->iteration;
    frame->format = c->pix_fmt;
    frame->width = c->width;
    frame->height = c->height;
    ret = av_frame_get_buffer(frame, 8);
    if (ret < 0) {
        fprintf(stderr,"could not allocate video frame buffer\n");
        return;
    }

    av_frame_make_writable(frame);
    //fprintf(stderr,"linesize is %d,width is %d\n",frame->linesize[0],c->width);

    for (y = 0; y < c->height; y++) {
        for (x = 0; x < c->width; x++) {
            //frame->data[0][y * frame->linesize[0] + x] = x + y + i * 3;
            int pos = (y * frame->width + x) * 3;
            uint8_t *p1 = &frame->data[0][pos];
            //uint8_t *p2 = &frame->data[1][pos];

            if (x < c->width / 2) {
                p1[0] = i * 2 % 255;
            } else {
                p1[1] = i % 255;
            }

        }
    }
    FT_GlyphSlot slot = ctx->face->glyph;
    draw_glyph(frame->data[0], bitmap.buffer, 100, 100, frame->width, bitmap.width,bitmap.rows);
    frame->pts = i;
    ctx->iteration++;

    encode(ctx,frame);
}
