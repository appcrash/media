CC=gcc
CFLAGS=`pkg-config --cflags libavformat libavcodec libavutil libswresample libavfilter freetype2`
LDFLAGS=`pkg-config --libs libavformat libavcodec libavutil libswresample libavfilter freetype2`

OBJS:=$(patsubst %.c,%.o,$(wildcard *.c))

# to generate compile_commands.json:
# pip3 install compiledb
# compiledb make

libcodec.a: $(OBJS)
	ar rcs $@ $^

%.o: %.c
	$(CC) -c -o $@ $< $(CFLAGS)

.PHONY: clean

clean:
	rm *.o *.a	
