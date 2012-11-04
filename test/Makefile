# Makefile for course work of Database System Design, 2012.09-11, Zhejiang University.
# Source code by:
#   Shicheng XU(lightxuzju@gmail.com)
#   Xun LOU(AquarHEAD@gmail.com)
#   Pengyu CHEN(cpy.prefers.you@gmail.com)
# COPYLEFT, ALL WRONGS RESERVED.
#
#

LD	= gccgo
CC	= gcc
GOC	= gccgo
LDFLAGS = 
CCFLAGS	= 
GOCFLAGS	= 

all: sample
	
sample: c_sample.o go_sample.o main.o
	$(LD) $(LDFLAGS) -o sample c_sample.o go_sample.o main.o

c_sample.o: c_sample.c
	$(CC) $(CCFLGS) -c -o c_sample.o c_sample.c

go_sample.o: go_sample.go
	$(GOC) $(GOCFLGAS) -c -o go_sample.o go_sample.go

main.o:	main.go
	$(GOC) $(GOCFLAGS) -c -o main.o main.go

clean_sample:
	rm -f sample c_sample.o go_sample.o main.o
	
clean: clean_sample
	
