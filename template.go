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

// ftpldata is a data for filling template for basic
// function.
type ftpldata struct {
	*bytes.Buffer
	lp, rp string // package
	lt, rt string // types
	fields fpairlist
}

func (d *ftpldata) Print(swap, debug bool, helperpkg string) error {
	d.header(swap)
	d.retstmt(swap)
	if err := d.fieldmap(swap, debug, helperpkg); err != nil {
		return err
	}
	d.footer()
	d.footer()

	return nil
}

func (d *ftpldata) header(swap bool) {
	lp, rp, lt, rt := d.lp, d.rp, d.lt, d.rt

	if rp != "" {
		rp += "."
	}

	if lp != "" {
		lp += "."
	}

	if swap {
		lp, rp = rp, lp
		lt, rt = rt, lt
	}

	fmt.Fprintf(d,
		// TODO(ekhabarov): Customizable 2 or To or whatever
		"func %[1]s2%[2]s(src %[3]s%[1]s) %[4]s%[2]s {\n",
		lt, rt, lp, rp,
	)
}

// reutrn
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

func (d *ftpldata) footer() {
	fmt.Fprintln(d, "}")
}
