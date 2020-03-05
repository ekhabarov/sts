package input

type A002 struct {
	I int     `sts:"I" foo:"II" bar:"i"`
	S string  `sts:"s" foo:"Str" bar:"s"` // invalid value "s", yet
	F float32 `sts:"F" foo:"FF" bar:"f"`
}
