package sts

import (
	"sort"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Field", func() {

	Describe("Pairlist", func() {

		Context("when call Len method", func() {

			It("returns length of the list", func() {
				fl := fpairlist{fpair{}, fpair{}, fpair{}}
				Expect(fl.Len()).To(Equal(3))
			})
		})

		Context("when call Less method", func() {

			It("compares slice elements", func() {
				fl := fpairlist{fpair{ord: 1}, fpair{ord: 2}}
				Expect(fl.Less(0, 1)).To(BeTrue())
				Expect(fl.Less(1, 0)).To(BeFalse())
			})
		})

		Context("when call Swap method", func() {

			It("swaps elements", func() {
				fl := fpairlist{fpair{ord: 1}, fpair{ord: 2}}
				Expect(fl[0].ord).To(Equal(uint8(1)))
				Expect(fl[1].ord).To(Equal(uint8(2)))
				fl.Swap(0, 1)
				Expect(fl[0].ord).To(Equal(uint8(2)))
				Expect(fl[1].ord).To(Equal(uint8(1)))
			})
		})
	})

	Describe("Field printer", func() {

		var (
			ar  = "right: src.left,"
			ars = "left: src.right,"

			cr  = "right: Rtype(src.left),"
			crs = "left: Ltype(src.right),"

			lf, rf = "left", "right"
			lt, rt = "Ltype", "Rtype"
		)

		DescribeTable("Print",
			func(p fpair, hp, exp, expSwapped string, experr error) {
				got, err := p.Print(false, false, hp)
				if experr == nil {
					Expect(err).NotTo(HaveOccurred())
				} else {
					Expect(err).To(MatchError(experr.Error()))
				}

				gotSwapped, err := p.Print(true, false, hp)
				if experr == nil {
					Expect(err).NotTo(HaveOccurred())
				} else {
					Expect(err).To(MatchError(experr.Error()))
				}

				Expect(got).To(Equal(exp))
				Expect(gotSwapped).To(Equal(expSwapped))
			},

			Entry("Empty struct", fpair{}, "", "", "", ErrEmptyLeftField),

			Entry("Empty right field", fpair{lf: lf}, "", "", "", ErrEmptRightField),

			// assignable && !convertable
			Entry("Field names only, assignable", fpair{
				lf: lf, rf: rf,
				assignable: true,
			}, "", ar, ars, nil),

			Entry("Field names with types, assignable", fpair{
				lf: lf, rf: rf, lt: lt, rt: rt,
				assignable: true,
			}, "", ar, ars, nil),

			Entry("Field names with types, left pointer, assignable", fpair{
				lf: lf, rf: rf, lt: lt, rt: rt,
				lp: true, rp: false,
				assignable: true,
			}, "", ar, ars, nil),

			Entry("Field names with types, right pointer, assignable", fpair{
				lf: lf, rf: rf, lt: lt, rt: rt,
				lp: false, rp: true,
				assignable: true,
			}, "", ar, ars, nil),

			Entry("Field names with types, both pointers, assignable", fpair{
				lf: lf, rf: rf, lt: lt, rt: rt,
				lp: true, rp: true,
				assignable: true,
			}, "", ar, ars, nil),

			// assignable && convertable
			Entry("Field names only, assignable & convertable", fpair{
				lf: lf, rf: rf, assignable: true,
				convertable: true,
			}, "", ar, ars, nil),

			Entry("Field names with types, assignable & convertable", fpair{
				lf: lf, rf: rf, lt: lt, rt: rt,
				assignable:  true,
				convertable: true,
			}, "", ar, ars, nil),

			Entry("Field names with types, left pointer, assignable & convertable",
				fpair{
					lf: lf, rf: rf, lt: lt, rt: rt,
					lp: true, rp: false,
					assignable:  true,
					convertable: true,
				}, "", ar, ars, nil),

			Entry("Field names with types, right pointer, assignable & convertable",
				fpair{
					lf: lf, rf: rf, lt: lt, rt: rt,
					lp: false, rp: true,
					assignable:  true,
					convertable: true,
				}, "", ar, ars, nil),

			Entry("Field names with types, both pointers, assignable & convertable",
				fpair{
					lf: lf, rf: rf, lt: lt, rt: rt,
					lp: true, rp: true,
					assignable:  true,
					convertable: true,
				}, "", ar, ars, nil),

			// !assignable && convertable
			Entry("Field names only, convertable", fpair{
				lf: lf, rf: rf,
				convertable: true,
			}, "", "", "", ErrEmptyLeftType),

			Entry("Field names only, empty right type, convertable", fpair{
				lf: lf, rf: rf, lt: lt,
				convertable: true,
			}, "", "", "", ErrEmptyRightType),

			Entry("Field names with types, convertable", fpair{
				lf: lf, rf: rf, lt: lt, rt: rt,
				convertable: true,
			}, "", cr, crs, nil),

			Entry("Field names with types, left pointer, convertable", fpair{
				lf: lf, rf: rf, lt: lt, rt: rt,
				lp: true, rp: false,
				convertable: true,
			}, "", cr, crs, nil),

			Entry("Field names with types, right pointer, convertable", fpair{
				lf: lf, rf: rf, lt: lt, rt: rt,
				lp: false, rp: true,
				convertable: true,
			}, "", cr, crs, nil),

			Entry("Field names with types, both pointers, convertable", fpair{
				lf: lf, rf: rf, lt: lt, rt: rt,
				lp: true, rp: true,
				convertable: true,
			}, "", cr, crs, nil),

			// !assignable && !convertable
			Entry("Field names only", fpair{
				lf: lf, rf: rf,
			}, "", "", "", ErrEmptyLeftType),

			Entry("Field names with types", fpair{
				lf: lf, rf: rf, lt: lt, rt: rt,
			}, "",
				"right: Ltype2Rtype(src.left),",
				"left: Rtype2Ltype(src.right),",
				nil,
			),

			Entry("Field names with types, left pointer", fpair{
				lf: lf, rf: rf, lt: lt, rt: rt,
				lp: true, rp: false,
			}, "",
				"right: LtypePtr2Rtype(src.left),",
				"left: Rtype2LtypePtr(src.right),",
				nil,
			),

			Entry("Field names with types, right pointer", fpair{
				lf: lf, rf: rf, lt: lt, rt: rt,
				lp: false, rp: true,
			}, "",
				"right: Ltype2RtypePtr(src.left),",
				"left: RtypePtr2Ltype(src.right),",
				nil,
			),

			Entry("Field names with types, both pointers", fpair{
				lf: lf, rf: rf, lt: lt, rt: rt,
				lp: true, rp: true,
			}, "",
				"right: LtypePtr2RtypePtr(src.left),",
				"left: RtypePtr2LtypePtr(src.right),",
				nil,
			),

			Entry("Field names with types & helpers in different package", fpair{
				lf: lf, rf: rf, lt: lt, rt: rt,
			}, "helpers",
				"right: helpers.Ltype2Rtype(src.left),",
				"left: helpers.Rtype2Ltype(src.right),",
				nil,
			),
		)
	})

	Describe("match function", func() {

		type input struct {
			tag     string
			expName string
			expType string
			expErr  string
			right   Fields
			vtags   []string
		}

		DescribeTable("returns field name and type",
			func(i input) {
				name, typ, err := match(i.tag, i.right, i.vtags)
				if i.expErr != "" {
					Expect(err).To(MatchError(i.expErr))
				} else {
					Expect(err).NotTo(HaveOccurred())
				}
				Expect(name).To(Equal(i.expName))
				Expect(typ).To(Equal(i.expType))
			},

			Entry("Empty tag", input{
				expErr: "empty tag",
			}),

			Entry("Empty fields map", input{
				tag: "sometag",
			}),

			Entry("Field with sts tag == field name on the right ", input{
				tag:     "Field1",
				expName: "Field1",
				expType: "string",
				right: Fields{
					"Field1": Field{Type: bt("string")},
					"Field2": Field{Type: bt("int")},
				},
			}),

			Entry("sts tag == json tag on the right and json is not valid tag",
				input{
					tag: "field_2",
					right: Fields{
						"Field1": Field{Type: bt("string")},
						"Field2": Field{Type: bt("int"), Tags: Tags{"json": "field_2"}},
					},
				}),

			Entry("sts tag == json tag on the right and json is valid tag",
				input{
					tag:     "field_2",
					expName: "Field2",
					expType: "int",
					right: Fields{
						"Field1": Field{Type: bt("string")},
						"Field2": Field{Type: bt("int"), Tags: Tags{"json": "field_2"}},
					},
					vtags: []string{"json"},
				}),

			Entry("sts tag with omitempty == json tag on the right and json is valid tag",
				input{
					tag:     "field_2,omitempty",
					expName: "Field2",
					expType: "int",
					right: Fields{
						"Field1": Field{Type: bt("string")},
						"Field2": Field{Type: bt("int"), Tags: Tags{"json": "field_2"}},
					},
					vtags: []string{"json"},
				}),

			Entry("sts tag == json tag on the right and json is valid tag with omitempty",
				input{
					tag:     "field_2",
					expName: "Field2",
					expType: "int",
					right: Fields{
						"Field1": Field{Type: bt("string")},
						"Field2": Field{Type: bt("int"), Tags: Tags{"json": "field_2,omitempty"}},
					},
					vtags: []string{"json"},
				}),

			Entry("bar tag == db tag on the right and db is valid tag", input{
				tag:     "field_3",
				expName: "Field3",
				expType: "int",
				right: Fields{
					"Field1": Field{Type: bt("string")},
					"Field2": Field{Type: bt("int"), Tags: Tags{"json": "field_2"}},
					"Field3": Field{Type: bt("int"), Tags: Tags{"db": "field_3"}},
				},
				vtags: []string{"db"},
			}),
		)
	})

	Describe("Link function", func() {

		var (
			readStruct = func(fname, sname string) Fields {
				data, err := Parse("./testdata/field/"+fname, []string{
					"sts", "json", "bar", "db", "foo",
				})
				Expect(err).NotTo(HaveOccurred())

				fields, ok := data.Structs[sname]
				Expect(ok).To(BeTrue())

				return fields
			}

			lf, rf Fields
		)

		BeforeEach(func() {
			lf = readStruct("left.go", "Left")
			rf = readStruct("right.go", "Right")
			Expect(lf).To(HaveLen(6))
			Expect(rf).To(HaveLen(6))
		})

		Context("when have two structures with tags", func() {

			It("returns filled pairs", func() {
				pairs, err := link(lf, rf, "sts", []string{"json"})
				Expect(err).NotTo(HaveOccurred())

				sort.Sort(pairs)

				Expect(pairs).To(HaveLen(3))

				Expect(pairs[0]).To(Equal(fpair{
					lf: "A", rf: "A",
					lt: "int", rt: "int",
					lp: false, rp: false,
					convertable: true, assignable: true,
					ord: 0,
				}))

				Expect(pairs[1]).To(Equal(fpair{
					lf: "B", rf: "B",
					lt: "string", rt: "string",
					lp: false, rp: false,
					convertable: true, assignable: true,
					ord: 1,
				}))

				Expect(pairs[2]).To(Equal(fpair{
					lf: "C", rf: "C",
					lt: "float32", rt: "float32",
					lp: false, rp: false,
					convertable: true, assignable: true,
					ord: 2,
				}))

			})

			It("fills field with source tag 'bar' and dest tag 'db'", func() {
				pairs, err := link(lf, rf, "bar", []string{"db"})
				Expect(err).NotTo(HaveOccurred())

				sort.Sort(pairs)

				Expect(pairs).To(HaveLen(1))

				Expect(pairs[0]).To(Equal(fpair{
					lf: "D", rf: "Double",
					lt: "int", rt: "int",
					lp: false, rp: false,
					convertable: true, assignable: true,
					ord: 3,
				}))

			})

			It("fills field with source tag 'foo'/omitemty and dest tag 'bar'", func() {
				pairs, err := link(lf, rf, "foo", []string{"bar"})
				Expect(err).NotTo(HaveOccurred())

				sort.Sort(pairs)

				Expect(pairs).To(HaveLen(2))

				Expect(pairs[0]).To(Equal(fpair{
					lf: "O", rf: "WOO",
					lt: "int", rt: "int",
					lp: false, rp: false,
					convertable: true, assignable: true,
					ord: 4,
				}))

				Expect(pairs[1]).To(Equal(fpair{
					lf: "OM", rf: "WOM",
					lt: "int", rt: "int",
					lp: false, rp: false,
					convertable: true, assignable: true,
					ord: 5,
				}))
			})

		})

	})

})
