package source

type CustomType int

const (
	A CustomType = iota
	B
)

type AC struct {
	AType CustomType `json:"a_type"`
}
