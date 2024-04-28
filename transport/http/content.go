package http

func isString[T any]() bool {
	var check T

	_, isTypeString := any(check).(*string)
	if !isTypeString {
		_, isTypeString = any(check).(string)
	}

	return isTypeString
}
