package devicetesting

import (
	"time"

	"queclink-go/base.device.service/utils"
)

//AssertType of message after unmarshaling
func AssertType(value interface{}, originValue interface{}) interface{} {
	switch value.(type) {
	case int:
		return getIntType(originValue)
	case int32:
		return getInt32Type(originValue)
	case int16:
		return getInt16Type(originValue)
	case int64:
		return getInt64Type(originValue)
	case float32:
		return getFloat32Type(originValue)
	case float64:
		return getFloat64Type(originValue)
	case byte:
		return getByteType(originValue)
	case *utils.JSONTime:
		return getJSONTimeTypePtr(originValue)
	case utils.JSONTime:
		return getJSONTimeType(originValue)
	default:
		return originValue
	}
}

func getIntType(value interface{}) int {
	switch value.(type) {
	case int:
		return value.(int)
	case int32:
		return int(value.(int32))
	case int16:
		return int(value.(int16))
	case int64:
		return int(value.(int64))
	case float32:
		return int(value.(float32))
	case float64:
		return int(value.(float64))
	case byte:
		return int(value.(byte))
	default:
		return -1
	}
}

func getInt32Type(value interface{}) int32 {
	switch value.(type) {
	case int:
		return int32(value.(int))
	case int32:
		return int32(value.(int32))
	case int16:
		return int32(value.(int16))
	case int64:
		return int32(value.(int64))
	case float32:
		return int32(value.(float32))
	case float64:
		return int32(value.(float64))
	case byte:
		return int32(value.(byte))
	default:
		return -1
	}
}

func getInt16Type(value interface{}) int16 {
	switch value.(type) {
	case int:
		return int16(value.(int))
	case int32:
		return int16(value.(int32))
	case int16:
		return int16(value.(int16))
	case int64:
		return int16(value.(int64))
	case float32:
		return int16(value.(float32))
	case float64:
		return int16(value.(float64))
	case byte:
		return int16(value.(byte))
	default:
		return -1
	}
}

func getInt64Type(value interface{}) int64 {
	switch value.(type) {
	case int:
		return int64(value.(int))
	case int32:
		return int64(value.(int32))
	case int16:
		return int64(value.(int16))
	case int64:
		return int64(value.(int64))
	case float32:
		return int64(value.(float32))
	case float64:
		return int64(value.(float64))
	case byte:
		return int64(value.(byte))
	default:
		return -1
	}
}

func getFloat32Type(value interface{}) float32 {
	switch value.(type) {
	case int:
		return float32(value.(int))
	case int32:
		return float32(value.(int32))
	case int16:
		return float32(value.(int16))
	case int64:
		return float32(value.(int64))
	case float32:
		return float32(value.(float32))
	case float64:
		return float32(value.(float64))
	case byte:
		return float32(value.(byte))
	default:
		return -1
	}
}

func getFloat64Type(value interface{}) float64 {
	switch value.(type) {
	case int:
		return float64(value.(int))
	case int32:
		return float64(value.(int32))
	case int16:
		return float64(value.(int16))
	case int64:
		return float64(value.(int64))
	case float32:
		return float64(value.(float32))
	case float64:
		return float64(value.(float64))
	case byte:
		return float64(value.(byte))
	default:
		return -1
	}
}

func getByteType(value interface{}) byte {
	switch value.(type) {
	case int:
		return byte(value.(int))
	case int32:
		return byte(value.(int32))
	case int16:
		return byte(value.(int16))
	case int64:
		return byte(value.(int64))
	case float32:
		return byte(value.(float32))
	case float64:
		return byte(value.(float64))
	case byte:
		return byte(value.(byte))
	default:
		return 255
	}
}

func getJSONTimeTypePtr(value interface{}) *utils.JSONTime {
	v := ""
	switch value.(type) {
	case []byte:
		{
			v = string(value.([]byte))
		}
	case string:
		{
			v = value.(string)
		}
	}

	datetime, terr := time.Parse("2006-01-02T15:04:05Z", v) //parse using yyyy-MM-ddTHH:mm:ssZ
	if terr != nil {
		datetime, _ = time.Parse("2006-01-02T15:04:05Z", "1989-01-01T00:00:00Z")
	}
	return &utils.JSONTime{Time: datetime}
}

func getJSONTimeType(value interface{}) utils.JSONTime {
	v := ""
	switch value.(type) {
	case []byte:
		{
			v = string(value.([]byte))
		}
	}

	datetime, terr := time.Parse("2006-01-02T15:04:05Z", v) //parse using yyyy-MM-ddTHH:mm:ssZ
	if terr != nil {
		datetime, _ = time.Parse("2006-01-02T15:04:05Z", "1989-01-01T00:00:00Z")
	}
	return utils.JSONTime{Time: datetime}
}
