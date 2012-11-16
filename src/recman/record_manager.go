package recman

import (
	"../common"
	"encoding/json"
	"os"
)

func Insert(dbfp string, rec common.Record) (int, error) {
	dbf, err := os.OpenFile(dbfp, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return -1, err
	}
	defer dbf.Close()
	enc := json.NewEncoder(dbf)
	enc.Encode(rec)
	return 1, nil
}

func DeleteAll(table_name string) error {
	return nil
}
