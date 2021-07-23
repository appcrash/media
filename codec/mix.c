#include "codec.h"

struct MixContext *mix_init_context(const char *param_string,int length)
{
    AVDictionary *dict = NULL;
    struct MixContext *mix_ctx = NULL;

    //av_log_set_level(AV_LOG_VERBOSE);
    parse_param_string(param_string,length,&dict);

    mix_ctx = av_malloc(sizeof(*mix_ctx));
    bzero(mix_ctx, sizeof(*mix_ctx));
    mix_ctx->out_buffer = buffer_alloc(102400);
    if (!mix_ctx->out_buffer) {
        goto error;
    }

    if (init_mix_filter_graph(mix_ctx, dict) < 0) {
        goto error;
    }

    /*
     * after filter graph prased, sample_fmt and sample_rate are definite,
     * cache them as they would not change
     */
    if (av_opt_get_sample_fmt(mix_ctx->bufsrc1_ctx,"sample_fmt",AV_OPT_SEARCH_CHILDREN,&mix_ctx->format1) < 0) {
        goto error;
    }
    if (av_opt_get_sample_fmt(mix_ctx->bufsrc2_ctx,"sample_fmt",AV_OPT_SEARCH_CHILDREN,&mix_ctx->format2) < 0) {
        goto error;
    }
    if (av_opt_get_int(mix_ctx->bufsrc1_ctx,"sample_rate",AV_OPT_SEARCH_CHILDREN,&mix_ctx->sample_rate1) < 0) {
        goto error;
    }
    if (av_opt_get_int(mix_ctx->bufsrc2_ctx,"sample_rate",AV_OPT_SEARCH_CHILDREN,&mix_ctx->sample_rate2) < 0) {
        goto error;
    }

    return mix_ctx;
error:
    if (mix_ctx) {
        av_free(mix_ctx);
    }
    return NULL;
}

void mix_iterate(struct MixContext *mix_ctx,char *src1,int len1,char *src2,int len2,int nb_samples,int *reason)
{
    AVFrame *frame1 = av_frame_alloc();
    AVFrame *frame2 = av_frame_alloc();
    AVFrame *output = av_frame_alloc();

    *reason = 0;
    if (!frame1 || !frame2 || !output) {
        PERR("allocate frame failed");
        goto error;
    }

    frame1->nb_samples = nb_samples;
    frame1->channels = 1;
    frame1->sample_rate = mix_ctx->sample_rate1;
    frame1->format = mix_ctx->format1;
    if (av_frame_get_buffer(frame1, 1) < 0) {
        PERR("get frame buffer failed");
        goto error;
    }
    if (len1 != frame1->linesize[0]) {
        PERR("mix input1 has wrong data length");
        goto error;
    }
    memcpy(frame1->extended_data[0],src1,len1);
    frame2->nb_samples = nb_samples;
    frame2->channels = 1;
    frame2->sample_rate = mix_ctx->sample_rate2;
    frame2->format = mix_ctx->format2;
    if (av_frame_get_buffer(frame2, 1) < 0) {
        PERR("get frame buffer failed");
        goto error;
    }
    if (len2 != frame2->linesize[0]) {
        PERR("mix input2 has wrong data length");
        goto error;
    }
    memcpy(frame2->extended_data[0],src2,len2);

    if (av_buffersrc_add_frame(mix_ctx->bufsrc1_ctx, frame1) < 0) {
        PERR("add frame to mix input1 failed");
        goto error;
    }
    if (av_buffersrc_add_frame(mix_ctx->bufsrc2_ctx, frame2) < 0) {
        PERR("add frame to mix input2 failed");
        goto error;
    }

    if (av_buffersink_get_frame(mix_ctx->bufsink_ctx, output) < 0) {
        PERR("get data from buffer sink failed");
        goto error;
    }
    buffer_fill(mix_ctx->out_buffer,(const char*)output->extended_data[0],output->linesize[0]);
    goto done;
error:
    *reason = -1;
done:
    av_frame_free(&frame1);
    av_frame_free(&frame2);
    av_frame_free(&output);
}

void mix_free(struct MixContext *mix_ctx)
{
    if (mix_ctx) {
        if (mix_ctx->filter_graph)  {
            avfilter_graph_free(&mix_ctx->filter_graph);
        }
        buffer_free(&mix_ctx->out_buffer);
        av_free(mix_ctx);
    }

}
