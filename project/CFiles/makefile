CFLAGS = -I./include  -Wall -DNON_MATLAB_PARSING -DMAX_EXT_API_CONNECTIONS=255 -fPIC
LDFLAGS = -lpthread -ldl

OS = $(shell uname -s)
ifeq ($(OS), Linux)
    CFLAGS += -D__linux
    LDFLAGS += -shared
    EXT = so
else
    CFLAGS += -D__APPLE__
    LDFLAGS += -dynamiclib -current_version 1.0 -install_name ../lib/libremoteApi.dylib
    EXT = dylib
endif

all:
	@rm -f ../lib/*.$(EXT)
	@rm -f *.o
	gcc $(CFLAGS) -c ../CFiles/extApi.c -o extApi.o
	gcc $(CFLAGS) -c ../CFiles/extApiPlatform.c -o extApiPlatform.o
	@mkdir -p ../lib
	gcc extApi.o extApiPlatform.o -o ../lib/libremoteApi.$(EXT) -lpthread -ldl $(CFLAGS) $(LDFLAGS)