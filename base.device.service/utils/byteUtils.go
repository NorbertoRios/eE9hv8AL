package utils

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"unsafe"
)

//GetBytesValue returns one of known type from byte
func GetBytesValue(data []byte, pos int, size int) (interface{}, error) {
	if len(data) < pos+size {
		return nil, errors.New("Data array out of bound")
	}

	switch size {
	case 1:
		return byte(data[pos]), nil
	case 2:
		return int32(data[pos])<<8 | int32(data[pos+1]), nil
	case 4:
		return int32(data[pos])<<24 | int32(data[pos+1])<<16 | int32(data[pos+2])<<8 | int32(data[pos+3]), nil
	default:
		slice := data[pos : pos+size]
		return slice, nil
	}
}

//BitIsSet check bit
func BitIsSet(value int64, i uint) bool {
	return (value & (1 << i)) != 0
}

//Reverse byte array
func Reverse(numbers []byte) {
	for i, j := 0, len(numbers)-1; i < j; i, j = i+1, j-1 {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
}

//AllIs checks array contains only value elements
func AllIs(bytes *[]byte, value byte) bool {
	all := true
	for _, e := range *bytes {
		if e != 0 {
			all = false
			break
		}
	}
	return all
}

//ByteToString converts byte array to hex string
func ByteToString(bytes []byte) string {
	result := ""
	for _, v := range bytes {
		result += fmt.Sprintf("%02X", v)
	}
	return result
}

//ConvertStringToByte return bytes set from string
func ConvertStringToByte(value string) ([]byte, error) {
	result := make([]byte, 0)

	for i := 0; i <= len(value)-2; i += 2 {
		h := value[i : i+2]
		decoded, err := hex.DecodeString(h)
		if err != nil {
			return nil, errors.New("[utils]ConvertStringToByte Convert hex to int error")
		}
		result = append(result, decoded[0])
	}
	return result, nil
}

//InsertNth insert symbol to specified position
func InsertNth(s string, n int, splitter rune) string {
	var buffer bytes.Buffer
	var n1 = n - 1
	var l1 = len(s) - 1
	for i, rune := range s {
		buffer.WriteRune(rune)
		if i%n == n1 && i != l1 {
			buffer.WriteRune(splitter)
		}
	}
	return buffer.String()
}

//Bool2i converts bool to integer
func Bool2i(b bool) byte {
	if b {
		return 1
	}
	return 0
}

const intSize int = int(unsafe.Sizeof(0))

//IsBigEndian checks system byte order
func IsBigEndian() (ret bool) {
	i := int(0x1)
	bs := (*[intSize]byte)(unsafe.Pointer(&i))
	if bs[0] == 0 {
		return true
	}
	return false
}

//GetLowestBits returns 4 lowest bits
func GetLowestBits(value byte) byte {
	return value & 15
}

//GetHighestBits returns 4 highest bits
func GetHighestBits(value byte) byte {
	return value >> 4
}
