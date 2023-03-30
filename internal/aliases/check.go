package aliases

func IsValid(alias string) bool {
	if len(alias) > 63 {
		return false
	}

	var valid = true
	for _, char := range alias {
		if char >= '0' && char <= '9' { // is an ASCII number
			continue
		}

		if char >= 'a' && char <= 'z' { // is an ASCII char between a and z
			continue
		}

		valid = false
		break
	}
	return valid
}
