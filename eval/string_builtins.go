package eval

import "lang/object"

func stringOtherwise(
	row *int,
	column *int,
	str object.Object,
	args ...object.Object,
) object.Object {
	s, _ := str.(*object.String)
	if len(args) != 1 {
		return newError(
			"[%d,%d] otherwise expected %d, got=%d",
			*row,
			*column,
			1,
			len(args),
		)
	}

	arg, ok := args[0].(*object.String)
	if !ok {
		return newError(
			"[%d,%d] otherwise expected it argument to be a STRING, got=%s",
			*row,
			*column,
			args[0].Type(),
		)
	}

	if s.Value == "" {
		return newString(arg.Value)
	} else {
		return s
	}
}
