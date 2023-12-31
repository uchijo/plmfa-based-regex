package eval

import (
	"errors"

	"github.com/uchijo/plmfa-based-regex/model"
)

type InputBuffer struct {
	Input string
}

func (ib InputBuffer) Len() int {
	return len(ib.Input)
}

func (ib InputBuffer) CanConsume(matcher model.CharContainer) (bool, string) {
	prefixLen := matcher.Len()
	inputLen := len(ib.Input)
	if inputLen < prefixLen {
		return false, ""
	}
	prefix := []rune(ib.Input)[:prefixLen]
	return matcher.WholeMatches(string(prefix)), string([]rune(ib.Input)[:matcher.Len()])
}

func (ib InputBuffer) Consumed(rawMatcher string) (InputBuffer, error) {
	matcher := stringContainer(rawMatcher)
	if ok, _ := ib.CanConsume(matcher); !ok {
		return InputBuffer{}, errors.New("cannot consume " + string(matcher) + " from " + ib.Input)
	}
	ib.Input = string([]rune(ib.Input)[matcher.Len():])
	return ib, nil
}

func (ib InputBuffer) Appended(prefix string) (InputBuffer, error) {
	ib.Input = prefix + ib.Input
	return ib, nil
}
