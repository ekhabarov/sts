package sts

import (
	"bytes"
	"fmt"
	"go/ast"
	"sort"
	"strings"
)

func newTemplate(lt, rt, lpkg, rpkg string, ff []fpair) ftpldata {
	return ftpldata{
		Buffer: &bytes.Buffer{},
		lt:     lt,
		rt:     rt,
		lpkg:   lpkg,
		rpkg:   rpkg,
		fields: ff,
	}
}

// ftpldata is a data for filling template for basic function.
type ftpldata struct {
	*bytes.Buffer
	lpkg, rpkg string // package
	lt, rt     string // types
	fields     fpairlist
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
	lpkg, rpkg := d.lpkg, d.rpkg

	if swap {
		lt, rt = rt, lt
		lpkg, rpkg = rpkg, lpkg
	}

	fmt.Fprintf(d, `%s; m := %s(*src); return %sm`,
		ifSoureNil(true),
		funcName(lt, rt, lpkg, rpkg),
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

func fnHeader(l, r, lpkg, rpkg string, lp, rp, list bool) string {
	sep := "2"

	if l == r {
		l = strings.ToUpper(lpkg[:1]) + lpkg[1:] + l
		r = strings.ToUpper(rpkg[:1]) + rpkg[1:] + r
	}

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

func funcName(l, r, lpkg, rpkg string, rest ...string) string {
	to := "2"
	if len(rest) > 0 {
		to = rest[0]
	}

	if l == r {
		l = strings.ToUpper(lpkg[:1]) + lpkg[1:] + l
		r = strings.ToUpper(rpkg[:1]) + rpkg[1:] + r
	}

	return fmt.Sprintf("%s%s%s", l, to, r)
}

func (d *ftpldata) PrintList(swap, lptr, rptr bool) {
	d.header(swap, lptr, rptr, true)

	lpkg, rpkg, lt, rt := d.lpkg, d.rpkg, d.lt, d.rt
	ltf, rtf := lt, rt

	if swap {
		lt, rt = rt, lt
		ltf, rtf = rtf, ltf
		lpkg, rpkg = rpkg, lpkg
		lptr, rptr = rptr, lptr
	}

	fname := funcName(ltf, rtf, lpkg, rpkg)

	fmt.Fprintf(d,
		`%s; res := make(%s, len(src)); for k, s := range src { %s }; return res`,
		ifSoureNil(true),
		typName(rpkg, rt, rptr, true),
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
	lpkg, rpkg, lt, rt := d.lpkg, d.rpkg, d.lt, d.rt

	if swap {
		lpkg, rpkg = rpkg, lpkg
		lt, rt = rt, lt
		lptr, rptr = rptr, lptr
	}

	fmt.Fprintf(d,
		"func %s(src %s) %s {\n",
		fnHeader(lt, rt, lpkg, rpkg, lptr, rptr, list),
		typName(lpkg, lt, lptr, list),
		typName(rpkg, rt, rptr, list),
	)
}

// retstmt prints return statement into buffer.
func (d *ftpldata) retstmt(swap bool) {
	rt, rpkg := d.rt, d.rpkg

	if swap {
		rt, rpkg = d.lt, d.lpkg
	}

	if rpkg != "" {
		rpkg += "."
	}

	fmt.Fprintf(d, "return %s%s {\n", rpkg, rt)
}

func (d *ftpldata) footer() {
	fmt.Fprintln(d, "}")
}
