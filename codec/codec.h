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
    AVCodecContext *codec_ctx;
    AVCodec *codec;

    AVPacket *packet;
    AVFrame *frame;
};



struct DecodedFrame* convert_format(char *pcma_payload,int plen);
struct Payload* read_media_file(const char* file_path);
int write_media_file(char *payload,int length,const char *file_path,int codec_id);
