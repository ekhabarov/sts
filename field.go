package sts

import (
	"errors"
	"fmt"
	"go/types"
	"strings"
)

var (
	ErrEmptyLeftType  = errors.New("left type is empty")
	ErrEmptyRightType = errors.New("right type is empty")
	ErrEmptyLeftField = errors.New("left field is empty")
	ErrEmptRightField = errors.New("right field is empty")
	ErrEmptyTag       = errors.New("empty tag")
)

// fpair represents pair of matched fields.
type fpair struct {
	lt, rt string // types
	lf, rf string // field
	lp, rp bool   // pointers
	// fields in the pair can be assigned one onto another.
	assignable bool
	// fields in th pair require conversion.
	convertable bool
	// Order in structure.
	ord uint8
}

type fpairlist []fpair

// Len is the number of elements in the collection.
func (pl fpairlist) Len() int {
	return len(pl)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (pl fpairlist) Less(i int, j int) bool {
	return pl[i].ord < pl[j].ord
}

// Swap swaps the elements with indexes i and j.
func (pl fpairlist) Swap(i int, j int) {
	pl[i], pl[j] = pl[j], pl[i]
}

// String prints fields map with package names.
func (p fpair) Print(swap, debug bool, helperpkg string) (string, error) {
	lf, rf, lt, rt, lp, rp := p.lf, p.rf, p.lt, p.rt, p.lp, p.rp

	switch {
	case lf == "":
		return "", ErrEmptyLeftField
	case rf == "":
		return "", ErrEmptRightField
	case !p.assignable && lt == "":
		return "", ErrEmptyLeftType
	case !p.assignable && rt == "":
		return "", ErrEmptyRightType
	}

	if swap {
		lt, rt = rt, lt // types
		lf, rf = rf, lf // field names
		lp, rp = rp, lp // pointers
	}

	switch {
	// assignable has precedence over convertable
	case p.assignable:
		lf = "src." + lf
	case p.convertable:

		lf = fmt.Sprintf("%s(src.%s)", rt, lf)
	default:
		lt = strings.Replace(strings.Title(lt), ".", "", -1)
		rt = strings.Replace(strings.Title(rt), ".", "", -1)

		if lp {
			lt += "Ptr"
		}

		if rp {
			rt += "Ptr"
		}

		lf = fmt.Sprintf("%s2%s(src.%s)", lt, rt, lf)

		if helperpkg != "" {
			lf = helperpkg + "." + lf
		}
	}

	tpl := "%[1]s: %[2]s,"
	args := []interface{}{rf, lf}

	if debug {
		tpl += " \t// %[3]s: %[4]s\t\tassignable: %t, \t\tconvertable: %t"
		args = append(args, rt, lt, p.assignable, p.convertable)
	}

	return fmt.Sprintf(tpl, args...), nil
}

// match returns field name and type for right (destination) side found by
// next rules:
//
// *	Try to find destination field comparing tag from left side with field name
//		on the right.
// *  If field not found, try to compare tag from the left with vtags on the
//		right.
//
// Case when we have one structure which should be mapped with 3 different
// structures. In this case it's possible to define 3 different tags on source
// (left) structure which will be mapped to field name or valid tag name on
// destination (right) structure.
func match(tag string, right Fields, vtags []string) (string, string, error) {
	if tag == "" {
		return "", "", ErrEmptyTag
	}

	// search by field name
	if f, ok := right[tag]; ok {
		if f.Type == nil {
			return "", "", fmt.Errorf("type for field %q is not found", f)
		}
		return tag, typName(f.Type.String()), nil
	}

	// search among valid tags on the right side.
	for n, f := range right {
		for t, v := range f.Tags {
			for _, vt := range vtags {
				if vt == t && tag == v {
					return n, typName(f.Type.String()), nil
				}
			}
		}
	}

	return "", "", nil
}

// link pairs field from different sides by field name or by tag on destination
// (right) structure.
func link(
	left, right Fields,
	sourceTag string,
	vtags []string,
) (fpairlist, error) {
	fp := fpairlist{}

	// go through all of the source fields and try
	for name, f := range left {
		rf, rt, err := match(f.Tags[sourceTag], right, vtags)
		if err != nil && err != ErrEmptyTag { // skip fields without tags.
			return nil, err
		}

		if rf == "" {
			continue
		}

		ff := right[rf]
		fp = append(fp, fpair{
			lf:          name,
			rf:          rf,
			lt:          typName(f.Type.String()),
			rt:          rt,
			lp:          f.IsPointer,
			rp:          left[rf].IsPointer,
			assignable:  types.AssignableTo(f.Type, ff.Type),
			convertable: types.ConvertibleTo(f.Type, ff.Type),
			ord:         f.Ord,
		})
	}

	return fp, nil
}
