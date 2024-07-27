package params

import (
	"github.com/veigaribo/qveen/util"
	"strings"
)

type ParamError struct {
	Path   []any
	Reason string
}

func MakeParamError(path []any, reason string) ParamError {
	return ParamError{
		Path:   path,
		Reason: reason,
	}
}

func (e ParamError) Error() string {
	var builder strings.Builder

	builder.WriteRune('`')
	util.WritePathString(e.Path, &builder)
	builder.WriteRune('`')
	builder.WriteRune(' ')
	builder.WriteString(e.Reason)

	return builder.String()
}
