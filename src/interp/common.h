/*
 * Source code for course work of Database System Design, 2012.09-11, Zhejiang University.
 * By:
 *   Shicheng XU(lightxuzju@gmail.com)
 *   Xun LOU(AquarHEAD@gmail.com)
 *   Pengyu CHEN(cpy.prefers.you@gmail.com)
 * COPYLEFT, ALL WRONGS RESERVED.
 */

#ifndef _COMMON_H
#define _COMMON_H

#define PROJECT "DDDDD: Diminutively Designed & Desperately Debugged Database"
#define VERSION "0.0.1"
#define AUTHOR  \
    AUTHOR0 "\n" \
    AUTHOR1 "\n" \
    AUTHOR2 "\n"
#define AUTHOR0 "Shicheng XU(lightxuzju@gmail.com)"
#define AUTHOR1 "Xun LOU(AquarHEAD@gmail.com)"
#define AUTHOR2 "Pengyu CHEN(cpy.prefers.you@gmail.com)"
#define COYPRIFHT   "COPYLEFT, ALL WRONGS RESERVED."

#define WELCOME_MSG \
    PROJECT "\n" \
    "version " VERSION "\n" \
    "Enter \"help\" for help\n" \
    "Enter SQL statements terminated with a \";\"" "\n"

#define INTERP_PROMPT0  "ddddd> "
#define INTERP_PROMPT1  "  ...> "

enum _DATATYPE
{
    DT_INT = 1,
    DT_STRING,
    DT_FLOAT,
};

enum _COL_ATTR
{
    COL_ATTR_NONE = 0,
    COL_ATTR_UNIQUE,
};

typedef struct __attribute__ ((__packed__)) _value_t
{
    int type;
    union
    {
        int int_t;
        char *str_t;
        float float_t;
    }   value;
    struct _value_t *next;
}   value_t;

typedef struct __attribute__((__packed__)) _datatype_t
{
    int meta_datatype; 
    int len;    /* for strings only */
    char test;
}   datatype_t;

typedef struct __attribute__((__packed__)) _column_t
{
    char *name;
    datatype_t *datatype;
    int attr;
    struct _column_t *next;
}   column_t;

typedef struct __attribute__((__packed__)) _create_table_t
{
    char *name;
    column_t *column_list;
    char *primary_key;
}   create_table_t;

typedef struct __attribute__((__packed__)) _drop_table_t
{
    char *name;
}   drop_table_t;

typedef struct __attribute__((__packed__)) _create_index_t
{
    char *index_name;
    char *table_name;
    char *col_name;
}   create_index_t;

typedef struct __attribute__((__packed__)) _drop_index_t
{
    char *name;
}   drop_index_t;

typedef struct __attribute__((__packed__)) _condition_t
{
    char *col_name;
    int operator;
    value_t *value;
    struct _condition_t *next;
}   condition_t;

typedef struct __attribute__((__packed__)) _select_t
{
    char *table_name;
    condition_t *condition_list;
}   select_t;

typedef struct __attribute__((__packed__)) _insert_into_t
{
    char *table_name;
    value_t *value_list;
}   insert_into_t;

typedef select_t delete_from_t;

typedef struct _execfile_t
{
    char *filename;
}   execfile_t;

extern int create_table_callback(create_table_t *param);
extern int drop_table_callback(drop_table_t *param);
extern int create_index_callback(create_index_t *param);
extern int drop_index_callback(drop_index_t *param);
extern int insert_into_callback(insert_into_t *param);
extern int delete_from_callback(delete_from_t *param);
extern int exit_callback();
extern int execfile_callback(execfile_t *param);

#endif /* _COMMON_H */
