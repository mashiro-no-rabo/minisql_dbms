/* A Bison parser, made by GNU Bison 2.6.5.  */

/* Bison interface for Yacc-like parsers in C
   
      Copyright (C) 1984, 1989-1990, 2000-2012 Free Software Foundation, Inc.
   
   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   
   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.
   
   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.  */

/* As a special exception, you may create a larger work that contains
   part or all of the Bison parser skeleton and distribute that work
   under terms of your choice, so long as that work isn't itself a
   parser generator using the skeleton or a modified version thereof
   as a parser skeleton.  Alternatively, if you modify or redistribute
   the parser skeleton itself, you may (at your option) remove this
   special exception, which will cause the skeleton and the resulting
   Bison output files to be licensed under the GNU General Public
   License without this special exception.
   
   This special exception was added by the Free Software Foundation in
   version 2.2 of Bison.  */

#ifndef YY_INTERP_INTERP_TAB_H_INCLUDED
# define YY_INTERP_INTERP_TAB_H_INCLUDED
/* Enabling traces.  */
#ifndef INTERPDEBUG
# if defined YYDEBUG
#  if YYDEBUG
#   define INTERPDEBUG 1
#  else
#   define INTERPDEBUG 0
#  endif
# else /* ! defined YYDEBUG */
#  define INTERPDEBUG 0
# endif /* ! defined YYDEBUG */
#endif  /* ! defined INTERPDEBUG */
#if INTERPDEBUG
extern int interpdebug;
#endif

/* Tokens.  */
#ifndef INTERPTOKENTYPE
# define INTERPTOKENTYPE
   /* Put the tokens into the symbol table, so that GDB and other debuggers
      know about them.  */
   enum interptokentype {
     VALUE_STRING = 258,
     MISC_IDENTIFIER = 259,
     VALUE_FLOAT = 260,
     VALUE_INTEGER = 261,
     MISC_PARENTHESIS_L = 262,
     MISC_PARENTHESIS_R = 263,
     MISC_COLON = 264,
     MISC_COMMA = 265,
     MISC_SEMICOLON = 266,
     MISC_COMMENT_SINGLE_LINE = 267,
     MISC_COMMENT_MULTI_LINE = 268,
     MISC_WHITESPACE = 269,
     MISC_ASTERISK = 270,
     MISC_UNKNOWN = 271,
     OPERATOR_EQ = 272,
     OPERATOR_NEQ = 273,
     OPERATOR_LT = 274,
     OPERATOR_GT = 275,
     OPERATOR_LEQ = 276,
     OPERATOR_GEQ = 277,
     DATATYPE_CHAR = 278,
     DATATYPE_FLOAT = 279,
     DATATYPE_INT = 280,
     KEYWORD_AND = 281,
     KEYWORD_CREATE = 282,
     KEYWORD_DELETE = 283,
     KEYWORD_DROP = 284,
     KEYWORD_FROM = 285,
     KEYWORD_INDEX = 286,
     KEYWORD_INSERT = 287,
     KEYWORD_INTO = 288,
     KEYWORD_KEY = 289,
     KEYWORD_ON = 290,
     KEYWORD_PRIMARY = 291,
     KEYWORD_SELECT = 292,
     KEYWORD_TABLE = 293,
     KEYWORD_UNIQUE = 294,
     KEYWORD_VALUES = 295,
     KEYWORD_WHERE = 296,
     DIRECTIVE_EXECFILE = 297,
     DIRECTIVE_QUIT = 298
   };
#endif


#if ! defined INTERPSTYPE && ! defined INTERPSTYPE_IS_DECLARED
typedef union INTERPSTYPE
{
/* Line 2042 of yacc.c  */
#line 26 "interp.y"

    int int_t;
    float float_t;
    char *str_t;
    void *ptr_t;


/* Line 2042 of yacc.c  */
#line 116 "interp.tab.h"
} INTERPSTYPE;
# define INTERPSTYPE_IS_TRIVIAL 1
# define interpstype INTERPSTYPE /* obsolescent; will be withdrawn */
# define INTERPSTYPE_IS_DECLARED 1
#endif

extern INTERPSTYPE interplval;

#ifdef YYPARSE_PARAM
#if defined __STDC__ || defined __cplusplus
int interpparse (void *YYPARSE_PARAM);
#else
int interpparse ();
#endif
#else /* ! YYPARSE_PARAM */
#if defined __STDC__ || defined __cplusplus
int interpparse (void);
#else
int interpparse ();
#endif
#endif /* ! YYPARSE_PARAM */

#endif /* !YY_INTERP_INTERP_TAB_H_INCLUDED  */
