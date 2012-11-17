package go_c_wrapper

import "fmt"
import "unsafe"
import "../common"
import "../core"

type __value_t struct { _type int32; value struct { int_t int32; }; next *__value_t; }
const _sizeof__value_t = 20
type _value_t struct { _type int32; value struct { int_t int32; }; next *__value_t; }
const _sizeof_value_t = 20
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
    fmt.Println("rte")
    return 0
}

func i__CreateTableCallback(param *_create_table_t) int {
    fmt.Println("rte")
    return 0
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
    fmt.Println(table.Name)
    core.CreateTable(&table)
    return 0
}

func DropTableCallback(param *int8) int{
    fmt.Println("asdasdasdasda")
    return 0
    core.DropTable(_CStringToString(param))
    return 0
}
