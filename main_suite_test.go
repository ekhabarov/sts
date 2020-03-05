package sts

import (
	"go/types"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSource(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Source Suite")
}

type BasicType struct {
	kind types.BasicKind
	info types.BasicInfo
	name string
}

func (b BasicType) String() string         { return b.name }
func (b BasicType) Underlying() types.Type { return b }

func (b *BasicType) Kind() types.BasicKind { return b.kind }
func (b *BasicType) Info() types.BasicInfo { return b.info }
func (b *BasicType) Name() string          { return b.name }

func bt(n string) types.Type {
	switch n {
	case "int":
		return &BasicType{
			kind: types.Int,
			info: types.IsInteger,
			name: n,
		}

	case "string":
		return &BasicType{
			kind: types.String,
			info: types.IsString,
			name: n,
		}
	}
	return nil
}
