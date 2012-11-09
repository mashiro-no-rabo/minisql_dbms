%{
/*
 * Source code for course work of Database System Design, 2012.09-11, Zhejiang University.
 * By:
 *   Shicheng XU(lightxuzju@gmail.com)
 *   Xun LOU(AquarHEAD@gmail.com)
 *   Pengyu CHEN(cpy.prefers.you@gmail.com)
 * COPYLEFT, ALL WRONGS RESERVED.
 */

#include "interp.h"

%}

%option case-insensitive
%option noyywrap
%option prefix="interp_line"

%%

 /* */
[ \n\t]*\n  return STATEMENT_FINISHED;
(?s:.)*;[ \t]*(\-\-[^\n]*)?\n  return STATEMENT_FINISHED; 

 /* */
(?s:.)*  return STATEMENT_UNFINISHED; 

%%

