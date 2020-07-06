package sts

import (
	"fmt"
	"go/types"

	tt "github.com/onsi/gomega/types"
)

func TypeMatcher(exp interface{}) tt.GomegaMatcher {
	return &typMatcher{
		expected: exp,
	}
}

type typMatcher struct {
	expected interface{}
}

func (t *typMatcher) Match(actual interface{}) (bool, error) {
	s, ok := actual.(types.Type)
	if !ok {
		return false, fmt.Errorf(
			"actual value of type %T should implement types.Type interface",
			actual,
		)
	}

	return baseType(s) == t.expected.(string), nil
}

func (t *typMatcher) FailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected %q is equal to %q", baseType(actual.(types.Type)), t.expected)
}

func (t *typMatcher) NegatedFailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected %q is not equal to %q", actual, t.expected)
}
