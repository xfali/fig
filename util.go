// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package fig

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
)

const (
	TagPrefixName = "figPx"
	TagName       = "fig"
)

func GetString(props Properties) func(key string, defaultValue string) string {
	return func(key string, defaultValue string) string {
		return props.Get(key, defaultValue)
	}
}

func GetBool(props Properties) func(key string, defaultValue bool) bool {
	return func(key string, defaultValue bool) bool {
		var v bool
		err := props.GetValue(key, &v)
		if err != nil {
			return defaultValue
		} else {
			return v
		}
	}
}

func GetInt(props Properties) func(key string, defaultValue int) int {
	return func(key string, defaultValue int) int {
		var v int
		err := props.GetValue(key, &v)
		if err != nil {
			return defaultValue
		} else {
			return v
		}
	}
}

func GetUint(props Properties) func(key string, defaultValue uint) uint {
	return func(key string, defaultValue uint) uint {
		var v uint
		err := props.GetValue(key, &v)
		if err != nil {
			return defaultValue
		} else {
			return v
		}
	}
}

func GetInt64(props Properties) func(key string, defaultValue int64) int64 {
	return func(key string, defaultValue int64) int64 {
		var v int64
		err := props.GetValue(key, &v)
		if err != nil {
			return defaultValue
		} else {
			return v
		}
	}
}

func GetUint64(props Properties) func(key string, defaultValue uint64) uint64 {
	return func(key string, defaultValue uint64) uint64 {
		var v uint64
		err := props.GetValue(key, &v)
		if err != nil {
			return defaultValue
		} else {
			return v
		}
	}
}

func GetFloat32(props Properties) func(key string, defaultValue float32) float32 {
	return func(key string, defaultValue float32) float32 {
		var v float32
		err := props.GetValue(key, &v)
		if err != nil {
			return defaultValue
		} else {
			return v
		}
	}
}

func GetFloat64(props Properties) func(key string, defaultValue float64) float64 {
	return func(key string, defaultValue float64) float64 {
		var v float64
		err := props.GetValue(key, &v)
		if err != nil {
			return defaultValue
		} else {
			return v
		}
	}
}

func LoadFile(filename string, reader ValueReader, loader ValueLoader) (Properties, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	prop := New()
	prop.SetValueReader(reader)
	prop.SetValueLoader(loader)
	err = prop.ReadValue(f)
	return prop, err
}

func LoadJsonFile(filename string) (Properties, error) {
	return LoadFile(filename, NewJsonReader(), NewJsonLoader())
}

func LoadYamlFile(filename string) (Properties, error) {
	return LoadFile(filename, NewYamlReader(), NewYamlLoader())
}

// param: prop 属性
// param: result 填充的struct
// result: result如果不为struct的指针返回错误，填充时异常返回错误
func Fill(prop Properties, result interface{}) error {
	return FillEx(prop, result, false)
}

// param: prop 属性
// param: result 填充的struct
// param: withField 是否根据field name填充
// result: result如果不为struct的指针返回错误，填充时异常返回错误
func FillEx(prop Properties, result interface{}, withField bool) error {
	return FillExWithTagName(prop, result, withField, TagPrefixName, TagName)
}

// param: prop 属性
// param: result 填充的struct
// param: withField 是否根据field name填充
// param: tagPxName tag前缀名，后续都使用tagPxName定义的名称做前缀
// param: tagName tag名
// result: result如果不为struct的指针返回错误，填充时异常返回错误
func FillExWithTagName(prop Properties, result interface{}, withField bool, tagPxName, tagName string) error {
	t := reflect.TypeOf(result)
	v := reflect.ValueOf(result)

	if t.Kind() != reflect.Ptr {
		return errors.New("result must be ptr")
	}
	t = t.Elem()
	v = v.Elem()

	if t.Kind() != reflect.Struct {
		return errors.New("result must be struct ptr")
	}

	prefix := ""
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(tagPxName)
		if tag != "" {
			prefix = tag
			continue
		}
		tag = field.Tag.Get(tagName)
		if tag != "" {
			if tag == "-" {
				continue
			}
		} else if withField {
			tag = field.Name
		}

		if tag != "" {
			if prefix != "" {
				tag = prefix + "." + tag
			}
			c := reflect.New(field.Type).Interface()
			err := prop.GetValue(tag, c)
			if err != nil {
				logf(err.Error())
			}
			fieldValue := v.Field(i)
			if fieldValue.CanSet() {
				fieldValue.Set(reflect.ValueOf(c).Elem())
			}
		}
	}
	return nil
}

func GetEnvs() map[string]string {
	s := os.Environ()
	ret := map[string]string{}
	for _, env := range s {
		env := strings.TrimSpace(env)
		if env != "" {
			pair := strings.Split(env, "=")
			if len(pair) == 2 {
				ret[pair[0]] = pair[1]
			}
		}
	}

	return ret
}

type logFunc func(format string, o ...interface{})

var logf logFunc = func(format string, o ...interface{}) {
	fmt.Printf(format, o...)
}

// 配置fig的内置log
func SetLog(log logFunc) {
	logf = log
}
