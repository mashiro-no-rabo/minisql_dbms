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
#include "interp.h"
#include "interp.tab.h"
#include "interp.yy.h"
#include "interp_line.yy.h"

#define STRBUF_INITSIZE 0x100000

/* the least power of 2 greater than x */
static inline unsigned int gclp2(unsigned int x)
{
    x |= x >> 1;
    x |= x >> 2;
    x |= x >> 4;
    x |= x >> 8;
    x |= x >> 16;
    return x + 1;
}


/* TODO: move this class string to a seperated new file? */

typedef struct __attribute__((packed)) _string
{
    int sz;
    int n;
    char *s;
}   string;

static inline string *string_new(int sz)
{
    string *ret = malloc(sizeof(string));
    ret->sz = sz;
    ret->n = 0;
    ret->s = malloc(sz);
    ret->s[0] = 0;
    return ret;
}

static inline void string_cat(string *s1, string *s2)
{
    char *s1tail = s1->s + s1->n;
    s1->n += s2->n;
    if (s1->n >= s1->sz)
    {
        s1->sz = gclp2(s1->n);
        s1->s = realloc(s1->s, s1->sz);
    }
    strncpy(s1tail, s2->s, s2->n + 1);
    return;
}

static inline void string_clear(string *str)
{
    str->n = 0;
    str->s[0] = 0;
    return;
}

static inline void string_delete(string *str)
{
    free(str->s);
    free(str);
    return;
}

int interp_init()
{
    fputs(WELCOME_MSG, stdout);
    return 0;
}


FILE *_fin = NULL;

int interp_main_loop()
{
    string *statement = string_new(STRBUF_INITSIZE);
    int statement_finished = 1;
    while (!feof(stdin))
    {
        FILE *fin = stdin;
        if (_fin)
            fin = _fin;
        else
            fputs(statement_finished ? INTERP_PROMPT0 : INTERP_PROMPT1, stdout);
        string *inputbuf = string_new(STRBUF_INITSIZE);
        while (!feof(fin) && (!inputbuf->n || inputbuf->s[inputbuf->n - 1] != '\n'))
        {
            fgets(inputbuf->s, inputbuf->sz, fin);
            inputbuf->n = strlen(inputbuf->s);
            string_cat(statement, inputbuf);
        }
        string_delete(inputbuf);
        YY_BUFFER_STATE yybufstate;
        yybufstate = interp_line_scan_string(statement->s);
        interp_line_switch_to_buffer(yybufstate);
        statement_finished = interp_linelex();
        interp_line_delete_buffer(yybufstate);
        if (feof(fin))
            statement_finished = 1;
        if (statement_finished)
        {
            if (_fin)
            {
                fclose(_fin);
                _fin = NULL;
            }
            puts("-----------------");
            puts(statement->s);
            YY_BUFFER_STATE yybufstate;
            yybufstate = interp_scan_string(statement->s);
            interp_switch_to_buffer(yybufstate);
            interpparse();
            interp_delete_buffer(yybufstate);
            string_clear(statement);
            
        }
    }
    string_delete(statement);
    return 0;
}
