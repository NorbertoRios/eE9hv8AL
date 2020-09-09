package utils

import "reflect"

//SliceContains checks slice for element existence
func SliceContains(target interface{}, list interface{}) bool {
	if reflect.TypeOf(list).Kind() == reflect.Slice || reflect.TypeOf(list).Kind() == reflect.Array {
		listvalue := reflect.ValueOf(list)
		for i := 0; i < listvalue.Len(); i++ {
			if target == listvalue.Index(i).Interface() {
				return true
			}
		}
	}
	return false
}

//ReverseSlise reverses slice
func ReverseSlise(input []byte) []byte {
	if len(input) == 0 {
		return input
	}
	return append(ReverseSlise(input[1:]), input[0])
}
