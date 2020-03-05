package whatever

import (
	"time"

	"github.com/ekhabarov/sts/examples/nulls"
)

type (
	MyStruct struct {
		T              time.Time
		PT             *time.Time
		ThirdPartyType nulls.Time
	}
)
