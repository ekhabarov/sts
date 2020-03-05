package output

import (
	"time"

	"github.com/ekhabarov/sts/examples/nulls"
)

func Bool2Int(b bool) int {
	if b {
		return 1
	}

	return 0
}

func Int2Bool(i int) bool {
	return i > 0
}

func TimeTime2NullsTime(t time.Time) nulls.Time {
	return nulls.Time{Time: t}
}

func TimeTimePtr2NullsTime(t *time.Time) nulls.Time {
	if t == nil {
		return nulls.Time{}
	}
	return nulls.Time{
		Time:  *t,
		Valid: true,
	}
}

func NullsTime2TimeTime(nt nulls.Time) time.Time {
	if nt.Valid {
		return nt.Time
	}
	return time.Time{}
}

func NullsTime2TimeTimePtr(nt nulls.Time) *time.Time {
	if nt.Valid {
		return &nt.Time
	}
	return &time.Time{}
}
