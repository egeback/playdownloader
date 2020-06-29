package utils

import "strconv"

//GetBoolValueFromString ...
func GetBoolValueFromString(value string, defaultValue bool) *bool {
	intValue, err := strconv.ParseBool(value)
	if err != nil {
		return &defaultValue
	}
	return &intValue
}
