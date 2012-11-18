package go_c_wrapper

import "fmt"
import "unsafe"
import "../common"
import "../core"
import "../catman"

type __value_t struct { value_type int32; int_t int32; str_t *int8; float_t float32; next *__value_t; }
const _sizeof__value_t = 28
type _value_t struct { value_type int32; int_t int32; str_t *int8; float_t float32; next *__value_t; }
const _sizeof_value_t = 28
type __datatype_t struct { meta_datatype int32; len int32; test int8; }
const _sizeof__datatype_t = 9
type _datatype_t struct { meta_datatype int32; len int32; test int8; }
const _sizeof_datatype_t = 9
type __column_t struct { name *int8; datatype *_datatype_t; attr int32; next *__column_t; }
const _sizeof__column_t = 28
type _column_t struct { name *int8; datatype *_datatype_t; attr int32; next *__column_t; }
const _sizeof_column_t = 28
type __create_table_t struct { name *int8; column_list *_column_t; primary_key *int8; }
const _sizeof__create_table_t = 24
type _create_table_t struct { name *int8; column_list *_column_t; primary_key *int8; }
const _sizeof_create_table_t = 24
type __drop_table_t struct { name *int8; }
const _sizeof__drop_table_t = 8
type _drop_table_t struct { name *int8; }
const _sizeof_drop_table_t = 8
type __create_index_t struct { index_name *int8; table_name *int8; col_name *int8; }
const _sizeof__create_index_t = 24
type _create_index_t struct { index_name *int8; table_name *int8; col_name *int8; }
const _sizeof_create_index_t = 24
type __drop_index_t struct { name *int8; }
const _sizeof__drop_index_t = 8
type _drop_index_t struct { name *int8; }
const _sizeof_drop_index_t = 8
type __condition_t struct { col_name *int8; operator int32; value *_value_t; next *__condition_t; }
const _sizeof__condition_t = 28
type _condition_t struct { col_name *int8; operator int32; value *_value_t; next *__condition_t; }
const _sizeof_condition_t = 28
type __select_t struct { table_name *int8; condition_list *_condition_t; }
const _sizeof__select_t = 16
type _select_t struct { table_name *int8; condition_list *_condition_t; }
const _sizeof_select_t = 16
type __insert_into_t struct { table_name *int8; value_list *_value_t; }
const _sizeof__insert_into_t = 16
type _insert_into_t struct { table_name *int8; value_list *_value_t; }
const _sizeof_insert_into_t = 16
type _delete_from_t struct { table_name *int8; condition_list *_condition_t; }
const _sizeof_delete_from_t = 16
type __execfile_t struct { filename *int8; }
const _sizeof__execfile_t = 8
type _execfile_t struct { filename *int8; }
const _sizeof_execfile_t = 8

func _CStringToString(pstr *int8) string{
    var bytearray []byte
    str := (*byte)(unsafe.Pointer(pstr))
    for ; *str != 0; str = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(str)) +
        unsafe.Sizeof(byte))) {
        bytearray = append(bytearray, *str)
    }
    return string(bytearray)
}

func CreateTableCallback(param *_create_table_t) int {
    var table common.Table
    table.Name = _CStringToString(param.name)
    for col := param.column_list; col != nil; col = (*_column_t)(col.next) {
        var column common.Column
        column.Name = _CStringToString(param.primary_key)
        column.Type = int(col.datatype.meta_datatype)
        column.Unique = col.attr == 1
        column.Length = int64(col.datatype.len)
        table.Columns = append(table.Columns, column)
    }
    pkey := _CStringToString(param.primary_key)
    table.PKey = -1
    for i, col := range table.Columns {
        if pkey == col.Name {
            table.PKey = i
        }
    }
    if table.PKey == -1 {
        fmt.Println("Invalid primary key: %s", pkey)
        return 1
    }
    core.CreateTable(&table)
    return 0
}

