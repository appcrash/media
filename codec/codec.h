#include <stdint.h>
#include <libavformat/avformat.h>
#include <libavcodec/avcodec.h>
#include <libavutil/avutil.h>

struct DecodedFrame
{
    char *data;
    int size;
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

    AVPacket *packet;
    AVFrame *frame;
};

#define PERR(format, ...) fprintf(stderr,"(%s:%d)#%s: "#format"\n",__FILE__,__LINE__,__FUNCTION__,## __VA_ARGS__)

struct DecodedFrame* convert_format(char *pcma_payload,int plen);
struct Payload* read_media_file(const char* file_path);
int write_media_file(char *payload,int length,const char *file_path,int codec_id,int duration);


struct TranscodeContext *transcode_init_context(const char *from_codec_name,const char *to_codec_name);
struct DecodedFrame *transcode_iterate(struct TranscodeContext *trans_ctx,const char *compressed_data,int compressed_size,int *reason);
void transcode_free(struct TranscodeContext *trans_ctx);
