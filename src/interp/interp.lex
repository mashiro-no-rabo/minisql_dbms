%{
/*
 * Source code for course work of Database System Design, 2012.09-11, Zhejiang University.
 * By:
 *   Shicheng XU(lightxuzju@gmail.com)
 *   Xun LOU(AquarHEAD@gmail.com)
 *   Pengyu CHEN(cpy.prefers.you@gmail.com)
 * COPYLEFT, ALL WRONGS RESERVED.
 */

#include <string.h>
#include "interp.tab.h"

#define yylval  interplval

%}

%option case-insensitive
%option noyywrap
%option prefix="interp"

%%

 /* reserved keywords */
and return yylval.int_t = KEYWORD_AND;
create  return yylval.int_t = KEYWORD_CREATE;
delete  return yylval.int_t = KEYWORD_DELETE;
drop    return yylval.int_t = KEYWORD_DROP;
from    return yylval.int_t = KEYWORD_FROM;
index   return yylval.int_t = KEYWORD_INDEX;
insert  return yylval.int_t = KEYWORD_INSERT;
into    return yylval.int_t = KEYWORD_INTO;
key return yylval.int_t = KEYWORD_KEY;
on  return yylval.int_t = KEYWORD_ON;
primary return yylval.int_t = KEYWORD_PRIMARY;
select  return yylval.int_t = KEYWORD_SELECT;
table   return yylval.int_t = KEYWORD_TABLE;
unique  return yylval.int_t = KEYWORD_UNIQUE;
values  return yylval.int_t = KEYWORD_VALUES;
where   return yylval.int_t = KEYWORD_WHERE;

 /* builtin data types */
char    return yylval.int_t = DATATYPE_CHAR;
float   return yylval.int_t = DATATYPE_FLOAT;
int     return yylval.int_t = DATATYPE_INT;
 /* interpreter directives */
execfile    return yylval.int_t = DIRECTIVE_EXECFILE;
quit    return yylval.int_t = DIRECTIVE_QUIT;

 /* operators */
\=   return yylval.int_t = OPERATOR_EQ;
\<\>    return yylval.int_t = OPERATOR_NEQ;
\<  return yylval.int_t = OPERATOR_LT;
\>  return yylval.int_t = OPERATOR_GT;
\<\= return yylval.int_t = OPERATOR_LEQ;
\>\= return yylval.int_t = OPERATOR_GEQ;

 /* literal values */
[0-9]+      yylval.int_t = strtol(yytext, NULL, 10); return VALUE_INTEGER;
0d[0-9]+    yylval.int_t = strtol(yytext, NULL, 10); return VALUE_INTEGER;
0b[01]+     yylval.int_t = strtol(yytext, NULL, 2); return VALUE_INTEGER;
0x[0-9a-zA-Z]+  yylval.int_t = strtol(yytext, NULL, 16); return VALUE_INTEGER; 
[-+]?[0-9]*\.?[0-9]+([eE][-+]?[0-9]+)?  yylval.float_t = strtof(yytext, NULL); return VALUE_FLOAT;
\"[^"]*\"   yylval.str_t = strdup(yytext + 1); yylval.str_t[strlen(yylval.str_t) - 1] = 0; return VALUE_STRING;
\'[^']*\'   yylval.str_t = strdup(yytext + 1); yylval.str_t[strlen(yylval.str_t) - 1] = 0; return VALUE_STRING;

 /* misc */
[_A-Za-z][_0-9A-Za-z]*  yylval.str_t = strdup(yytext); return MISC_IDENTIFIER;
\(  return yylval.int_t = MISC_PARENTHESIS_L;
\)  return yylval.int_t = MISC_PARENTHESIS_R;
\,  return yylval.int_t = MISC_COMMA;
\;  return yylval.int_t = MISC_SEMICOLON;
\*  return yylval.int_t = MISC_ASTERISK;
--[^\n]*    //return yylval.int_t = MISC_COMMENT_SINGLE_LINE;
\/\*(\/|[.\n]*[^\*]\/)*[^\/]*\*\/   //return yylval.int_t = MISC_COMMENT_MULTI_LINE; 

[ \t\n\r]+    /* white spaces, do nothing */
[^\=\<\>\(\)\,\*\; \n\t]+   return yylval.int_t = MISC_UNKNOWN; 

%%

