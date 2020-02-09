package internal

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/recoilme/pudge"
)

type DB struct {
	Project *pudge.Db
	Sample  *pudge.Db
	Markup  *pudge.Db
}

func decodeBinary(raw []byte, makeValuePtr func() interface{}) interface{} {
	buf := bytes.NewBuffer(raw)
	dec := gob.NewDecoder(buf)
	key := makeValuePtr()
	_ = dec.Decode(key)
	fmt.Println("key: ", key)
	return key
}
