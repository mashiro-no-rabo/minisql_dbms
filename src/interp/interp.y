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
static inline column_t *new_column(char *colname, datatype_t *datatype, int col_attr);
static inline condition_t *new_condition(char *colname, int operator, value_t *value);

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
    MISC_ASTERISK MISC_UNKNOWN

    /* operators */
    OPERATOR_EQ OPERATOR_NEQ OPERATOR_LT OPERATOR_GT OPERATOR_LEQ OPERATOR_GEQ

    /* build-in data types */
    DATATYPE_CHAR DATATYPE_FLOAT DATATYPE_INT

    /* reserved keywords in SQL */
    KEYWORD_AND KEYWORD_CREATE KEYWORD_DELETE KEYWORD_DROP KEYWORD_FROM 
    KEYWORD_INDEX KEYWORD_INSERT KEYWORD_INTO KEYWORD_KEY KEYWORD_ON 
    KEYWORD_PRIMARY KEYWORD_SELECT KEYWORD_TABLE KEYWORD_UNIQUE KEYWORD_VALUES 
    KEYWORD_WHERE

    /* directives in the interpreter */
    DIRECTIVE_EXECFILE DIRECTIVE_QUIT

%type <int_t>
    operator 

%type <ptr_t>
    datatype
    create_table
    create_table_column
    create_table_column_list
    drop_table
    create_index
    drop_index
    select
    condition_list
    condition
    insert_into
    value_list
    value
    delete_from
    quit
    execfile

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
    | drop_table { drop_table_callback($1); }
    | create_index { create_index_callback($1); }
    | drop_index { drop_index_callback($1); }
    | select { select_callback($1); }
    | insert_into { insert_into_callback($1); }
    | delete_from { delete_from_callback($1); }
    | quit { quit_callback($1); }
    | execfile { execfile_callback($1); }
    ;

datatype:
    DATATYPE_INT { $$ = new_datatype(DT_INT, 1); }
    | DATATYPE_FLOAT { $$ = new_datatype(DT_FLOAT, 1); }
    | DATATYPE_CHAR MISC_PARENTHESIS_L VALUE_INTEGER MISC_PARENTHESIS_R { $$ = new_datatype(DT_STRING, $3); }
    ;

create_table:
    KEYWORD_CREATE KEYWORD_TABLE MISC_IDENTIFIER MISC_PARENTHESIS_L
    create_table_column_list KEYWORD_PRIMARY KEYWORD_KEY 
    MISC_PARENTHESIS_L MISC_IDENTIFIER MISC_PARENTHESIS_R
    MISC_PARENTHESIS_R
    {
        create_table_t *ret = malloc(sizeof(create_table_t));
        ret->name = $3;
        ret->column_list = $5;
        ret->primary_key = $9;
        $$ = ret;
    }
    ;

drop_table:
    KEYWORD_DROP KEYWORD_TABLE MISC_IDENTIFIER 
    {
        drop_table_t *ret = malloc(sizeof(drop_table_t));
        ret ->name = $3; 
        $$ = ret;
    }
    ;

create_index:
    KEYWORD_CREATE KEYWORD_INDEX MISC_IDENTIFIER KEYWORD_ON MISC_IDENTIFIER
    MISC_PARENTHESIS_L MISC_IDENTIFIER MISC_PARENTHESIS_R
    {
        create_index_t *ret = malloc(sizeof(create_index_t));
        ret->index_name = $3;
        ret->table_name = $5;
        ret->col_name = $7;
        $$ = ret;
    }
    ;

drop_index:
    KEYWORD_DROP KEYWORD_INDEX MISC_IDENTIFIER
    {
        drop_index_t *ret = malloc(sizeof(drop_index_t));
        ret->name = $3;
        $$ = ret;
    }
    ;

select:
    KEYWORD_SELECT MISC_ASTERISK KEYWORD_FROM MISC_IDENTIFIER
    {
        select_t *ret = malloc(sizeof(select_t));
        ret->table_name = $4;
        ret->condition_list = NULL;
        $$ = ret; 
    }
    | KEYWORD_SELECT MISC_ASTERISK KEYWORD_FROM MISC_IDENTIFIER 
    KEYWORD_WHERE condition_list
    {
        select_t *ret = malloc(sizeof(select_t));
        ret->table_name = $4;
        ret->condition_list = $6;
        $$ = ret; 
    }
    ;

