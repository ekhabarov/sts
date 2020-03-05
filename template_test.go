package sts

import (
	. "github.com/onsi/ginkgo"
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

	Context("when header method called", func() {

		Context("swapped = false", func() {
			It("writes filled function header into output", func() {
				data.header(false)
				Expect(data.String()).To(Equal("func lt2rt(src lp.lt) rp.rt {\n"))
			})
		})

		Context("swapped = true", func() {
			It("writes filled function header into output", func() {
				data.header(true)
				Expect(data.String()).To(Equal("func rt2lt(src rp.rt) lp.lt {\n"))
			})
		})
	})

	Context("when retstmt method called", func() {

		Context("swapped = false", func() {
			It("writes filled function return statement into output", func() {
				data.retstmt(false)
				Expect(data.String()).To(Equal("return rp.rt {\n"))
			})
		})

		Context("swapped = true", func() {
			It("writes filled function header into output", func() {
				data.retstmt(true)
				Expect(data.String()).To(Equal("return lp.lt {\n"))
			})
		})
	})

	Context("when fieldmap method called", func() {

		Context("swapped = false", func() {
			It("writes filled fields into output", func() {
				// fieldmap/print method with debug=true tested in field_test.go
				err := data.fieldmap(false, false, "")
				Expect(err).NotTo(HaveOccurred())

				Expect(data.String()).To(Equal("Frf: src.Flf,\n"))
			})
		})

		Context("swapped = true", func() {
			It("writes filled function header into output", func() {
				err := data.fieldmap(true, false, "")
				Expect(err).NotTo(HaveOccurred())

				Expect(data.String()).To(Equal("Flf: src.Frf,\n"))
			})
		})
	})

	Context("when footer method called", func() {
		It("writes } into output", func() {
			data.footer()
			Expect(data.String()).To(Equal("}\n"))
		})

	})

})
