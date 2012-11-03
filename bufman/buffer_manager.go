package bufman

import (
	"fmt"
	"io/ioutil"
)

type page struct {
	status  int
	file    string
	content []byte
	pin     bool
}

// type info struct {
// 	path string
// 	page *page
// }

var (
	// fps [100]info
	buf [100]page
)

// func init_buffer(size int) (Buffer, error) {
//
// }

func BufWriteRecord(tab string, c []byte) error {

}

func BufWriteIndex(tab string, idx string, c []byte) error {

}

func BufReadIndex(tab string, idx string) ([]byte, error) {

}

func BufReadRecord(tab string) ([]byte, error) {

}

func BufCommit() error {
	err := nil
	for i, p := range buf {
		if p.status == 1 {
			err = ioutil.WriteFile(p.file, p.content, 0600)
		}
	}
	return err
}
