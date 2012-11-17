package go_c_wrapper

import "unsafe"
import "../common"
import "../core"

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
    table.Columns = make(map[string]common.Column)
    for col := param.column_list; col != nil; col = (*_column_t)(col.next) {
        var column common.Column
        column.Type = int(col.datatype.meta_datatype)
        column.Unique = col.attr == 1
        column.Length = int64(col.datatype.len)
        table.Columns[_CStringToString(col.name)] = column
    }
    table.PKey = _CStringToString(param.primary_key)
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
    index_key := _CStringToString(param.col_name)
    core.CreateIndex(table_name, index_name, index_key)
    return 0
}

func DropIndexCallback(param *_drop_index_t) int{
    core.DropIndex(_CStringToString(param.name))
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
    var rec common.Record
    core.Insert(table_name, rec)
    return 0
}
