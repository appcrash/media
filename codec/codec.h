#include <stdint.h>
#include <libavformat/avformat.h>
#include <libavcodec/avcodec.h>
#include <libavutil/avutil.h>
#include <libswresample/swresample.h>
#include <libavutil/audio_fifo.h>

struct DecodedFrame
{
    char *data;
    int size;
};

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

    SwrContext *resample_ctx;
    uint8_t **resample_buff;  // pointer of an array of pointer, the number of array is channel number
    int resample_max_sample_number;

    AVAudioFifo *fifo_queue;

    struct DataBuffer *out_buffer;
};

#define PERR(format, ...) fprintf(stderr,"(%s:%d)#%s: "#format"\n",__FILE__,__LINE__,__FUNCTION__,## __VA_ARGS__)

struct DecodedFrame* convert_format(char *pcma_payload,int plen);
struct Payload* read_media_file(const char* file_path);
int write_media_file(char *payload,int length,const char *file_path,int codec_id,int duration);

/*
 * initialize trasncoding context, with src/dst codec names and sample properties
 * @param to_sample_bitrate is optional, ignored when set to 0
 */
struct TranscodeContext *transcode_init_context(const char *from_codec_name,int from_sample_rate,
                                                const char *to_codec_name,int to_sample_rate,int to_sample_bitrate);
/*
 *
 */
void transcode_iterate(struct TranscodeContext *trans_ctx,char *compressed_data,int compressed_size,int *reason);
void transcode_free(struct TranscodeContext *trans_ctx);
