package whatever

type (
	Embedded struct {
		CS string
	}

	WithEmbedded struct {
		I int
		S string
		Embedded
	}
)
