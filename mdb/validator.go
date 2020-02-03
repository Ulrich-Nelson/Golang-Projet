package mdb

import "regexp"

func ValidateID(id string) bool {
	regex := regexp.MustCompile(`^[a-fA-F0-9]{32}$`)
	return regex.MatchString(id)
}

func ValidateName(name string) bool {
	re := regexp.MustCompile(`[a-zA-Z_]{1,16}`)
	return re.Match([]byte(name))
}

func ValidatePassword(password string) bool {
	return validatePrintableAscii(password, 32)
}

func validatePrintableAscii(s string, maxLen int) bool {
	if len(s) > maxLen {
		return false
	}
	for _, c := range s {
		if int(c) < 32 || int(c) > 255 {
			return false
		}
	}
	return true
}
