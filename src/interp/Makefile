# Source code for course work of Database System Design, 2012.09-11, Zhejiang University.
# By:
#   Shicheng XU(lightxuzju@gmail.com)
#   Xun LOU(AquarHEAD@gmail.com)
#   Pengyu CHEN(cpy.prefers.you@gmail.com)
# COPYLEFT, ALL WRONGS RESERVED.
#
#

LD	= ld 
CC	= gcc
GOC	= gccgo
YACC	= bison #/usr/local/bin/bison
LEX	= flex
LDFLAGS = -g
CCFLAGS	= -g
GOCFLAGS	= -g
OBJS	= interp.tab.o interp.yy.o c_go_wrapper.o _interp.o interp_line.yy.o go_c_wrapper.o

# override the builtin rules by GNU make
%.c: %.y
%.c: %.l

all: interp.o

interp.o: $(OBJS)
	$(LD) -r -o interp.o $(OBJS)
	
interp.tab.o: interp.tab.c
	$(CC) $(CCFLGAS) -c -o interp.tab.o interp.tab.c

interp.tab.c interp.tab.h: interp.y common.h
	$(YACC) -d -o interp.tab.c interp.y

interp.yy.o: interp.yy.c 
	$(CC) $(CCFLAGS) -c -o interp.yy.o interp.yy.c

interp.yy.c interp.yy.h: interp.lex interp.tab.h
	$(LEX) --header-file=interp.yy.h -o interp.yy.c interp.lex

interp_line.yy.o: interp_line.yy.c
	$(CC) $(CCFLAGS) -c -o interp_line.yy.o interp_line.yy.c

interp_line.yy.c interp_line.yy.h: interp_line.lex interp.h
	$(LEX) --header-file=interp_line.yy.h -o interp_line.yy.c interp_line.lex

c_go_wrapper.o: c_go_wrapper.c common.h interp.tab.h
	$(CC) $(CFLAGS) -c -o c_go_wrapper.o c_go_wrapper.c

_interp.o: interp.c interp.h interp.tab.h common.h interp_line.yy.h
	$(CC) $(CFLAGS) -c -o _interp.o interp.c	

go_c_wrapper.o: go_c_wrapper.go
	$(GOC) $(GOCFLAGS) -c -o go_c_wrapper.o go_c_wrapper.go

clean:
	rm -f interp.o $(OBJS) interp.tab.[hc] interp.yy.[hc] interp_line.yy.[hc]
	
