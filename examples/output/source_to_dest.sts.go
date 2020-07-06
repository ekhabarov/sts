// Code generated by sts v0.0.6. DO NOT EDIT.

package output

import (
	"github.com/ekhabarov/sts/examples"
	"github.com/ekhabarov/sts/examples/dest"
)

func Source2Dest(src examples.Source) dest.Dest {
	return dest.Dest{
		I64:       int64(src.I1),
		B:         Int2Bool(src.I2),
		Nt:        TimeTimePtr2NullsTime(src.PT),
		JsonField: src.JJ,
		DB:        int64(src.D),
	}
}
func Dest2Source(src dest.Dest) examples.Source {
	return examples.Source{
		I1: int(src.I64),
		I2: Bool2Int(src.B),
		PT: NullsTime2TimeTimePtr(src.Nt),
		JJ: src.JsonField,
		D:  int32(src.DB),
	}
}
func SourcePtr2DestPtr(src *examples.Source) *dest.Dest {
	if src == nil {
		return nil
	}
	m := Source2Dest(*src)
	return &m
}
func DestPtr2SourcePtr(src *dest.Dest) *examples.Source {
	if src == nil {
		return nil
	}
	m := Dest2Source(*src)
	return &m
}
func SourceList2DestList(src []examples.Source) []dest.Dest {
	if src == nil {
		return nil
	}
	res := make([]dest.Dest, len(src))
	for k, s := range src {
		p := Source2Dest(s)
		res[k] = p
	}
	return res
}
func DestList2SourceList(src []dest.Dest) []examples.Source {
	if src == nil {
		return nil
	}
	res := make([]examples.Source, len(src))
	for k, s := range src {
		p := Dest2Source(s)
		res[k] = p
	}
	return res
}
func SourceList2DestPtrList(src []examples.Source) []*dest.Dest {
	if src == nil {
		return nil
	}
	res := make([]*dest.Dest, len(src))
	for k, s := range src {
		p := Source2Dest(s)
		res[k] = &p
	}
	return res
}
func DestPtrList2SourceList(src []*dest.Dest) []examples.Source {
	if src == nil {
		return nil
	}
	res := make([]examples.Source, len(src))
	for k, s := range src {
		p := Dest2Source(*s)
		res[k] = p
	}
	return res
}
func SourcePtrList2DestList(src []*examples.Source) []dest.Dest {
	if src == nil {
		return nil
	}
	res := make([]dest.Dest, len(src))
	for k, s := range src {
		p := Source2Dest(*s)
		res[k] = p
	}
	return res
}
func DestList2SourcePtrList(src []dest.Dest) []*examples.Source {
	if src == nil {
		return nil
	}
	res := make([]*examples.Source, len(src))
	for k, s := range src {
		p := Dest2Source(s)
		res[k] = &p
	}
	return res
}
func SourcePtrList2DestPtrList(src []*examples.Source) []*dest.Dest {
	if src == nil {
		return nil
	}
	res := make([]*dest.Dest, len(src))
	for k, s := range src {
		p := Source2Dest(*s)
		res[k] = &p
	}
	return res
}
func DestPtrList2SourcePtrList(src []*dest.Dest) []*examples.Source {
	if src == nil {
		return nil
	}
	res := make([]*examples.Source, len(src))
	for k, s := range src {
		p := Dest2Source(*s)
		res[k] = &p
	}
	return res
}
