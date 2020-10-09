package dest

import "github.com/powerflyco/sts/examples/nulls"

// Dest, techincally it's some isolated structure we cannot change.
type Dest struct {
	I         int
	S         string
	I64       int64
	B         bool
	Nt        nulls.Time
	JsonField string `json:"json_field"`
	DB        int64  `db:"db_field"`
}
