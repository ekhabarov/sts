// Code generated by sts v0.0.6. DO NOT EDIT.

package examples

func subsrc2anysub(src subsrc) anysub {
	return anysub{
		Test: src.Field1,
	}
}
func anysub2subsrc(src anysub) subsrc {
	return subsrc{
		Field1: src.Test,
	}
}
func subsrcPtr2anysubPtr(src *subsrc) *anysub {
	if src == nil {
		return nil
	}
	m := subsrc2anysub(*src)
	return &m
}
func anysubPtr2subsrcPtr(src *anysub) *subsrc {
	if src == nil {
		return nil
	}
	m := anysub2subsrc(*src)
	return &m
}
func subsrcList2anysubList(src []subsrc) []anysub {
	if src == nil {
		return nil
	}
	res := make([]anysub, len(src))
	for k, s := range src {
		p := subsrc2anysub(s)
		res[k] = p
	}
	return res
}
func anysubList2subsrcList(src []anysub) []subsrc {
	if src == nil {
		return nil
	}
	res := make([]subsrc, len(src))
	for k, s := range src {
		p := anysub2subsrc(s)
		res[k] = p
	}
	return res
}
func subsrcList2anysubPtrList(src []subsrc) []*anysub {
	if src == nil {
		return nil
	}
	res := make([]*anysub, len(src))
	for k, s := range src {
		p := subsrc2anysub(s)
		res[k] = &p
	}
	return res
}
func anysubPtrList2subsrcList(src []*anysub) []subsrc {
	if src == nil {
		return nil
	}
	res := make([]subsrc, len(src))
	for k, s := range src {
		p := anysub2subsrc(*s)
		res[k] = p
	}
	return res
}
func subsrcPtrList2anysubList(src []*subsrc) []anysub {
	if src == nil {
		return nil
	}
	res := make([]anysub, len(src))
	for k, s := range src {
		p := subsrc2anysub(*s)
		res[k] = p
	}
	return res
}
func anysubList2subsrcPtrList(src []anysub) []*subsrc {
	if src == nil {
		return nil
	}
	res := make([]*subsrc, len(src))
	for k, s := range src {
		p := anysub2subsrc(s)
		res[k] = &p
	}
	return res
}
func subsrcPtrList2anysubPtrList(src []*subsrc) []*anysub {
	if src == nil {
		return nil
	}
	res := make([]*anysub, len(src))
	for k, s := range src {
		p := subsrc2anysub(*s)
		res[k] = &p
	}
	return res
}
func anysubPtrList2subsrcPtrList(src []*anysub) []*subsrc {
	if src == nil {
		return nil
	}
	res := make([]*subsrc, len(src))
	for k, s := range src {
		p := anysub2subsrc(*s)
		res[k] = &p
	}
	return res
}
