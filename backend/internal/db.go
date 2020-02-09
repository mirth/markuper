package internal

import "github.com/recoilme/pudge"

type DB struct {
	Project *pudge.Db
	Sample  *pudge.Db
	Markup  *pudge.Db
}
