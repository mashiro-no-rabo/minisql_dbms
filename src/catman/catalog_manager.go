package catman

import (
	"../common"
	"encoding/json"
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
	idxs, err := os.OpenFile(common.DataDir+"/"+table_name+"/index", os.O_EXCL, 0700)
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

func TableInfo(table_name string) (*common.Table, error) {
	tab := new(common.Table)
	sf, err := os.OpenFile(common.DataDir+"/"+table_name+"/"+"schema.dbf", os.O_RDONLY, 0600)
	if err != nil {
		common.ErrLogger.Println(err)
		return tab, err
	}
	dec := json.NewDecoder(sf)

	err = dec.Decode(tab)
	if err != nil {
		return tab, err
	}
	return tab, nil
}
