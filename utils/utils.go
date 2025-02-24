package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func SliceToStructs[T any](data [][]any) ([]*T, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	result := make([]*T, len(data))
	if len(data) == 0 {
		return result, nil
	}
	var temp T
	targetType := reflect.TypeOf(temp)
	if targetType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("目标类型必须是结构体，当前类型是: %v", targetType.Kind())
	}
	for index, row := range data {
		if row == nil {
			continue
		}
		ptr := reflect.New(targetType)
		newStruct := ptr.Elem()
		if len(row) > targetType.NumField() {
			return nil, fmt.Errorf("数据列数(%d)超过结构体字段数(%d)", len(row), targetType.NumField())
		}
		for i, value := range row {
			if value == nil {
				continue
			}
			field := newStruct.Field(i)

			if !field.CanSet() {
				return nil, fmt.Errorf("无法设置字段 %s", targetType.Field(i).Name)
			}
			val := reflect.ValueOf(value)
			if val.Type().ConvertibleTo(field.Type()) {
				field.Set(val.Convert(field.Type()))
			} else {
				if convertedVal, err := convertValue(value, field.Type()); err == nil {
					field.Set(convertedVal)
				} else {
					return nil, fmt.Errorf("无法将类型 %v 转换为 %v ERROR: %v", val.Type(), field.Type(), err)
				}
			}
		}
		result[index] = ptr.Interface().(*T)
	}

	return result, nil
}

func convertValue(value any, targetType reflect.Type) (reflect.Value, error) {
	val := reflect.ValueOf(value)
	if value == nil {
		return reflect.Zero(targetType), nil
	}

	switch targetType.Kind() {
	case reflect.String:
		// 任何类型都可以转换为字符串
		return reflect.ValueOf(fmt.Sprintf("%v", value)), nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var intVal int64
		switch v := value.(type) {
		case string:
			if i, err := strconv.ParseInt(v, 10, 64); err == nil {
				intVal = i
			} else if f, err := strconv.ParseFloat(v, 64); err == nil {
				intVal = int64(f)
			} else {
				return reflect.Value{}, fmt.Errorf("无法将字符串 %v 转换为整数", v)
			}
		case float32, float64:
			intVal = int64(reflect.ValueOf(v).Float())
		case bool:
			if v {
				intVal = 1
			}
		default:
			return reflect.Value{}, fmt.Errorf("无法将类型 %T 转换为整数", value)
		}
		return reflect.ValueOf(intVal).Convert(targetType), nil

	case reflect.Float32, reflect.Float64:
		var floatVal float64
		switch v := value.(type) {
		case string:
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				floatVal = f
			} else {
				return reflect.Value{}, fmt.Errorf("无法将字符串 %v 转换为浮点数", v)
			}
		case int, int8, int16, int32, int64:
			floatVal = float64(reflect.ValueOf(v).Int())
		case bool:
			if v {
				floatVal = 1.0
			}
		default:
			return reflect.Value{}, fmt.Errorf("无法将类型 %T 转换为浮点数", value)
		}
		return reflect.ValueOf(floatVal).Convert(targetType), nil

	case reflect.Bool:
		var boolVal bool
		switch v := value.(type) {
		case string:
			// 字符串转布尔值
			switch strings.ToLower(v) {
			case "1", "t", "true", "yes", "y", "on":
				boolVal = true
			case "0", "f", "false", "no", "n", "off", "":
				boolVal = false
			default:
				return reflect.Value{}, fmt.Errorf("无法将字符串 %v 转换为布尔值", v)
			}
		case int, int8, int16, int32, int64:
			boolVal = reflect.ValueOf(v).Int() != 0
		case float32, float64:
			boolVal = reflect.ValueOf(v).Float() != 0
		default:
			return reflect.Value{}, fmt.Errorf("无法将类型 %T 转换为布尔值", value)
		}
		return reflect.ValueOf(boolVal), nil

	case reflect.Slice:
		switch val.Kind() {
		case reflect.String:
			if targetType.Elem().Kind() == reflect.Uint8 {
				return reflect.ValueOf([]byte(val.String())), nil
			}
		case reflect.Slice:
			newSlice := reflect.MakeSlice(targetType, val.Len(), val.Cap())
			for i := 0; i < val.Len(); i++ {
				elem := val.Index(i).Interface()
				convertedElem, err := convertValue(elem, targetType.Elem())
				if err != nil {
					return reflect.Value{}, fmt.Errorf("转换切片元素失败: %v", err)
				}
				newSlice.Index(i).Set(convertedElem)
			}
			return newSlice, nil
		}

	case reflect.Struct:
		if targetType == reflect.TypeOf(time.Time{}) {
			switch v := value.(type) {
			case string:
				layouts := []string{
					time.RFC3339,
					"2006-01-02 15:04:05",
					"2006-01-02",
					time.RFC822,
					time.RFC850,
					time.DateOnly,
					time.TimeOnly,
					"2006-01-02T15:04:05.00",
					"2006-01-02T15:04:05",
				}
				for _, layout := range layouts {
					if t, err := time.Parse(layout, v); err == nil {
						return reflect.ValueOf(t), nil
					}
				}
				return reflect.Value{}, fmt.Errorf("无法将字符串 %v 解析为时间", v)
			case int64:
				return reflect.ValueOf(time.Unix(v, 0)), nil
			}
		}
	}

	return reflect.Value{}, fmt.Errorf("无法将类型 %T 转换为 %v", value, targetType)
}
