package field

type Right struct {
	A      int
	B      string
	C      float32 `json:"crc"`
	Double int     `db:"double"`
	WOO    int     `bar:"wo"`
	WOM    int     `bar:"wom,omitenmpty"`
}
