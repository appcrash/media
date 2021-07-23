#include "codec.h"
#include "libavutil/opt.h"

/************************************************************************************************************************
Here is a BNF description of the filtergraph syntax:

NAME             ::= sequence of alphanumeric characters and '_'
FILTER_NAME      ::= NAME["@"NAME]
LINKLABEL        ::= "[" NAME "]"
LINKLABELS       ::= LINKLABEL [LINKLABELS]
FILTER_ARGUMENTS ::= sequence of chars (possibly quoted)
FILTER           ::= [LINKLABELS] FILTER_NAME ["=" FILTER_ARGUMENTS] [LINKLABELS]
FILTERCHAIN      ::= FILTER [,FILTERCHAIN]
FILTERGRAPH      ::= [sws_flags=flags;] FILTERCHAIN [;FILTERGRAPH]

libavfilter internals:

struct:
AVFilterInOut:  represent LINKLABEL, it has fields of name, pad_index, filter_ctx, mainly to make input correspond
with output, as well as both of them have the same name, then the two filter contexts are connected.

function:
*avfilter_graph_parse_ptr* parse the filter graph and connect filter context.

How the parsing is done:

loop{
  1. parse_inputs()
  2. parse_filter()
  3. link_filter_inouts()
  4. parse_outputs()
}

For example: a graph description string is like:  [in1] [in2] [in3] filter=<options...>  [out1] [out2]; ...
some variables during parsing:

*open_inputs*:  a linked list of pending AVFilterInOut whose corresponding output not found yet
*open_outputs*: a linked list of peinding AVFilterInOut whose corresponding input not found yet
*curr_inputs*: a linked list of current filter's input

step1: parse_inputs()
parse current filter's input LINKLABEL and build AVFilterInOut instances, set their names and pad index
in the order of appearance of description string, and insert them to curr_inputs. If any AVFilterInOut in
the open_outputs has the same name of an input, pull it from open_outputs instead of creating new one.

After function returns:
[in1] [in2] [in3] filter=<options...>  [out1] [out2]; ...
                  ^
curr_inputs:  [in1] -> [in2] -> [in3]

step2: parse_filter()
forward the pointer and get the filter name and options, create and initialize the filter in filter graph.

After function returns:
[in1] [in2] [in3] filter=<options...>  [out1] [out2]; ...
                                      ^
Now the current filter context is definite, we know how many input/output pads and can relate the filter context to
AVFilterInOut slot in open_inputs or open_outputs

step3: link_filter_inouts()
iterate over curr_inputs list, if the element's filter_ctx is not null, means it has definitely an
output of <filter_ctx>, meanwhile it is also an input of current filter context, so just connect
the two filter context. If the corresponding output not found, just append it to the pending input.

As current filter info is known, allocate AVFilterInOut for its output pads, then append these allocated AVFilterInOut
into curr_inputs. These newly created AVFilterInOut have been related to current filter context and pad index,
but have no name. Their names would be found in the next step.


open_inputs: [old open_inputs] -> [in1] -> [in2] -> [in3]
curr_inputs: [out1] [out2]    note: all output here are allocated but without name

step4: parse_outputs()
forward the pointer again to parse remaining output LINKLABELs.
iterate AVFilterInOut over curr_inputs. note is step3 it is actually output of current filter context. The parsed names can
be assigned to them now. And this is the right time to find corresponding input AVFilterInOut in the open_inputs. Once matched,
connect previous filter_ctx(pointed by AVFilterInOut found in open_inputs) and current filter context. If not found, just insert
it to open_outputs' head. Also relating current filter context to these AVFilterInOut before inserting them. After this function
returns, loop to step1, it would parse the inputs of next filter's input and got the corresponding output slot in open_outputs.
So insert the output slot to the head of open_outputs instead of appending to the end for performance reason.

Conclusion:

Any input or output with the same name would find each other once parsing done. The net result is filter contexts are connected as
long as you set the correct name of AVFilterInOut and pass them to *avfilter_graph_parse_ptr*. However the argument names of this
function may cause confusion, so make it clear:

inputs: they are not the input slots to the filter graph, but something that can be plugged into any output pad of any
filter context defined in the graph

outputs: they are not the output slots to the filter graph, but something that can be plugged into any input pad of any
filter context defined in the graph

Usually, a filter graph has implicit label link on both endpoints [in]/[out], and is not used alone but along with abuffer/abuffersink

                [in] ----->  filter graph ---> [out]

                  |                              |
      [abuffer]<--+                              +--> [abuffersink]
      (parameter outputs)                             (parameter inputs)


Feed data to abuffer, retrieve filtered data from abuffersink.
************************************************************************************************************************/


