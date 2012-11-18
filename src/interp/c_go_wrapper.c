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
#include "interp.yy.h"

void static inline print_value(value_t *value)
{
    switch (value->value_type)
    {
    case VALUE_STRING:
        printf("%s", value->str_t);
        break;
    case VALUE_INTEGER:
        printf("%d", value->int_t);
        break;
    case VALUE_FLOAT:
        printf("%f", value->float_t);
        break;
    default:
        fprintf(stderr, "wtf %d\n", value->value_type);
        break;
    }
    return;
}

extern int go_CreateTableCallback(create_table_t *param) __asm__ ("go.go_c_wrapper.CreateTableCallback");
int create_table_callback(create_table_t *param)
{
    /*
    printf("Create table %s\n", param->name);
    column_t *col = param->column_list;
    while (col)
    {
        printf("\tTable column %s of type %d with attribute %s\n", col->name, 
            col->datatype->meta_datatype, col->attr == COL_ATTR_UNIQUE ? "unique" : "none");
        col = col->next;
    }
    printf("\tPrimary key %s\n", param->primary_key);
    */
    return go_CreateTableCallback(param);
}

extern int go_DropTableCallback(drop_table_t *param) __asm__ ("go.go_c_wrapper.DropTableCallback");
int drop_table_callback(drop_table_t *param)
{
    //printf("Drop table %s\n", param->name);
    return go_DropTableCallback(param);
}

extern int go_CreateIndexCallback(create_index_t *param) __asm__ ("go.go_c_wrapper.CreateIndexCallback");
int create_index_callback(create_index_t *param)
{
    //printf("Create index %s on %s (%s)\n", param->index_name, param->table_name, param->col_name);
    return go_CreateIndexCallback(param);
}

extern int go_DropIndexCallback(drop_index_t *param) __asm__ ("go.go_c_wrapper.DropIndexCallback");
int drop_index_callback(drop_index_t *param)
{
    //printf("Drop index %s\n", param->name);    
    return go_DropIndexCallback(param);
}

extern int go_SelectCallback(select_t *param) __asm__ ("go.go_c_wrapper.SelectCallback");
int select_callback(select_t *param)
{
    /*
    printf("Select * from %s\n", param->table_name);
    condition_t *cond = param->condition_list;
    while (cond)
    {
        printf("\tCondition %s %d ", cond->col_name, cond->operator);
        print_value(cond->value);
        printf("\n");
        cond = cond->next;
    }
    */
    return go_SelectCallback(param);
}

extern int go_InsertIntoCallback(insert_into_t *param) __asm__ ("go.go_c_wrapper.InsertIntoCallback");
int insert_into_callback(insert_into_t *param)
{
    /*
    printf("Insert into %s\n", param->table_name);
    value_t *value = param->value_list;
    while (value)
    {
        printf("\tValue: ");
        print_value(value);
        printf("\n");
        value = value->next;
    }
    */
    return go_InsertIntoCallback(param);
}

extern int go_DeleteFromCallback(delete_from_t *param) __asm__ ("go.go_c_wrapper.DeleteFromCallback");
int delete_from_callback(delete_from_t *param)
{
    /*
    printf("Delete from %s\n", param->table_name);
    condition_t *cond = param->condition_list;
    while (cond)
    {
        printf("\tCondition: %s %d ", cond->col_name, cond->operator);
        print_value(cond->value);
        printf("\n");
        cond = cond->next;
    }
    */
    return go_DeleteFromCallback(param);
}

/* rewrite me to comminicate with Go */
int quit_callback(void *param)
{
    /* your finalizing codes */
    exit(0);
    return 0;
}

int execfile_callback(execfile_t *param)
{
    //printf("Execfile %s\n", param->filename);
    FILE *fin = fopen(param->filename, "r");
    if (!fin)
    {
        fprintf(stdout, "Error reading file: %s\n", param->filename);
        return 0;
    }
    YY_BUFFER_STATE bufs = interp_create_buffer (fin, YY_BUF_SIZE);
    interppush_buffer_state(bufs);
    interpparse();
    interppop_buffer_state();
    fclose(fin);
    return 0;
}
