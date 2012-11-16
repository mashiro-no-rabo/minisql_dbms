package catman

import (
	"../common"
	"os"
)

func AllTables() ([]string, error) {
	dbs, err := os.OpenFile(common.DataDir, os.O_EXCL, 0700)
	if err != nil {
		common.ErrLogger.Println(err)
		return nil, err
	}
	names, err := dbs.Readdirnames(-1)
	if err != nil {
		common.ErrLogger.Println(err)
		return nil, err
	}
	return names, nil
}

func TableIndexes(table_name string) ([]string, error) {
	idxs, err := os.OpenFile(common.DataDir+"/"+table_name, os.O_EXCL, 0700)
	if err != nil {
		common.ErrLogger.Println(err)
		return nil, err
	}
	names, err := idxs.Readdirnames(-1)
	if err != nil {
		common.ErrLogger.Println(err)
		return nil, err
	}
	return names, nil
}