/*
 * prepare filter graph and connect abuffer/abuffersink to both endpoints
 */
int init_transcode_filter_graph(struct TranscodeContext *trans_ctx,const char *graph_desc_str)
{
    AVCodecContext *encode_ctx = trans_ctx->encode_ctx;
    AVCodecContext *decode_ctx = trans_ctx->decode_ctx;
    AVFilterInOut *inputs = avfilter_inout_alloc();
    AVFilterInOut *outputs = avfilter_inout_alloc();
    char args[256];

    if (!decode_ctx->channel_layout) {
        decode_ctx->channel_layout = av_get_default_channel_layout(decode_ctx->channels);
    }
    trans_ctx->filter_graph = avfilter_graph_alloc();
    snprintf(args, sizeof(args), "time_base=1/%d:sample_rate=%d:sample_fmt=%s:channel_layout=0x%"PRIx64,
             decode_ctx->sample_rate,decode_ctx->sample_rate,
             av_get_sample_fmt_name(decode_ctx->sample_fmt),decode_ctx->channel_layout);
    if (avfilter_graph_create_filter(&trans_ctx->bufsrc_ctx, avfilter_get_by_name("abuffer"),
                                     "in", args, NULL, trans_ctx->filter_graph) < 0) {
        PERR("create abuffer filter failed");
        goto error;
    }
    if (avfilter_graph_create_filter(&trans_ctx->bufsink_ctx, avfilter_get_by_name("abuffersink"),
                                     "out", NULL, NULL, trans_ctx->filter_graph) < 0) {
        PERR("create abuffersink filter failed");
        goto error;
    }

    /* abuffersink has different way to set options, as its attributes are plural */
    if (av_opt_set_bin(trans_ctx->bufsink_ctx,"sample_fmts",(uint8_t*)&encode_ctx->sample_fmt,
                       sizeof(encode_ctx->sample_fmt),AV_OPT_SEARCH_CHILDREN) < 0) {
        PERR("set abuffersink sample_fmts failed");
        goto error;
    }
    if (av_opt_set_bin(trans_ctx->bufsink_ctx,"sample_rates",(uint8_t*)&encode_ctx->sample_rate,
                       sizeof(encode_ctx->sample_rate),AV_OPT_SEARCH_CHILDREN) < 0) {
        PERR("set abuffersink sample_rate failed");
        goto error;
    }
    if (!encode_ctx->channel_layout) {
        encode_ctx->channel_layout = av_get_default_channel_layout(encode_ctx->channels);
    }
    if (av_opt_set_bin(trans_ctx->bufsink_ctx,"channel_layouts",(uint8_t*)&encode_ctx->channel_layout,
                       sizeof(encode_ctx->channel_layout),AV_OPT_SEARCH_CHILDREN) < 0) {
        PERR("set abuffersink channel_layout failed");
        goto error;
    }

    /* inputs/outputs  refer conclusion above */
    inputs->name = av_strdup("out");
    inputs->filter_ctx = trans_ctx->bufsink_ctx;
    inputs->pad_idx = 0;
    inputs->next = NULL;

    outputs->name = av_strdup("in");
    outputs->filter_ctx = trans_ctx->bufsrc_ctx;
    outputs->pad_idx = 0;
    outputs->next = NULL;

    if (avfilter_graph_parse_ptr(trans_ctx->filter_graph, graph_desc_str,&inputs,&outputs,NULL) < 0) {
        PERR("parse filter graph error");
        goto error;
    }
    if (avfilter_graph_config(trans_ctx->filter_graph, NULL) < 0) {
        PERR("config filter graph error");
        goto error;
    }

    /*
     *  encoder such as amr needs every frame size being fixed, make sure frame
     *  pulled from sink buffer is in fixed-size
     */
    if (!(encode_ctx->codec->capabilities & AV_CODEC_CAP_VARIABLE_FRAME_SIZE)) {
        //printf("set sink frame size is %d\n",encode_ctx->frame_size);
        av_buffersink_set_frame_size(trans_ctx->bufsink_ctx, encode_ctx->frame_size);
    }
    avfilter_inout_free(&inputs);
    avfilter_inout_free(&outputs);

    return 0;
error:
    avfilter_inout_free(&inputs);
    avfilter_inout_free(&outputs);
    return -1;
}


