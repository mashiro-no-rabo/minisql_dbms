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
LDFLAGS = 
CCFLAGS	= 
GOCFLAGS	= 

# override the builtin rules by GNU make
%.c: %.y
%.c: %.l

all: Makefile ddddd
	
ddddd: interp.tab.o interp.yy.o main.o c_go_wrapper.o interp.o interp_line.yy.o
	$(LD) $(LDFLAGS) -o ddddd interp.tab.o interp.yy.o main.o c_go_wrapper.o interp.o interp_line.yy.o

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

main.o:	main.go
	$(GOC) $(GOCFLAGS) -c -o main.o main.go

c_go_wrapper.o: c_go_wrapper.c common.h interp.tab.h
	$(CC) $(CFLAGS) -c -o c_go_wrapper.o c_go_wrapper.c

interp.o: interp.c interp.h interp.tab.h common.h interp_line.yy.h
	$(CC) $(CFLAGS) -c -o interp.o interp.c	

clean:
	rm -f ddddd *.o interp.tab.[hc] interp.yy.[hc] interp_line.yy.[hc]
	
