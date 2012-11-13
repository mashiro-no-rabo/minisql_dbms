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
    column_t *columns;
    char *primary_key;
}   create_table_t;

typedef struct __attribute__((__packed__)) _drop_table_t
{
    char *name;
}   drop_table_t;

extern int create_table_callback(create_table_t *param);
extern int drop_table_callback(drop_table_t *param);

#endif /* _COMMON_H */
