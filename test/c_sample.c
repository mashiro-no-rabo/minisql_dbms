/* 
 * Sample source for course work of Database System Design, 2012.09-11, Zhejiang University.
 * Source code by Pengyu CHEN(cpy.prefers.you@gmail.com).
 * COPYLEFT, ALL WRONGS RESERVED.
 */

#include <stdio.h>

extern void go_func_sample() __asm ("go.sample.GoFuncSample");

void c_func_sample()
{
    puts("This is message from a C function.");
    go_func_sample();
    return;
}

