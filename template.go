package sts

import (
	"bytes"
	"fmt"
	"go/ast"
	"sort"
)

func newTemplate(lt, rt, lp, rp string, ff []fpair) ftpldata {
	return ftpldata{
		Buffer: &bytes.Buffer{},
		lt:     lt,
		rt:     rt,
		lp:     lp,
		rp:     rp,
		fields: ff,
	}
}

// ftpldata is a data for filling template for basic function.
type ftpldata struct {
	*bytes.Buffer
	lp, rp string // package
	lt, rt string // types
	fields fpairlist
}

// Print adds main conversion functions into buffer, which contains matched
// fields.
func (d *ftpldata) Print(swap, debug bool, helperpkg string) error {
	d.header(swap, false, false, false)
	d.retstmt(swap)
	if err := d.fieldmap(swap, debug, helperpkg); err != nil {
		return err
	}
	d.footer()
	d.footer()

	return nil
}

func (d *ftpldata) fieldmap(swap, debug bool, helperpkg string) error {
	sort.Sort(d.fields)

	for _, f := range d.fields {
		if !ast.IsExported(f.lf) || !ast.IsExported(f.rf) {
			continue
		}
		p, err := f.Print(swap, debug, helperpkg)
		if err != nil {
			return err
		}

		fmt.Fprintln(d, p)
	}

	return nil
}

// PrintWithPointer prints functions for pointer-to-pointer conversion.
func (d *ftpldata) PrintWithBothPointers(swap bool) {
	d.header(swap, true, true, false)

	lt, rt := d.lt, d.rt

	if swap {
		lt, rt = rt, lt
	}

	fmt.Fprintf(d, `%s; m := %s(*src); return %sm`,
		ifSoureNil(true),
		funcName(lt, rt),
		"&",
	)

	d.footer()
}

func ifSoureNil(p bool) string {
	if p {
		return "if src == nil { return nil }"
	}
	return ""
}

func fn(l, r string, lp, rp, list bool) string {
	sep := "2"

	if lp {
		l += "Ptr"
	}

	if rp {
		r += "Ptr"
	}

	if list {
		l += "List"
		r += "List"
	}

	return fmt.Sprintf("%s%s%s", l, sep, r)
}

func funcName(left, right string, rest ...string) string {
	to := "2"
	if len(rest) > 0 {
		to = rest[0]
	}
	return fmt.Sprintf("%s%s%s", left, to, right)
}

func (d *ftpldata) PrintList(swap, lptr, rptr bool) {
	d.header(swap, lptr, rptr, true)

	lp, rp, lt, rt := d.lp, d.rp, d.lt, d.rt
	ltf, rtf := lt, rt

	if swap {
		lt, rt = rt, lt
		ltf, rtf = rtf, ltf
		lp, rp = rp, lp
		lptr, rptr = rptr, lptr
	}

	fname := funcName(ltf, rtf)

	fmt.Fprintf(d,
		`%s; res := make(%s, len(src)); for k, s := range src { %s }; return res`,
		ifSoureNil(true),
		typName(rp, rt, rptr, true),
		ptrLoopBody(lptr, rptr, fname),
	)
	d.footer()
}

// ttypName returns type name for template.
func typName(packageName, typ string, ptr, list bool) string {
	b, s, dot := "", "", ""

	if ptr {
		s = "*"
	}

	if list {
		b = "[]"
	}

	if packageName != "" {
		dot = "."
	}

	return fmt.Sprintf("%s%s%s%s%s", b, s, packageName, dot, typ)
}

// pptrLoopBody builds for loop body for List functions.
func ptrLoopBody(lp, rp bool, name string) string {
	t := "p := %s(%ss); res[k] = %sp"

	s, p := "", ""

	if lp {
		s = "*"
	}
	if rp {
		p = "&"
	}

	return fmt.Sprintf(t, name, s, p)
}

func (d *ftpldata) header(swap, lptr, rptr, list bool) {
	lp, rp, lt, rt := d.lp, d.rp, d.lt, d.rt

	if swap {
		lp, rp = rp, lp
		lt, rt = rt, lt
		lptr, rptr = rptr, lptr
	}

	fmt.Fprintf(d,
		"func %s(src %s) %s {\n",
		fn(lt, rt, lptr, rptr, list),
		typName(lp, lt, lptr, list),
		typName(rp, rt, rptr, list),
	)
}

// retstmt prints return statement into buffer.
func (d *ftpldata) retstmt(swap bool) {
	rt, rp := d.rt, d.rp

	if swap {
		rt, rp = d.lt, d.lp
	}

	if rp != "" {
		rp += "."
	}

	fmt.Fprintf(d, "return %s%s {\n", rp, rt)
}

func (d *ftpldata) footer() {
	fmt.Fprintln(d, "}")
}
