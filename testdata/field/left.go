package field

type Left struct {
	A  int     `sts:"A"`
	B  string  `sts:"B"`
	C  float32 `sts:"crc"`
	D  int     `bar:"double"`
	O  int     `foo:"wo,omitempty"`
	OM int     `foo:"wom"`
}
