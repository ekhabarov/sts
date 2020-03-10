package examples

import (
	"time"
)

//go:generate sts -src $GOFILE:Source -dst ./dest/dest.go:Dest -o ./output -dt json,db
type Source struct {
	I  int
	S  string
	I1 int        `sts:"I64"`
	I2 int        `sts:"B"`  // types.Basic
	PT *time.Time `sts:"Nt"` // types.Named
	JJ string     `sts:"json_field"`
	D  int32      `sts:"db_field"`
	R  Doer
}

//go:generate sts -src $GOFILE:subsrc -dst $GOFILE:anysub -o .
type subsrc struct {
	Field1 string `sts:"Test"`
}

type anysub struct {
	Test string
}
