/*
 * Source code for course work of Database System Design, 2012.09-11, Zhejiang University.
 * By:
 *   Shicheng XU(lightxuzju@gmail.com)
 *   Xun LOU(AquarHEAD@gmail.com)
 *   Pengyu CHEN(cpy.prefers.you@gmail.com)
 * COPYLEFT, ALL WRONGS RESERVED.
 */

#ifndef _INTERP_H
#define _INTERP_H

enum
{
    STATEMENT_UNFINISHED = 0,
    STATEMENT_FINISHED = 1,
};

extern int interp_init();
extern int interp_main_loop();

#endif /* _INTERP_H */
