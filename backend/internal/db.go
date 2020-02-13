package internal

import (
	"bytes"
	"encoding/gob"

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

	return key
}
