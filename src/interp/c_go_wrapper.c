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

#include "common.h"
#include "interp.tab.h"

void static inline print_value(value_t *value)
{
    switch (value->type)
    {
    case VALUE_STRING:
        printf("%s", value->value.str_t);
        break;
    case VALUE_INTEGER:
        printf("%d", value->value.int_t);
        break;
    case VALUE_FLOAT:
        printf("%f", value->value.float_t);
        break;
    default:
        fprintf(stderr, "wtf %d\n", value->type);
        break;
    }
    return;
}

extern int go_CreateTableCallback(create_table_t *param) __asm__ ("go.go_c_wrapper.CreateTableCallback");
int create_table_callback(create_table_t *param)
{
    return go_CreateTableCallback(param);
}

/* rewrite me to comminicate with Go */
int drop_table_callback(drop_table_t *param)
{
    printf("Drop table %s\n", param->name);
    return 0;
}

/* rewrite me to comminicate with Go */
int create_index_callback(create_index_t *param)
{
    printf("Create index %s on %s (%s)\n", param->index_name, param->table_name, param->col_name);
    return 0;
}

/* rewrite me to comminicate with Go */
int drop_index_callback(drop_index_t *param)
{
    printf("Drop index %s\n", param->name);    
    return 0;
}

/* rewrite me to comminicate with Go */
int select_callback(select_t *param)
{
    printf("Select * from %s\n", param->table_name);
    condition_t *cond = param->condition_list;
    while (cond)
    {
        printf("\tCondition %s %d ", cond->col_name, cond->operator);
        print_value(cond->value);
        printf("\n");
        cond = cond->next;
    }
    return 0;
}

/* rewrite me to comminicate with Go */
int insert_into_callback(insert_into_t *param)
{
    printf("Insert into %s\n", param->table_name);
    value_t *value = param->value_list;
    while (value)
    {
        printf("\tValue: ");
        print_value(value);
        printf("\n");
        value = value->next;
    }
    return 0;
}

/* rewrite me to comminicate with Go */
int delete_from_callback(delete_from_t *param)
{
    printf("Delete from %s\n", param->table_name);
    condition_t *cond = param->condition_list;
    while (cond)
    {
        printf("\tCondition: %s %d ", cond->col_name, cond->operator);
        print_value(cond->value);
        printf("\n");
        cond = cond->next;
    }
    return 0;
}

/* rewrite me to comminicate with Go */
int quit_callback(void *param)
{
    exit(0);
    return 0;
}

int execfile_callback(execfile_t *param)
{
    printf("Execfile %s\n", param->filename);
    return 0;
}
