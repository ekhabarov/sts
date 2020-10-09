package whatever

import (
	"time"

	"github.com/powerflyco/sts/examples/nulls"
)

type (
	TimeTime struct {
		T              time.Time
		PT             *time.Time
		ThirdPartyType nulls.Time
	}
)
