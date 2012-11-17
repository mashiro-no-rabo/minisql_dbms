# Source code for course work of Database System Design, 2012.09-11, Zhejiang University.
# By:
#   Shicheng XU(lightxuzju@gmail.com)
#   Xun LOU(AquarHEAD@gmail.com)
#   Pengyu CHEN(cpy.prefers.you@gmail.com)
# COPYLEFT, ALL WRONGS RESERVED.
#
#

.PHONY: all clean 

all: src
	$(MAKE) -C src
	cp src/ddddd bin/

clean:
	rm -f bin/ddddd
	$(MAKE) clean -C src
