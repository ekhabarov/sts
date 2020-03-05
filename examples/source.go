package examples

import (
	"time"

	"github.com/ekhabarov/sts/examples/nulls"
)

//go:generate sts -src $GOFILE:Source -dst $GOFILE:Dest -o ./output -dt json,db
type Source struct {
	I  int
	S  string
	I1 int        `sts:"I64"`
	I2 int        `sts:"B"`  // types.Basic
	PT *time.Time `sts:"Nt"` // types.Named
	JJ string     `sts:"json_field"`
	D  int32      `sts:"db_field"`
}

// dest, techincally it's some isolated structure we cannot change.
type Dest struct {
	I         int
	S         string
	I64       int64
	B         bool
	Nt        nulls.Time
	JsonField string `json:"json_field"`
	DB        int64  `db:"db_field"`
}

//go:generate sts -src $GOFILE:subsrc -dst $GOFILE:anysub -o .
type subsrc struct {
	Field1 string `sts:"Test"`
}

type anysub struct {
	Test string
}
