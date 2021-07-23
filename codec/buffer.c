#include "codec.h"


struct DataBuffer *buffer_alloc(int capacity)
{
    struct DataBuffer *buff = NULL;

    if (capacity <= 0) {
        PERR("wrong capacity when allocate data buffer");
        goto error;
    }

    buff = av_malloc(sizeof(*buff));
    if (!buff) {
        PERR("allocate data buffer failed");
        goto error;
    }
    bzero(buff,sizeof(*buff));
    buff->data = av_malloc(capacity);
    if (!buff->data) {
        PERR("allocate data for buffer failed");
        goto error;
    }
    buff->capacity = capacity;
    return buff;
error:
    buffer_free(&buff);
    return NULL;
}

static int buffer_ensure_capacity(struct DataBuffer *buff,int new_size,int save_old_data)
{
    if (new_size > buff->capacity) {
        int newcap = FFMAX(buff->capacity * 2,new_size);
        uint8_t *newbuff = av_malloc(newcap);
        if (!newbuff) {
            PERR("enlarge data buffer failed");
            return -1;
        }
        if (save_old_data) {
            memcpy(newbuff,buff->data,buff->size);
        }
        av_free(buff->data);
        buff->data = newbuff;
        buff->capacity = new_size;
    }
    return 0;
}

int buffer_fill(struct DataBuffer *buff,const char *data,int size)
{
    if (buffer_ensure_capacity(buff,size,0) < 0) {
        return -1;
    }
    memcpy(buff->data,data,size);
    buff->size = size;
    return 0;
}

int buffer_append(struct DataBuffer *buff,const char *data,int size)
{
    if (buffer_ensure_capacity(buff, buff->size + size, 1) < 0) {
        return -1;
    }
    memcpy(&buff->data[buff->size],data,size);
    buff->size += size;
    return 0;
}

void buffer_free(struct DataBuffer **buff)
{
    if (buff && *buff) {
        if ((*buff)->data) {
            av_free((*buff)->data);
        }
        av_free(*buff);
        *buff = NULL;
    }
}
