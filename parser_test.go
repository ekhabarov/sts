package sts

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gs "github.com/onsi/gomega/gstruct"
)

func loadFile(path string) File {
	str, err := Parse("./testdata/parser-input/"+path, nil)
	Expect(err).NotTo(HaveOccurred())

	return *str
}

var _ = Describe("Parser", func() {

	Describe("Errors", func() {

		Context("not-exists.go", func() {

			It("uses package name only", func() {
				path := "./testdata/parser-input/not-exists.go"
				str, err := Parse(path, nil)
				Expect(err).To(MatchError(fmt.Sprintf("file %q doesn't exist", path)))
				Expect(str).To(BeNil())
			})
		})
	})

	Describe("Load and parse", func() {

		Context("empty.go", func() {

			It("uses package name only", func() {

				str := loadFile("empty.go")

				Expect(str).To(gs.MatchAllFields(gs.Fields{
					"Package": Equal("whatever"),
					"Structs": BeNil(),
				}))
			})
		})

		Context("non-struct-types.go", func() {

			It("uses package name only", func() {

				str := loadFile("non-struct-types.go")

				Expect(str).To(gs.MatchAllFields(gs.Fields{
					"Package": Equal("whatever"),
					"Structs": BeNil(),
				}))
			})
		})

		Context("one-struct.go", func() {

			It("return info about one structure", func() {

				str := loadFile("one-struct.go")

				Expect(str).To(gs.MatchAllFields(gs.Fields{
					"Package": Equal("whatever"),
					"Structs": gs.MatchAllKeys(gs.Keys{
						"OneStruct": gs.MatchAllKeys(gs.Keys{

							"I": gs.MatchAllFields(gs.Fields{
								"Type":      TypeMatcher("int"),
								"IsPointer": BeFalse(),
								"Tags":      BeZero(),
								"Ord":       Equal(uint8(0)),
							}),

							"PI": gs.MatchAllFields(gs.Fields{
								"Type":      TypeMatcher("int"),
								"IsPointer": BeTrue(),
								"Tags":      BeZero(),
								"Ord":       Equal(uint8(1)),
							}),

							"S": gs.MatchAllFields(gs.Fields{
								"Type":      TypeMatcher("string"),
								"IsPointer": BeFalse(),
								"Tags":      BeZero(),
								"Ord":       Equal(uint8(2)),
							}),

							"PS": gs.MatchAllFields(gs.Fields{
								"Type":      TypeMatcher("string"),
								"IsPointer": BeTrue(),
								"Tags":      BeZero(),
								"Ord":       Equal(uint8(3)),
							}),
						}),
					}),
				}))
			})
		})

		Context("two-structs-one-is-embedded-into-another.go", func() {

			It("return info about two structures and embedded field", func() {

				str := loadFile("two-structs-one-is-embedded-into-another.go")

				Expect(str).To(gs.MatchAllFields(gs.Fields{
					"Package": Equal("whatever"),
					"Structs": gs.MatchAllKeys(gs.Keys{

						"WithEmbedded": gs.MatchAllKeys(gs.Keys{

							"I": gs.MatchAllFields(gs.Fields{
								"Type":      TypeMatcher("int"),
								"IsPointer": BeFalse(),
								"Tags":      BeZero(),
								"Ord":       Equal(uint8(0)),
							}),

							"S": gs.MatchAllFields(gs.Fields{
								"Type":      TypeMatcher("string"),
								"IsPointer": BeFalse(),
								"Tags":      BeZero(),
								"Ord":       Equal(uint8(1)),
							}),

							"embedded_0": gs.MatchAllFields(gs.Fields{
								"Type":      TypeMatcher("whatever.Embedded"),
								"IsPointer": BeFalse(),
								"Tags":      BeZero(),
								"Ord":       Equal(uint8(2)),
							}),
						}),

						"Embedded": gs.MatchAllKeys(gs.Keys{
							"CS": gs.MatchAllFields(gs.Fields{
								"Type":      TypeMatcher("string"),
								"IsPointer": BeFalse(),
								"Tags":      BeZero(),
								"Ord":       Equal(uint8(0)),
							}),
						}),
					}),
				}))
			})

		})

		Context("two-independent-structs.go", func() {

			It("return info about two structures", func() {

				str := loadFile("two-independent-structs.go")

				Expect(str).To(gs.MatchAllFields(gs.Fields{
					"Package": Equal("whatever"),
					"Structs": gs.MatchAllKeys(gs.Keys{

						"Second": gs.MatchAllKeys(gs.Keys{

							"Intf": gs.MatchAllFields(gs.Fields{
								"Type":      TypeMatcher("int"),
								"IsPointer": BeFalse(),
								"Tags":      BeZero(),
								"Ord":       Equal(uint8(0)),
							}),

							"Strf": gs.MatchAllFields(gs.Fields{
								"Type":      TypeMatcher("string"),
								"IsPointer": BeTrue(),
								"Tags":      BeZero(),
								"Ord":       Equal(uint8(1)),
							}),
						}),

						"First": gs.MatchAllKeys(gs.Keys{

							"Intf": gs.MatchAllFields(gs.Fields{
								"Type":      TypeMatcher("int"),
								"IsPointer": BeTrue(),
								"Tags":      BeZero(),
								"Ord":       Equal(uint8(0)),
							}),

							"StringF": gs.MatchAllFields(gs.Fields{
								"Type":      TypeMatcher("string"),
								"IsPointer": BeFalse(),
								"Tags":      BeZero(),
								"Ord":       Equal(uint8(1)),
							}),
						}),
					}),
				}))
			})
		})

		Context("one-struct-fields-are-of-type-SelectorExpr.go", func() { //nolint

			It("return info about one structure", func() {

				str := loadFile("one-struct-fields-are-of-type-SelectorExpr.go")

				Expect(str).To(gs.MatchAllFields(gs.Fields{
					"Package": Equal("whatever"),
					"Structs": gs.MatchAllKeys(gs.Keys{

						"Selector": gs.MatchAllKeys(gs.Keys{

							"Intf": gs.MatchAllFields(gs.Fields{
								"Type":      TypeMatcher("int"),
								"IsPointer": BeFalse(),
								"Tags":      BeZero(),
								"Ord":       Equal(uint8(0)),
							}),

							"Strf": gs.MatchAllFields(gs.Fields{
								"Type":      TypeMatcher("string"),
								"IsPointer": BeFalse(),
								"Tags":      BeZero(),
								"Ord":       Equal(uint8(1)),
							}),

							"CreatedAt": gs.MatchAllFields(gs.Fields{
								"Type":      TypeMatcher("time.Time"),
								"IsPointer": BeFalse(),
								"Tags":      BeZero(),
								"Ord":       Equal(uint8(2)),
							}),
						}),
					}),
				}))
			})
		})

		Context("one-struct-fields-are-of-slice-type.go", func() {

			It("return info about one structure", func() {

				str := loadFile("one-struct-fields-are-of-slice-type.go")

				Expect(str).To(gs.MatchAllFields(gs.Fields{
					"Package": Equal("whatever"),
					"Structs": gs.MatchAllKeys(gs.Keys{

						"Slices": gs.MatchAllKeys(gs.Keys{

							"Intf": gs.MatchAllFields(gs.Fields{
								"Type":      TypeMatcher("int"),
								"IsPointer": BeFalse(),
								"Tags":      BeZero(),
								"Ord":       Equal(uint8(0)),
							}),

							"IntSlice": gs.MatchAllFields(gs.Fields{
								"Type":      TypeMatcher("int"),
								"IsPointer": BeFalse(),
								"Tags":      BeZero(),
								"Ord":       Equal(uint8(1)),
							}),
						}),
					}),
				}))
			})
		})

		Context("one-struct-fields-are-of-struct-slice-type.go", func() {

			It("return info about one structure", func() {

				str := loadFile("one-struct-fields-are-of-struct-slice-type.go")

				Expect(str).To(gs.MatchAllFields(gs.Fields{
					"Package": Equal("whatever"),
					"Structs": gs.MatchAllKeys(gs.Keys{

						"TimeSlices": gs.MatchAllKeys(gs.Keys{

							"Intf": gs.MatchAllFields(gs.Fields{
								"Type":      TypeMatcher("int"),
								"IsPointer": BeTrue(),
								"Tags":      BeZero(),
								"Ord":       Equal(uint8(0)),
							}),

							"TimeSlice": gs.MatchAllFields(gs.Fields{
								"Type":      TypeMatcher("time.Time"),
								"IsPointer": BeFalse(),
								"Tags":      BeZero(),
								"Ord":       Equal(uint8(1)),
							}),
						}),
					}),
				}))
			})
		})

		Context("one-struct-fields-are-of-unsupported-slice-type.go", func() {

			It("return info about one structure", func() {

				str := loadFile("one-struct-fields-are-of-unsupported-slice-type.go")

				Expect(str).To(gs.MatchAllFields(gs.Fields{
					"Package": Equal("whatever"),
					"Structs": gs.MatchAllKeys(gs.Keys{
						"UnsupportedSlices": HaveLen(0),
					}),
				}))
			})
		})

		Context("one-struct-field-is-of-types-time-time.go", func() { //nolint

			It("return info about one structure", func() {

				str := loadFile("one-struct-field-is-of-types-time-time.go")

				Expect(str).To(gs.MatchAllFields(gs.Fields{
					"Package": Equal("whatever"),
					"Structs": gs.MatchAllKeys(gs.Keys{

						"TimeTime": gs.MatchAllKeys(gs.Keys{

							"T": gs.MatchAllFields(gs.Fields{
								"Type":      TypeMatcher("time.Time"),
								"IsPointer": BeFalse(),
								"Tags":      BeZero(),
								"Ord":       Equal(uint8(0)),
							}),

							"PT": gs.MatchAllFields(gs.Fields{
								"Type":      TypeMatcher("time.Time"),
								"IsPointer": BeTrue(),
								"Tags":      BeZero(),
								"Ord":       Equal(uint8(1)),
							}),

							"ThirdPartyType": gs.MatchAllFields(gs.Fields{
								"Type":      TypeMatcher("nulls.Time"),
								"IsPointer": BeFalse(),
								"Tags":      BeZero(),
								"Ord":       Equal(uint8(2)),
							}),
						}),
					}),
				}))
			})
		})

		Context("one-struct-with-unsupported-type.go", func() {

			It("return info about one structure", func() {

				str := loadFile("one-struct-with-unsupported-type.go")

				Expect(str).To(gs.MatchAllFields(gs.Fields{
					"Package": Equal("whatever"),
					"Structs": gs.MatchAllKeys(gs.Keys{
						"UnsupportedTypes": HaveLen(0),
					}),
				}))
			})

		})
	})

	Describe("Lookup", func() {

		Context("when call Lookup with existing struct", func() {

			It("returns set of fields", func() {
				str, err := Parse("./testdata/parser-input/two-independent-structs.go", nil)
				Expect(err).NotTo(HaveOccurred())

				fields, err := Lookup(str, "Second")
				Expect(err).NotTo(HaveOccurred())

				Expect(fields).To(gs.MatchAllKeys(gs.Keys{

					"Intf": gs.MatchAllFields(gs.Fields{
						"Type":      TypeMatcher("int"),
						"IsPointer": BeFalse(),
						"Tags":      BeZero(),
						"Ord":       Equal(uint8(0)),
					}),

					"Strf": gs.MatchAllFields(gs.Fields{
						"Type":      TypeMatcher("string"),
						"IsPointer": BeTrue(),
						"Tags":      BeZero(),
						"Ord":       Equal(uint8(1)),
					}),
				}))

			})
		})

		Context("when call Lookup with non-existing struct", func() {

			It("returns set of fields", func() {
				str, err := Parse("./testdata/parser-input/two-independent-structs.go", nil)
				Expect(err).NotTo(HaveOccurred())

				fields, err := Lookup(str, "NotExists")
				Expect(err).To(MatchError(`structure "NotExists" not found`))
				Expect(fields).To(BeNil())
			})
		})
	})
})
