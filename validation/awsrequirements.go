package validate

import "regexp"

var ipRe = regexp.MustCompile(
	`^(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9]?[0-9])(\.(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9]?[0-9])){3}$`,
)

var objectkeyRe = regexp.MustCompile(
	`^[A-Za-z0-9!._*'()-]{1,1024}$`,
)

func isIPv4(s string) bool {
	return ipRe.MatchString(s)
}

func BucketnameValidation(bucket_name string) bool {
	if !(3 <= len(bucket_name) && len(bucket_name) <= 63) {
		return false
	}

	for i := 0; i < len(bucket_name); i++ {
		sym := bucket_name[i]
		if !(('a' <= sym && sym <= 'z') || ('0' <= sym && sym <= '9') || sym == '.' || sym == '-') {
			return false
		}
		if i+1 < len(bucket_name) && bucket_name[i] == '.' && bucket_name[i+1] == '.' {
			return false
		}
	}

	end := len(bucket_name) - 1
	if bucket_name[0] == '.' || bucket_name[0] == '-' || bucket_name[end] == '.' || bucket_name[end] == '-' {
		return false
	}

	if isIPv4(bucket_name) {
		return false
	}

	return true
}

func ObejectkeyValidation(object_key string) bool {
	return objectkeyRe.MatchString(object_key)
}
