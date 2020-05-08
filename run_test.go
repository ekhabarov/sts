package sts

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Run", func() {

	Describe("split", func() {

		DescribeTable("split",
			func(input, exp1, exp2 string) {
				got1, got2 := split(input)
				Expect(got1).To(Equal(exp1))
				Expect(got2).To(Equal(exp2))
			},

			Entry("Empty stirng", "", "", ""),
			Entry("String without semicolon", "abc", "", ""),
			Entry("Semicolon", ":", "", ""),
			Entry("Semicolon with 1st part", "abc:", "abc", ""),
			Entry("Semicolon with 2nd part", ":def", "", "def"),
			Entry("String with 2 parts divided by :", "abc:def", "abc", "def"),
			Entry("String with 3 parts divided by :", "abc:def:zzz", "", ""),
		)
	})

	Describe("pkgFromPath", func() {

		DescribeTable("pkgFromPath",
			func(input, expected string) {
				got := pkgFromPath(input)
				Expect(got).To(Equal(expected))
			},

			Entry("Empty string", "", ""),
			Entry("String without slash", "abc", "abc"),
			Entry("Root", "/", ""),
			Entry("1st level", "/abc", "abc"),
			Entry("2nd level", "/abc/def", "def"),
			Entry("3rd level", "/abc/def/ggg", "ggg"),
			Entry("Relative with dot", "./abc/def", "def"),
			Entry("Relative", "abc/def", "def"),
			Entry("Deep down", "/a/b/c/d/f/e/g/h/i/j/k/l/m/n/o/p", "p"),
		)
	})

	Describe("Run", func() {

		type input struct {
			left         string
			right        string
			sourceTag    string
			destTags     string
			outputDir    string
			helperPkg    string
			version      string
			expectedName string
			expectedErr  string
		}

		DescribeTable("cases",

			func(in input) {
				fname, content, err := Run(in.left, in.right, in.sourceTag, in.destTags,
					in.outputDir, in.helperPkg, in.version, false)
				if in.expectedErr == "" {
					Expect(err).NotTo(HaveOccurred())
				} else {
					Expect(err).To(MatchError(in.expectedErr))
				}

				Expect(fname).To(Equal(in.expectedName))

				if in.expectedName == "" {
					return
				}

				gldn := "./testdata/run/" + in.expectedName + ".golden"
				Expect(gldn).To(BeAnExistingFile())

				fc, err := ioutil.ReadFile(gldn)
				Expect(err).NotTo(HaveOccurred())

				Expect(string(content)).To(Equal(string(fc)))
			},

			Entry("Empty params",
				input{
					expectedErr: `incorrect source, format is "/path/to/file.go:struct_name"`, //nolint
				},
			),

			Entry("Source only",
				input{
					left:        "a.go:A",
					expectedErr: `incorrect destination, format is "/path/to/file.go:struct_name"`, //nolint
				},
			),

			Entry("001: Field without tags",
				input{
					left:         "./testdata/run/input/source/001_a.go:A001",
					right:        "./testdata/run/input/dest/001_b.go:B001",
					sourceTag:    "sts",
					outputDir:    ".",
					helperPkg:    "helpers",
					version:      "0.0.1",
					expectedName: "a001_to_b001.sts.go",
				},
			),

			Entry("002: Some field with tags",
				input{
					left:         "./testdata/run/input/002_a.go:A002",
					right:        "./testdata/run/input/002_b.go:B002",
					sourceTag:    "sts",
					outputDir:    ".",
					helperPkg:    "helpers",
					version:      "0.0.2",
					expectedName: "a002_to_b002.sts.go",
				},
			),

			Entry("003: Source struct not found",
				input{
					left:        "./testdata/run/input/source/001_a.go:None",
					right:       "./testdata/run/input/dest/001_b.go:B001",
					sourceTag:   "sts",
					outputDir:   ".",
					expectedErr: `source structure "None" not found: `,
				},
			),

			Entry("004: Destination struct not found",
				input{
					left:        "./testdata/run/input/source/001_a.go:A001",
					right:       "./testdata/run/input/dest/001_b.go:None",
					sourceTag:   "sts",
					outputDir:   ".",
					expectedErr: `destination structure "None" not found: `,
				},
			),

			Entry("005: Source 'foo' tag is mapped to field name on destination", input{
				left:         "./testdata/run/input/002_a.go:A002",
				right:        "./testdata/run/input/foo.go:Foo",
				sourceTag:    "foo",
				outputDir:    ".",
				version:      "0.0.5",
				expectedName: "a002_to_foo.sts.go",
			}),

			Entry("006: Source 'foo' tag is mapped to 'bar' on destination", input{
				left:         "./testdata/run/input/002_a.go:A002",
				right:        "./testdata/run/input/bar.go:Bar",
				sourceTag:    "bar",
				destTags:     "bar",
				outputDir:    ".",
				version:      "0.0.6",
				expectedName: "a002_to_bar.sts.go",
			}),
		)
	})

})
