#include "codec.h"

void parse_param_string(const char *param_string,int length,AVDictionary **dict)
{
    AVDictionary *d = NULL;
    const char *start = param_string;
    const char *end = &param_string[length];
    const char *next;
    char *type_name = NULL,*value = NULL;

    if (!dict) {
        PERR("parse_param_string with null dictionary string");
        return;
    }

    /*
     *  parse line by line
     */
    while(start < end) {
        start += strspn(start," \t");
        next = av_strnstr(start, ":", end - start);
        if (!next) {
            break;
        }
        /* extract type(key) name */
        type_name = av_strndup(start,next - start);
        start = next + 1; /* skip ":" */
        if (start >= end) {
            goto cleanup;
        }
        start += strspn(start," \t");  /* forward to values */
        if (start == end) {
            /* a line start with "key:" but value is empty */
            PERR("empty value of param desc, key:%s",type_name);
            goto cleanup;
        }
        next = av_strnstr(start,"\n",end - start);
        if (!next) {
            /* can not find line end, maybe already in the last line */
            next = end;
        }
        value = av_strndup(start,next - start);
        /* one line parsed, insert to dict */
        av_dict_set(&d,type_name,value,0);
        type_name = value = NULL;
        start = next + 1; /* skip "\n" */
    }

cleanup:
    if (type_name) {
        av_free(type_name);
    }
    if (value) {
        av_free(value);
    }

    *dict = d;
}