insert_into:
    KEYWORD_INSERT KEYWORD_INTO MISC_IDENTIFIER KEYWORD_VALUES
    MISC_PARENTHESIS_L value_list MISC_PARENTHESIS_R
    {
        insert_into_t *ret = malloc(sizeof(insert_into_t));
        ret->table_name = $3;
        ret->value_list = $6;
        $$ = ret;
    }
    ;

delete_from:
    KEYWORD_DELETE KEYWORD_FROM MISC_IDENTIFIER
    {
        delete_from_t *ret = malloc(sizeof(delete_from_t));
        ret->table_name = $3;
        ret->condition_list = NULL;
        $$ = ret; 
    }
    | KEYWORD_DELETE KEYWORD_FROM MISC_IDENTIFIER KEYWORD_WHERE condition_list
    {
        delete_from_t *ret = malloc(sizeof(delete_from_t));
        ret->table_name = $3;
        ret->condition_list = $5;
        $$ = ret; 
    }
    ;

quit:
    DIRECTIVE_QUIT { $$ = NULL; }
    ;

execfile:
    DIRECTIVE_EXECFILE MISC_IDENTIFIER
    {
        execfile_t *ret = malloc(sizeof(execfile_t));
        ret->filename = $2;
        $$ = ret;
    }
    ;

value_list:
    value { $$ = $1; }
    | value MISC_COMMA value_list 
    {
        value_t *head_value = $1;
        head_value->next = $3;
        $$ = head_value;
    }
    ;

condition_list:
    condition { $$ = $1; }
    | condition KEYWORD_AND condition_list
    {
        condition_t *head_condition = $1;
        head_condition->next = $3;
        $$ = head_condition;
    }
    ;

condition:
    MISC_IDENTIFIER operator value { $$ = new_condition($1, $2, $3); }
    ;

create_table_column_list:
    /* empty column list. to avoid bug caused by state reduce */ { $$ = NULL; }
    | create_table_column MISC_COMMA create_table_column_list 
    {   
        column_t *head_col = $1;   
        head_col->next = $3;
        $$ = head_col;
    }
    ;

create_table_column:
    MISC_IDENTIFIER datatype KEYWORD_UNIQUE { $$ = new_column($1, $2, COL_ATTR_UNIQUE); }
    | MISC_IDENTIFIER datatype { $$ = new_column($1, $2, COL_ATTR_NONE); }
    ;

operator:
    OPERATOR_EQ { $$ = OP_EQ; }
    | OPERATOR_NEQ { $$ = OP_NEQ; }
    | OPERATOR_LT { $$ = OP_LT; }
    | OPERATOR_GT { $$ = OP_GT; }
    | OPERATOR_LEQ { $$ = OP_LEQ; }
    | OPERATOR_GEQ { $$ = OP_GEQ; }
    ;

value:
    VALUE_INTEGER 
    { 
        value_t *ret = malloc(sizeof(value_t));
        ret->value_type = VALUE_INTEGER;
        ret->int_t = $1;
        ret->next = NULL;
        $$ = ret;
    }
    | VALUE_STRING 
    { 
        value_t *ret = malloc(sizeof(value_t));
        ret->value_type = VALUE_STRING;
        ret->str_t = $1;
        ret->next = NULL;
        $$ = ret;
    }
    | VALUE_FLOAT 
    { 
        value_t *ret = malloc(sizeof(value_t));
        ret->value_type = VALUE_FLOAT;
        ret->float_t = $1;
        ret->next = NULL;
        $$ = ret;
    }
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

static inline column_t *new_column(char *colname, datatype_t *datatype, int col_attr)
{
    column_t *ret = malloc(sizeof(column_t));
    ret->name = colname;
    ret->datatype = datatype;
    ret->attr = col_attr;
    ret->next = NULL;
    return ret;
}

static inline condition_t *new_condition(char *colname, int operator, value_t *value)
{
    condition_t *ret = malloc(sizeof(condition_t));
    ret->col_name = colname;
    ret->operator = operator;
    ret->value = value;
    ret->next = NULL;
    return ret;
}

