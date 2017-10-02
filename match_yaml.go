package matcho

import (
	"fmt"
	"strings"

	"github.com/kr/pretty"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
	yaml "gopkg.in/yaml.v2"
)

func MatchYAMLWithDiffs(yaml interface{}) types.GomegaMatcher {
	return &MatchYAMLWithDiffsMatcher{
		YAMLToMatch: yaml,
	}
}

type MatchYAMLWithDiffsMatcher struct {
	YAMLToMatch interface{}
	Diffs       []string
}

func (matcher *MatchYAMLWithDiffsMatcher) Match(actual interface{}) (success bool, err error) {
	actualString, expectedString, err := matcher.toStrings(actual)
	if err != nil {
		return false, err
	}

	var aval interface{}
	var eval interface{}

	if err := yaml.Unmarshal([]byte(actualString), &aval); err != nil {
		return false, fmt.Errorf("Actual '%s' should be valid YAML, but it is not.\nUnderlying error:%s", actualString, err)
	}
	if err := yaml.Unmarshal([]byte(expectedString), &eval); err != nil {
		return false, fmt.Errorf("Expected '%s' should be valid YAML, but it is not.\nUnderlying error:%s", expectedString, err)
	}

	diffs := pretty.Diff(aval, eval)
	matcher.Diffs = diffs

	if len(diffs) > 0 {
		return false, nil
	}

	return true, nil
}

func (matcher *MatchYAMLWithDiffsMatcher) FailureMessage(actual interface{}) (message string) {
	actualString, expectedString, _ := matcher.toNormalisedStrings(actual)
	message = format.Message(actualString, "to match YAML of", expectedString)
	message = fmt.Sprintf("%s\nDiffs: (expected vs actual)\n%s\n", message, strings.Join(matcher.Diffs, "\n"))
	return message
}

func (matcher *MatchYAMLWithDiffsMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	actualString, expectedString, _ := matcher.toNormalisedStrings(actual)
	message = format.Message(actualString, "not to match YAML of", expectedString)
	message = fmt.Sprintf("%s\nDiffs: (expected vs actual)\n%s\n", message, strings.Join(matcher.Diffs, "\n"))

	return message
}

func (matcher *MatchYAMLWithDiffsMatcher) toNormalisedStrings(actual interface{}) (actualFormatted, expectedFormatted string, err error) {
	actualString, expectedString, err := matcher.toStrings(actual)
	return normalise(actualString), normalise(expectedString), err
}

func normalise(input string) string {
	var val interface{}
	err := yaml.Unmarshal([]byte(input), &val)
	if err != nil {
		panic(err) // guarded by Match
	}
	output, err := yaml.Marshal(val)
	if err != nil {
		panic(err) // guarded by Unmarshal
	}
	return strings.TrimSpace(string(output))
}

func (matcher *MatchYAMLWithDiffsMatcher) toStrings(actual interface{}) (actualFormatted, expectedFormatted string, err error) {
	actualString, ok := toString(actual)
	if !ok {
		return "", "", fmt.Errorf("MatchYAMLWithDiffsMatcher matcher requires a string, stringer, or []byte.  Got actual:\n%s", format.Object(actual, 1))
	}
	expectedString, ok := toString(matcher.YAMLToMatch)
	if !ok {
		return "", "", fmt.Errorf("MatchYAMLWithDiffsMatcher matcher requires a string, stringer, or []byte.  Got expected:\n%s", format.Object(matcher.YAMLToMatch, 1))
	}

	return actualString, expectedString, nil
}

func toString(a interface{}) (string, bool) {
	aString, isString := a.(string)
	if isString {
		return aString, true
	}

	aBytes, isBytes := a.([]byte)
	if isBytes {
		return string(aBytes), true
	}

	aStringer, isStringer := a.(fmt.Stringer)
	if isStringer {
		return aStringer.String(), true
	}

	return "", false
}
