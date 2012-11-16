package recman

import (
	"../common"
	"encoding/json"
	"os"
)

func Insert(dbf *os.File, rec common.Record) error {
	enc := json.NewEncoder(dbf)
	enc.Encode(rec)
	return nil
}
