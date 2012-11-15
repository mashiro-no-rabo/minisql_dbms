package catman

import (
	"../common"
	"os"
)

func AllTables() ([]string, error) {
	dbs, err := os.OpenFile(common.DataDir, os.O_EXCL, 0700)
	if err != nil {
		return nil, err
	}
	names, err := dbs.Readdirnames(-1)
	if err != nil {
		common.ErrLogger.Println(err)
		return nil, err
	}
	return names, nil
}

func AddTable() error {
	return nil
}

func DelTable() error {
	return nil
}