func DropTableCallback(param *_drop_table_t) int{
    core.DropTable(_CStringToString(param.name))
    return 0
}

func CreateIndexCallback(param *_create_index_t) int{
    table_name := _CStringToString(param.table_name)
    index_name := _CStringToString(param.index_name)
    index_key := -1
    iname := _CStringToString(param.col_name)
    table, err := catman.TableInfo(table_name)
    if err != nil {
        fmt.Println("Invalid table name: %s", table_name)
        return 1
    }
    for i, col := range table.Columns {
        if iname == col.Name {
            index_key = i
        }
    }
    if index_key == -1 {
        fmt.Println("Invalid index key: %s", iname)
        return 1
    }
    core.CreateIndex(table_name, index_name, index_key)
    return 0
}

func DropIndexCallback(param *_drop_index_t) int{
    iname := _CStringToString(param.name)
    tname := ""
    ikey := -1
    tablist, err := catman.AllTables()
    if err != nil {
        fmt.Println("Error getting table list.")
        return 1
    }
    for _, tablename := range tablist {
        table, _ := catman.TableInfo(tablename)
        for i, col := range table.Columns {
            if iname == col.Name {
                tname = table.Name
                ikey = i
            }
        }
    }
    if ikey == -1 {
         fmt.Println("Invalid index key: %s", iname)
        return 1
    }
    core.DropIndex(tname, iname)
    return 0
}

func SelectCallback(param *_select_t) int{
    table_name := _CStringToString(param.table_name)
    var cond []common.Condition
    for pcond := param.condition_list; pcond != nil; pcond =
       (*_condition_t)(pcond.next) {
        var tmpcond common.Condition
        tmpcond.ColName = _CStringToString(pcond.col_name)
        tmpcond.Op = int(pcond.operator)
        tmpcond.ValueType = int(pcond.value.value_type)
        switch tmpcond.ValueType{
        case common.IntCol:
            tmpcond.ValueInt = (int)(pcond.value.int_t)
            break
        case common.StrCol:
            tmpcond.ValueString = _CStringToString(pcond.value.str_t)
            break
        case common.FltCol:
            tmpcond.ValueFloat = (float64)(pcond.value.float_t)
            break
        default:
            break
        }
        cond = append(cond, tmpcond)
    }
    core.Select(table_name, cond)
    return 0
}

func InsertIntoCallback(param *_insert_into_t) int{
    table_name := _CStringToString(param.table_name)
    var vals []common.CellValue
    for v := param.value_list; v != nil; v = (*_value_t)(v.next) {
        switch v.value_type {
        case common.IntCol:
            vals = append(vals, (common.IntVal)(v.int_t))
            break
        case common.StrCol:
            vals = append(vals, (common.StrVal)(_CStringToString(v.str_t)))
            break
        case common.FltCol:
            vals = append(vals, (common.FltVal)(v.float_t))
            break
        default:
            break
        }
    } 
    core.Insert(table_name, vals)
    return 0
}

func DeleteFromCallback(param *_delete_from_t) int{
    table_name := _CStringToString(param.table_name)
    var cond []common.Condition
    for pcond := param.condition_list; pcond != nil; pcond =
       (*_condition_t)(pcond.next) {
        var tmpcond common.Condition
        tmpcond.ColName = _CStringToString(pcond.col_name)
        tmpcond.Op = int(pcond.operator)
        tmpcond.ValueType = int(pcond.value.value_type)
        switch tmpcond.ValueType{
        case common.IntCol:
            tmpcond.ValueInt = (int)(pcond.value.int_t)
            break
        case common.StrCol:
            tmpcond.ValueString = _CStringToString(pcond.value.str_t)
            break
        case common.FltCol:
            tmpcond.ValueFloat = (float64)(pcond.value.float_t)
            break
        default:
            break
        }
        cond = append(cond, tmpcond)
    }
    core.Delete(table_name, cond)
    return 0
}

