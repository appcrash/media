#pragma once

#include <stdint.h>
#include <libavformat/avformat.h>
#include <libavcodec/avcodec.h>
#include <libavutil/avutil.h>
#include <libswresample/swresample.h>
#include <libswscale/swscale.h>
#include <libavutil/imgutils.h>
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

/*
 * transcode from A codec to B codec
 */
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

/*
 * mix two audio streams into one
 */
struct MixContext
{
    AVFilterGraph *filter_graph;
    AVFilterContext *bufsrc1_ctx;
    AVFilterContext *bufsrc2_ctx;
    AVFilterContext *bufsink_ctx;

    enum AVSampleFormat format1,format2;
    int64_t sample_rate1,sample_rate2;

    struct DataBuffer *out_buffer;
};

/*
 * record file frame by frame
 */
struct RecordContext
{
    AVFormatContext *ctx;
};




#define PERR(format, ...) fprintf(stderr,"(%s:%d)#%s: "#format"\n",__FILE__,__LINE__,__FUNCTION__,## __VA_ARGS__)

/*
 * format apis
 */
struct Payload* read_media_file(const char* file_path);
int write_media_file(char *payload,int length,const char *file_path,int codec_id,int duration);
struct RecordContext *record_init_context(const char *file_path,const char *params);
void record_iterate(struct RecordContext *ctx,const char *buff,int32_t frame_delimits[],int nb_frame);
void record_free(struct RecordContext *ctx);


/*
 * utility apis
 *
 * the buffer class is used to minimize allocate/free overhead when exchanging
 * data between golang/c. allocate buffer and reuse it until no longer needed
 *
 * buffer_fill: fill data from the start of buffer, overwrite existing data
 * buffer_append: keep old data, append data to its end
 *
 * both methods would take care of relocate underlying memory once the capacity
 * can not fit the writing data's size
 */
struct DataBuffer *buffer_alloc(int capacity);
int buffer_fill(struct DataBuffer *buff,const char *data,int size);
int buffer_append(struct DataBuffer *buff,const char *data,int size);
void buffer_free(struct DataBuffer **buff);

/*
 * param apis
 *
 * the parameter string contains lines of key/value item, ending with "\n" each line
 * the parsing function would put them to AVDictionary.
 *
 * the goal of param apis is to minimize Golang/C binding, as new features added more
 * and more options have be defined and handled, but the parameter passing should be
 * kept simple and clean.
 */

/*
 * @param str[in]  the param description string, for example:
 *   encoder:pcm_alaw  sample_rate=8000,channels=1 \n
 *   decoder:amrnb     sample_rate=8000,channels=1 \n
 *   filter_graph: resample
 *
 *   each line is of <key>:<value> \n, the last line's \n is optional
 *
 * @param dict[out] parsed string would be saved in this AVDictionary as key/value pair
 */
void parse_param_string(const char *str,int length,AVDictionary **dict);

/*
 * initialize trasncoding context, with encoder/decoder names and sample properties
 * @param param_string description string for encoder/decoder/filter_graph
 */
struct TranscodeContext *transcode_init_context(const char *param_string,int length);
/*
 * @param compressed_data  the source encoded audio data
 * @param reason[out]  set to 0 only on success
 */
void transcode_iterate(struct TranscodeContext *trans_ctx,char *compressed_data,int compressed_size,int *reason);
/*
 * free all allocated resources in this context
 */
void transcode_free(struct TranscodeContext *trans_ctx);


struct MixContext *mix_init_context(const char *param_string,int length);
void mix_iterate(struct MixContext *mix_ctx,char *src1,int len1,char *src2,int len2,int samples1,int samples2,int *reason);
void mix_free(struct MixContext *mix_ctx);



/*
 * filter setup routines
 */
int init_transcode_filter_graph(struct TranscodeContext *trans_ctx,const char *graph_desc_str);
int init_mix_filter_graph(struct MixContext *mix_ctx,AVDictionary *dict);
