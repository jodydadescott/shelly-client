package util

func CompareBool(a, b *bool) bool {

	if a == nil {
		return b == nil
	}

	if b == nil {
		return false
	}

	return *a == *b
}

func CompareString(a, b *string) bool {

	if a == nil {
		return b == nil
	}

	if b == nil {
		return false
	}

	return *a == *b
}

func CompareInt(a, b *int) bool {

	if a == nil {
		return b == nil
	}

	if b == nil {
		return false
	}

	return *a == *b
}

func CompareFloat64(a, b *float64) bool {

	if a == nil {
		return b == nil
	}

	if b == nil {
		return false
	}

	return *a == *b
}

func CompareFloat32(a, b *float32) bool {

	if a == nil {
		return b == nil
	}

	if b == nil {
		return false
	}

	return *a == *b
}

func CompareStringSlice(a, b []string) bool {

	if a == nil {
		return b == nil
	}

	if b == nil {
		return false
	}

	has := func(x string, s []string) bool {
		for _, v := range s {
			if v == x {
				return true
			}
		}
		return false
	}

	for _, v := range a {
		if !has(v, b) {
			return false
		}
	}

	for _, v := range b {
		if !has(v, a) {
			return false
		}
	}

	return true
}
