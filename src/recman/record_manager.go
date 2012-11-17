package recman

import (
	"../common"
	"encoding/json"
	"os"
)

func Insert(dbfp string, rec common.Record) (int64, error) {
	dbf, err := os.OpenFile(dbfp, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return -1, err
	}
	defer dbf.Close()
	offset := dbf.Stat().Size()
	enc := json.NewEncoder(dbf)
	enc.Encode(rec)
	return offset, nil
}

func DeleteAll(table_name string) error {
	return nil
}
