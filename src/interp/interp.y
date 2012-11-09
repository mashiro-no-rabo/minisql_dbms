%{
/*
 * Source code for course work of Database System Design, 2012.09-11, Zhejiang University.
 * By:
 *   Shicheng XU(lightxuzju@gmail.com)
 *   Xun LOU(AquarHEAD@gmail.com)
 *   Pengyu CHEN(cpy.prefers.you@gmail.com)
 * COPYLEFT, ALL WRONGS RESERVED.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "common.h"

static void yyerror(char const *err);
static inline datatype_t *new_datatype(int meta_datatype, int len);
static inline column_t *new_column(char *colname, datatype_t *datatype);

%}

%define api.prefix interp

%union /* typedef for yylval */
{
    int int_t;
    float float_t;
    char *str_t;
    void *ptr_t;
}

%token <str_t>
    VALUE_STRING MISC_IDENTIFIER

%token <float_t>
    VALUE_FLOAT

%token <int_t>
    VALUE_INTEGER

%token <int_t>
    /* misc */
    MISC_PARENTHESIS_L MISC_PARENTHESIS_R MISC_COLON MISC_COMMA MISC_SEMICOLON
    MISC_COMMENT_SINGLE_LINE MISC_COMMENT_MULTI_LINE MISC_WHITESPACE 
    MISC_UNKNOWN

    /* operators */
    OPERATOR_EQ OPERATOR_NEQ OPERATOR_LT OPERATOR_GT OPERATOR_LEQ OPERATOR_GEQ

    /* build-in data types */
    DATATYPE_CHAR DATATYPE_FLOAT DATATYPE_INT

    /* reserved keywords in SQL */
    KEYWORD_CREATE KEYWORD_DELETE KEYWORD_DROP KEYWORD_FROM KEYWORD_INDEX 
    KEYWORD_INSERT KEYWORD_INTO KEYWORD_KEY KEYWORD_ON KEYWORD_PRIMARY 
    KEYWORD_SELECT KEYWORD_TABLE KEYWORD_UNIQUE KEYWORD_VALUES KEYWORD_WHERE

    /* directives in the interpreter */
    DIRECTIVE_EXECFILE DIRECTIVE_QUIT
    
%type <ptr_t>
    datatype
    create_table
    create_table_column
    create_table_column_list

%start all

%%

all:
    statement_list ;

statement_list:
    /* base case */
    | statement_list statement MISC_SEMICOLON
    ;

statement:  
    /* empty statement */
    | create_table { create_table_callback($1); }
    ;

datatype:
    DATATYPE_INT { $$ = new_datatype(DT_INT, 1); }
    | DATATYPE_FLOAT { $$ = new_datatype(DT_FLOAT, 1); }
    | DATATYPE_CHAR MISC_PARENTHESIS_L VALUE_INTEGER MISC_PARENTHESIS_R { $$ = new_datatype(DT_STRING, $3); }
    ;

create_table:
    KEYWORD_CREATE KEYWORD_TABLE MISC_IDENTIFIER MISC_PARENTHESIS_L
    create_table_column_list 
    KEYWORD_PRIMARY KEYWORD_KEY 
    MISC_PARENTHESIS_L MISC_IDENTIFIER MISC_PARENTHESIS_R
    MISC_PARENTHESIS_R
    {
        create_table_t *ret = malloc(sizeof(create_table_t));
        ret->table_name = strdup($3);
        ret->columns = $5;
        ret->primary_key = strdup($9);
        $$ = ret;
    }
    ;

create_table_column_list:
    /* */ { $$ = NULL; }
    | create_table_column MISC_COMMA create_table_column_list 
    {   
        column_t *next_col = $1;   
        next_col->next = $3;
        $$ = next_col;
    }
    ;

create_table_column:
    MISC_IDENTIFIER datatype { $$ = new_column($1, $2); }
    ;

%%

static void yyerror(char const *err)
{
    fprintf(stderr, "error: %s\n", err);
    return;
}

static inline datatype_t *new_datatype(int meta_datatype, int len)
{
    datatype_t *ret = malloc(sizeof(datatype_t));
    ret->meta_datatype = meta_datatype;
    ret->len = len;
    return ret;
}

static inline column_t *new_column(char *colname, datatype_t *datatype)
{
    column_t *ret = malloc(sizeof(column_t));
    ret->colname = strdup(colname);
    ret->datatype = datatype;
    ret->next = NULL;
    return ret;
}