static void config_mix_input(AVFilterGraph *filter_graph,AVFilterInOut **inputs,AVFilterContext **ctx,
                             const char *opts,const char *filter_name,const char *inout_name)
{
    AVFilterInOut *inout;

    if (!inputs) {
        return;
    }
    if (avfilter_graph_create_filter(ctx, avfilter_get_by_name("abuffer"), filter_name, opts, NULL, filter_graph) < 0) {
        PERR("create abuffer filter for mix input failed");
        return;
    }
    inout = avfilter_inout_alloc();
    inout->name = av_strdup(inout_name);
    inout->filter_ctx = *ctx;
    inout->pad_idx = 0;
    inout->next = NULL;
    if (*inputs) {
        (*inputs)->next = inout;
    } else {
        *inputs = inout;
    }
}

static void config_mix_output(AVFilterGraph *filter_graph,AVFilterInOut **outputs,AVFilterContext **ctx,
                              const char *opts,const char *filter_name,const char *inout_name)
{
    AVFilterInOut *inout;
    AVDictionary *dict = NULL;
    AVDictionaryEntry *t = NULL;
    enum AVSampleFormat fmts[1];
    uint64_t channel_layouts[1];
    int sample_rate[1];

    if (avfilter_graph_create_filter(ctx, avfilter_get_by_name("abuffersink"), filter_name, NULL, NULL, filter_graph) < 0) {
        PERR("create abuffer filter for mix output failed");
        return;
    }
    if (av_dict_parse_string(&dict, opts, "=", ":", 0) < 0) {
        PERR("output options of mix context invalid");
        av_dict_free(&dict);
        return;
    }

    while ((t = av_dict_get(dict, "", t, AV_DICT_IGNORE_SUFFIX))) {
        if (!av_strcasecmp("sample_rate", t->key)) {
            sample_rate[0] = atoi(t->value);
        } else if (!av_strcasecmp("sample_fmt", t->key)){
            fmts[0] = atoi(t->value);
        } else if (!av_strcasecmp("channel_layout", t->key)) {
            channel_layouts[0] = atoi(t->value);
        }
    }
    av_dict_free(&dict);

    av_opt_set_bin(*ctx,"sample_fmts",(uint8_t*)fmts,sizeof(uint64_t),AV_OPT_SEARCH_CHILDREN);
    av_opt_set_bin(*ctx,"sample_rates",(uint8_t*)sample_rate,sizeof(int),AV_OPT_SEARCH_CHILDREN);
    av_opt_set_bin(*ctx,"channel_layouts",(uint8_t*)channel_layouts,sizeof(uint64_t),AV_OPT_SEARCH_CHILDREN);

    inout = avfilter_inout_alloc();
    inout->name = av_strdup(inout_name);
    inout->filter_ctx = *ctx;
    inout->pad_idx = 0;
    inout->next = NULL;
    *outputs = inout;
}

/*
 * use amix filter, prepare two source buffers and one sink buffer
 */
int init_mix_filter_graph(struct MixContext *mix_ctx,AVDictionary *dict)
{
    AVFilterInOut *inputs = NULL;
    AVFilterInOut *outputs = NULL;
    int configed_src = 0;
    AVDictionaryEntry *t = NULL;
    int ret = 0;
    static const char *mixgraph_desc = "[in1][in2] amix=inputs=2 ";

    mix_ctx->filter_graph = avfilter_graph_alloc();
    while ((t = av_dict_get(dict, "", t, AV_DICT_IGNORE_SUFFIX))) {
        /* setup inputs/outputs */
        if (!av_strcasecmp(t->key, "input1")) {
            /* NOTE: *inout_name* must be the same in the filter description string */
            config_mix_input(mix_ctx->filter_graph, &outputs, &mix_ctx->bufsrc1_ctx, t->value, "input1", "in1");
            configed_src++;
        } else if (!av_strcasecmp(t->key, "input2")) {
            config_mix_input(mix_ctx->filter_graph, &outputs, &mix_ctx->bufsrc2_ctx, t->value, "input2", "in2");
            configed_src++;
        } else if (!av_strcasecmp(t->key, "output")) {
            /* NOTE: "out" is the implicit name of amix filter output link */
            config_mix_output(mix_ctx->filter_graph,&inputs,&mix_ctx->bufsink_ctx,t->value,"output","out");
        }
    }

    if (configed_src != 2 || NULL == outputs) {
        PERR("inputs or outputs config error");
        goto error;
    }
    if (avfilter_graph_parse_ptr(mix_ctx->filter_graph, mixgraph_desc,&inputs,&outputs,NULL) < 0) {
        PERR("parse mix filter graph error");
        goto error;
    }
    if (avfilter_graph_config(mix_ctx->filter_graph, NULL) < 0) {
        PERR("config mix filter graph error");
        goto error;
    }

    goto done;
error:
    ret = -1;
done:
    avfilter_inout_free(&inputs);
    avfilter_inout_free(&outputs);
    return ret;
}
