# Source code for course work of Database System Design, 2012.09-11, Zhejiang University.
# By:
#   Shicheng XU(lightxuzju@gmail.com)
#   Xun LOU(AquarHEAD@gmail.com)
#   Pengyu CHEN(cpy.prefers.you@gmail.com)
# COPYLEFT, ALL WRONGS RESERVED.
#
#

LD	= gccgo
CC	= gcc
GOC	= gccgo
YACC	= bison
LEX	= flex
LDFLAGS = -g
CCFLAGS	= 
GOCFLAGS	= 
DIRS	= bufman catman common core interp main recman
OBJS	= bufman.o catman.o common.o core.o interp.o main.o recman.o

# override the builtin rules by GNU make
%.c: %.y
%.c: %.l

.PHONY: all $(DIRS)

all: Makefile ddddd $(DIRS)
	
ddddd: $(OBJS)
	$(LD) $(LDFLAGS) -o ddddd $(OBJS)

bufman.o: bufman
	$(MAKE) -C bufman
	cp bufman/bufman.o .

catman.o: catman common.o
	$(MAKE) -C catman
	cp catman/catman.o .

common.o: common
	$(MAKE) -C common
	cp common/common.o .

core.o: core catman.o common.o recman.o
	$(MAKE) -C core
	cp core/core.o .
	
interp.o: interp catman.o common.o core.o
	$(MAKE) -C interp 
	cp interp/interp.o .

main.o: main bufman.o catman.o common.o core.o interp.o recman.o
	$(MAKE) -C main
	cp main/main.o .

recman.o: recman common.o
	$(MAKE) -C recman
	cp recman/recman.o .
	
clean:
	rm -f ddddd $(OBJS)
	for d in $(DIRS); do $(MAKE) clean -C $$d ; done
	
