package dest

type CustomType int32

const (
	A CustomType = iota
	B
)

type DC struct {
	DType CustomType `json:"a_type"`
}
