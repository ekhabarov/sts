package whatever

type (
	Embedded struct {
		CS string
	}

	MyStruct struct {
		I int
		S string
		Embedded
	}
)
