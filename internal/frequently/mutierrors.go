package frequently

import (
	"fmt"
	"strings"
)

type ErrorList struct {
	Errs []error
}

func (el ErrorList) Error() string {
	var str strings.Builder
	/*if len(el.Errs) > 0 {
		fmt.Fprintf(&str, "multiple errors(%d): ", len(el.Errs))
	}*/
	for i, err := range el.Errs {
		fmt.Fprintf(&str, "Err(%d/%d): ", i+1, len(el.Errs))
		str.WriteString(err.Error() + "\n")
	}
	return str.String()
}
