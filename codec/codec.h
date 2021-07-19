#pragma once

#include <stdint.h>
#include <libavformat/avformat.h>
#include <libavcodec/avcodec.h>
#include <libavutil/avutil.h>
#include <libswresample/swresample.h>
#include <libavfilter/avfilter.h>
#include <libavfilter/buffersrc.h>
#include <libavfilter/buffersink.h>
#include <libavutil/audio_fifo.h>
#include <libavutil/avstring.h>
#include <libavutil/opt.h>


struct DataBuffer
{
    uint8_t *data;
    int size;
    int capacity;               /* max size of this buffer can hold */
};


struct Payload
{
    char *data;
    int size;
    int64_t bitrate;
};

struct TranscodeContext
{
    AVCodecContext *decode_ctx;
    AVCodecContext *encode_ctx;

    AVFilterGraph *filter_graph;
    AVFilterContext *bufsrc_ctx;
    AVFilterContext *bufsink_ctx;

    AVAudioFifo *fifo_queue;

    struct DataBuffer *out_buffer;
    uint8_t is_draining;
};

#define PERR(format, ...) fprintf(stderr,"(%s:%d)#%s: "#format"\n",__FILE__,__LINE__,__FUNCTION__,## __VA_ARGS__)

struct DecodedFrame* convert_format(char *pcma_payload,int plen);
struct Payload* read_media_file(const char* file_path);
int write_media_file(char *payload,int length,const char *file_path,int codec_id,int duration);

/*
 * initialize trasncoding context, with encoder/decoder names and sample properties
 * @param param_string description string for encoder/decoder/filter_graph, for example:
 *   encoder:pcm_alaw  sample_rate=8000,channels=1 \n
 *   decoder:amrnb     sample_rate=8000,channels=1 \n
 *   filter_graph: resample
 * use av_opt_* APIs to initialize these context provided by the description string
 */
struct TranscodeContext *transcode_init_context(const char *param_string,int length);
/*
 *
 */
void transcode_iterate(struct TranscodeContext *trans_ctx,char *compressed_data,int compressed_size,int *reason);
void transcode_free(struct TranscodeContext *trans_ctx);


int init_filter_graph(struct TranscodeContext *trans_ctx,const char *graph_desc_str);
