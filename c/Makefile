
CC ?=		gcc
CFLAGS ?=

OBJS=	main.o filevercmp.o

all: filevercmp

clean:
	rm -f filevercmp ${OBJS}

.PHONY: all clean

filevercmp.o: filevercmp.c filevercmp.h
	$(CC) $(CPPFLAGS) $(CFLAGS) -c filevercmp.c

main.o: main.c filevercmp.h
	$(CC) $(CPPFLAGS) $(CFLAGS) -c main.c

filevercmp: ${OBJS}
	$(CC) $(LDFLAGS) $(CPPFLAGS) $(CFLAGS) -o filevercmp ${OBJS}
