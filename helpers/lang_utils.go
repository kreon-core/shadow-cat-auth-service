package helpers

func IsEmptyString(s *string) bool {
	return s == nil || *s == ""
}

func IsBlankString(s *string) bool {
	if IsEmptyString(s) {
		return true
	}

	for _, r := range *s {
		if r != ' ' && r != '\n' && r != '\t' && r != '\r' {
			return false
		}
	}
	return true
}

func OrDefault[T comparable](value, defaultValue T) T {
	var zero T
	if value == zero {
		return defaultValue
	}
	return value
}
