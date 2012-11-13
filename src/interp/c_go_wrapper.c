/*
 * Source code for course work of Database System Design, 2012.09-11, Zhejiang University.
 * By:
 *   Shicheng XU(lightxuzju@gmail.com)
 *   Xun LOU(AquarHEAD@gmail.com)
 *   Pengyu CHEN(cpy.prefers.you@gmail.com)
 * COPYLEFT, ALL WRONGS RESERVED.
 */

#include <stdio.h>

#include "common.h"

int create_table_callback(create_table_t *param)
{
    printf("Create table: %s\n", param->name);
    column_t *col = param->columns;
    while (col)
    {
        printf("\tTable column: %s of type %d with attribute %s\n", col->name, 
            col->datatype->meta_datatype, col->attr == COL_ATTR_UNIQUE ? "unique" : "none");
        col = col->next;
    }
    printf("\tPrimary key: %s\n", param->primary_key);
    return 0;
}

int drop_table_callback(drop_table_t *param)
{
    printf("Drop table: %s\n", param->name);
    return 0;
}
