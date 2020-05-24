package sts

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Template", func() {
	var (
		data ftpldata
	)

	BeforeEach(func() {
		data = newTemplate("lt", "rt", "lp", "rp", []fpair{
			{
				lt:          "flt",
				rt:          "frt",
				lf:          "Flf",
				rf:          "Frf",
				lp:          true,
				rp:          true,
				assignable:  true,
				convertable: true,
				ord:         1,
			},
		})
	})

	Context("when newftpldata called", func() {

		It("returns initialized ftpldata object", func() {

			Expect(data.Buffer).NotTo(BeNil())

			Expect(data.lt).To(Equal("lt"))
			Expect(data.rt).To(Equal("rt"))
			Expect(data.lp).To(Equal("lp"))
			Expect(data.rp).To(Equal("rp"))

			f := data.fields[0]

			Expect(f.lt).To(Equal("flt"))
			Expect(f.rt).To(Equal("frt"))
			Expect(f.lf).To(Equal("Flf"))
			Expect(f.rf).To(Equal("Frf"))
			Expect(f.lp).To(BeTrue())
			Expect(f.rp).To(BeTrue())
			Expect(f.convertable).To(BeTrue())
			Expect(f.assignable).To(BeTrue())
			Expect(f.ord).To(Equal(uint8(1)))
		})
	})

	DescribeTable("Header func",
		func(swap, lpp, rpp, list bool, expected string) {
			data.header(swap, lpp, rpp, list)
			Expect(data.String()).To(Equal(expected))
		},

		Entry("", false, false, false, false, "func lt2rt(src lp.lt) rp.rt {\n"),
		Entry("", true, false, false, false, "func rt2lt(src rp.rt) lp.lt {\n"),

		Entry("", false, true, false, false, "func ltPtr2rt(src *lp.lt) rp.rt {\n"),
		Entry("AA", true, true, false, false, "func rt2ltPtr(src rp.rt) *lp.lt {\n"),

		Entry("", false, false, true, false, "func lt2rtPtr(src lp.lt) *rp.rt {\n"),
		Entry("", true, false, true, false, "func rtPtr2lt(src *rp.rt) lp.lt {\n"),

		Entry("", false, true, true, false, "func ltPtr2rtPtr(src *lp.lt) *rp.rt {\n"),
		Entry("", true, true, true, false, "func rtPtr2ltPtr(src *rp.rt) *lp.lt {\n"),

		Entry("", false, false, false, true, "func ltList2rtList(src []lp.lt) []rp.rt {\n"),
		Entry("", true, false, false, true, "func rtList2ltList(src []rp.rt) []lp.lt {\n"),

		Entry("", false, false, true, true, "func ltList2rtPtrList(src []lp.lt) []*rp.rt {\n"),
		Entry("", true, false, true, true, "func rtPtrList2ltList(src []*rp.rt) []lp.lt {\n"),

		Entry("", false, true, true, true, "func ltPtrList2rtPtrList(src []*lp.lt) []*rp.rt {\n"),
		Entry("", true, true, true, true, "func rtPtrList2ltPtrList(src []*rp.rt) []*lp.lt {\n"),
	)

	DescribeTable("retstmt func",
		func(swap bool, expected string) {
			data.retstmt(swap)
			Expect(data.String()).To(Equal(expected))
		},

		Entry("swapped == false", false, "return rp.rt {\n"),
		Entry("swapped == true", true, "return lp.lt {\n"),
	)

	DescribeTable("Fieldmap func",
		func(swap bool, hp, expected string) {
			err := data.fieldmap(swap, false, "")
			Expect(err).NotTo(HaveOccurred())

			Expect(data.String()).To(Equal(expected))
		},

		// fieldmap/print method with debug=true tested in field_test.go
		Entry("swapped == false", false, "", "Frf: src.Flf,\n"),
		Entry("swapped == true", true, "", "Flf: src.Frf,\n"),
	)

	Context("when footer method called", func() {
		It("writes } into output", func() {
			data.footer()
			Expect(data.String()).To(Equal("}\n"))
		})
	})

	DescribeTable("fn func",
		func(lp, rp, list bool, expected string) {
			n := fn("left", "right", lp, rp, list)
			Expect(n).To(Equal(expected))
		},

		Entry("Non-pointers", false, false, false, "left2right"),
		Entry("Left pointer", true, false, false, "leftPtr2right"),
		Entry("Right pointer", false, true, false, "left2rightPtr"),
		Entry("Both pointers", true, true, false, "leftPtr2rightPtr"),

		Entry("Non-pointers: list", false, false, true, "leftList2rightList"),
		Entry("Left pointer: list", true, false, true, "leftPtrList2rightList"),
		Entry("Right pointer: list", false, true, true, "leftList2rightPtrList"),
		Entry("Both pointers: list", true, true, true, "leftPtrList2rightPtrList"),
	)

	Describe("funcName", func() {

		Context("when two args passed", func() {
			It("return name with default divider", func() {
				name := funcName("left", "right")
				Expect(name).To(Equal("left2right"))
			})
		})

		Context("when three args passed", func() {
			It("return name with divider equals to third arg", func() {
				name := funcName("left", "right", "To")
				Expect(name).To(Equal("leftToright"))
			})
		})
	})

	DescribeTable("ifSoureNil",
		func(p bool, expected string) {
			i := ifSoureNil(p)
			Expect(i).To(Equal(expected))
		},

		Entry("Non-pointer", false, ""),
		Entry("Pointer", true, "if src == nil { return nil }"),
	)

	DescribeTable("ptrLoopBody",
		func(lp, rp bool, expected string) {
			b := ptrLoopBody(lp, rp, "name")
			Expect(b).To(Equal(expected))
		},

		Entry("src => dst", false, false, "p := name(s); res[k] = p"),
		Entry("src => *dst", false, true, "p := name(s); res[k] = &p"),
		Entry("*src => dst", true, false, "p := name(*s); res[k] = p"),
		Entry("*src => *dst", true, true, "p := name(*s); res[k] = &p"),
	)

	DescribeTable("typName",
		func(pn, t string, p, l bool, expected string) {
			n := typName(pn, t, p, l)
			Expect(n).To(Equal(expected))
		},

		Entry("name", "", "type", false, false, "type"),
		Entry("pointer", "", "type", true, false, "*type"),
		Entry("list", "", "type", false, true, "[]type"),
		Entry("list of pointers", "", "type", true, true, "[]*type"),

		Entry("name", "package", "type", false, false, "package.type"),
		Entry("pointer", "package", "type", true, false, "*package.type"),
		Entry("list", "package", "type", false, true, "[]package.type"),
		Entry("list of pointers", "package", "type", true, true, "[]*package.type"),
	)

	DescribeTable("PrintWithBothPointers func",
		func(swap bool, expected string) {
			data.PrintWithBothPointers(swap)
			Expect(data.String()).To(Equal(expected))
		},

		Entry("swapped == false", false, `func ltPtr2rtPtr(src *lp.lt) *rp.rt {
if src == nil { return nil }; m := lt2rt(*src); return &m}
`),
		Entry("swapped == true", true, `func rtPtr2ltPtr(src *rp.rt) *lp.lt {
if src == nil { return nil }; m := rt2lt(*src); return &m}
`),
	)

	DescribeTable("PrintList",
		func(swap, lp, rp bool, expected string) {
			data.PrintList(swap, lp, rp)
			Expect(data.String()).To(Equal(expected))

		},
		Entry("src => dst", false, false, false, `func ltList2rtList(src []lp.lt) []rp.rt {
if src == nil { return nil }; res := make([]rp.rt, len(src)); for k, s := range src { p := lt2rt(s); res[k] = p }; return res}
`),

		Entry("src => *dst", false, false, true, `func ltList2rtPtrList(src []lp.lt) []*rp.rt {
if src == nil { return nil }; res := make([]*rp.rt, len(src)); for k, s := range src { p := lt2rt(s); res[k] = &p }; return res}
`),

		Entry("*src => dst", false, true, false, `func ltPtrList2rtList(src []*lp.lt) []rp.rt {
if src == nil { return nil }; res := make([]rp.rt, len(src)); for k, s := range src { p := lt2rt(*s); res[k] = p }; return res}
`),

		Entry("*src => *dst", false, true, true, `func ltPtrList2rtPtrList(src []*lp.lt) []*rp.rt {
if src == nil { return nil }; res := make([]*rp.rt, len(src)); for k, s := range src { p := lt2rt(*s); res[k] = &p }; return res}
`),
	)

})
