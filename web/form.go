package web

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/liwp-stephen/null"
	"github.com/liwp-stephen/null/zero"
	"github.com/swanwish/go-common/logs"
)

func (ctx HandlerContext) PopulateForm(v interface{}) error {
	objT := reflect.TypeOf(v)
	objV := reflect.ValueOf(v)
	if !isStructPtr(objT) {
		return fmt.Errorf("%v must be  a struct pointer", v)
	}
	objT = objT.Elem()
	objV = objV.Elem()

	return populateFormToStruct(ctx, objT, objV)
}

func isStructPtr(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}

var sliceOfInts = reflect.TypeOf([]int(nil))
var sliceOfStrings = reflect.TypeOf([]string(nil))

func populateFormToStruct(ctx HandlerContext, objT reflect.Type, objV reflect.Value) error {
	for i := 0; i < objT.NumField(); i++ {
		fieldV := objV.Field(i)
		if !fieldV.CanSet() {
			continue
		}

		fieldT := objT.Field(i)
		if fieldT.Anonymous && fieldT.Type.Kind() == reflect.Struct {
			err := populateFormToStruct(ctx, fieldT.Type, fieldV)
			if err != nil {
				return err
			}
			continue
		}

		tags := strings.Split(fieldT.Tag.Get("form"), ",")
		var tag string
		if len(tags) == 0 || len(tags[0]) == 0 {
			tag = fieldT.Name
		} else if tags[0] == "-" {
			continue
		} else {
			tag = tags[0]
		}

		value := ctx.FormValue(tag)
		if len(value) == 0 {
			continue
		}

		switch fieldT.Type.Kind() {
		case reflect.Bool:
			if strings.ToLower(value) == "on" || strings.ToLower(value) == "1" || strings.ToLower(value) == "yes" {
				fieldV.SetBool(true)
				continue
			}
			if strings.ToLower(value) == "off" || strings.ToLower(value) == "0" || strings.ToLower(value) == "no" {
				fieldV.SetBool(false)
				continue
			}
			b, err := strconv.ParseBool(value)
			if err != nil {
				return err
			}
			fieldV.SetBool(b)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			x, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}
			fieldV.SetInt(x)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			x, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return err
			}
			fieldV.SetUint(x)
		case reflect.Float32, reflect.Float64:
			x, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return err
			}
			fieldV.SetFloat(x)
		case reflect.Interface:
			fieldV.Set(reflect.ValueOf(value))
		case reflect.String:
			fieldV.SetString(value)
		case reflect.Struct:
			fieldType := fieldT.Type.String()
			var realValue interface{} = nil
			switch fieldType {
			case "time.Time":
				format := time.RFC3339
				if len(tags) > 1 {
					format = tags[1]
				}
				t, err := time.ParseInLocation(format, value, time.Local)
				if err != nil {
					return err
				}
				realValue = t
			case "zero.String", "null.String":
				if fieldType == "zero.String" {
					realValue = zero.StringFrom(value)
				} else {
					realValue = null.StringFrom(value)
				}
			case "zero.Int", "null.Int":
				x, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return err
				}
				if fieldType == "zero.Int" {
					realValue = zero.IntFrom(x)
				} else {
					realValue = null.IntFrom(x)
				}
			case "zero.Float":
				x, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return err
				}
				if fieldType == "zero.Float" {
					realValue = zero.FloatFrom(x)
				} else {
					realValue = null.FloatFrom(x)
				}
			case "zero.Bool":
				if strings.ToLower(value) == "on" || strings.ToLower(value) == "1" || strings.ToLower(value) == "yes" {
					realValue = true
				} else if strings.ToLower(value) == "off" || strings.ToLower(value) == "0" || strings.ToLower(value) == "no" {
					realValue = false
				} else {
					b, err := strconv.ParseBool(value)
					if err != nil {
						return err
					}
					realValue = b
				}
			default:
				logs.Errorf("Unhandle struct %s", fieldT.Type.String())
			}
			if realValue != nil {
				fieldV.Set(reflect.ValueOf(realValue))
			}
		case reflect.Slice:
			if fieldT.Type == sliceOfInts {
				formVals := ctx.FormValues(tag)
				fieldV.Set(reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(int(1))), len(formVals), len(formVals)))
				for i := 0; i < len(formVals); i++ {
					val, err := strconv.Atoi(formVals[i])
					if err != nil {
						return err
					}
					fieldV.Index(i).SetInt(int64(val))
				}
			} else if fieldT.Type == sliceOfStrings {
				formVals := ctx.FormValues(tag)
				fieldV.Set(reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf("")), len(formVals), len(formVals)))
				for i := 0; i < len(formVals); i++ {
					fieldV.Index(i).SetString(formVals[i])
				}
			}
		default:
			logs.Errorf("Unhandle field %v", fieldT)
		}
	}
	return nil
}
