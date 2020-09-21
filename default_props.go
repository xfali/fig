// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package fig

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"
	"text/template"
)

type Value = map[string]interface{}

type Opt func(ctx *DefaultProperties) error

type DefaultProperties struct {
	Value *Value
	Env   map[string]string

	reader ValueReader
	loader ValueLoader

	cache map[string]interface{}
	lock  sync.RWMutex
}

var Default Properties = New()

func New(opts ...Opt) *DefaultProperties {
	ret := &DefaultProperties{
		Value:  nil,
		reader: NewYamlReader(),
		loader: NewYamlLoader(),
		cache:  map[string]interface{}{},
	}

	for _, opt := range opts {
		err := opt(ret)
		if err != nil {
			logf("opt err! : %s\n", err.Error())
			return nil
		}
	}

	return ret
}

func SetValueReader(r ValueReader) Opt {
	return func(ctx *DefaultProperties) error {
		ctx.SetValueReader(r)
		return nil
	}
}

func SetValueLoader(l ValueLoader) Opt {
	return func(ctx *DefaultProperties) error {
		ctx.SetValueLoader(l)
		return nil
	}
}

func SetValue(r io.Reader) Opt {
	return func(ctx *DefaultProperties) error {
		return ctx.ReadValue(r)
	}
}

func (ctx *DefaultProperties) SetValueReader(r ValueReader) {
	ctx.reader = r
}

func (ctx *DefaultProperties) SetValueLoader(l ValueLoader) {
	ctx.loader = l
}

func (ctx *DefaultProperties) ReadValue(r io.Reader) error {
	ctx.lock.Lock()
	defer ctx.lock.Unlock()

	ctx.cache = map[string]interface{}{}
	ctx.Env = GetEnvs()

	if ctx.reader != nil {
		r, err := ctx.ExecTemplate(r)
		if err != nil {
			return err
		}
		v, err := ctx.reader.Read(r)
		if err != nil {
			return err
		}

		ctx.Value = v
	}
	return nil
}

// A.B.C
func (ctx *DefaultProperties) Get(key string, defaultValue string) string {
	//if key == "" {
	//	return defaultValue
	//}

	ctx.lock.Lock()
	defer ctx.lock.Unlock()

	if v, ok := ctx.cache[key]; ok {
		if ret, ok := v.(string); ok {
			return ret
		}
	}

	tempKey := "{{ ." + key + "}}"
	tpl, ok := template.New("").Option("missingkey=error").Parse(tempKey)
	if ok != nil {
		logf("key: %s not found(parse error)", key)
		return defaultValue
	}
	b := strings.Builder{}
	err := tpl.Execute(&b, ctx.Value)
	if err != nil {
		return defaultValue
	}

	ret := b.String()
	ctx.cache[key] = ret
	return ret
}

// 依赖于ValueReader的序列化和反序列化方式
func (ctx *DefaultProperties) GetValue(key string, result interface{}) error {
	//if key == "" {
	//	return fmt.Errorf("key is empty")
	//}
	ctx.lock.Lock()
	defer ctx.lock.Unlock()

	if v, ok := ctx.cache[key]; ok {
		if ret, ok := v.(string); ok {
			err := ctx.loader.Deserialize(ret, result)
			if err != nil {
				return fmt.Errorf("Unmarshal from cache error: %s, data: %s ", err.Error(), ret)
			}
			return nil
		}
	}

	tempKey := "{{ load_value ." + key + "}}"
	tpl, ok := template.New("").Option("missingkey=error").Funcs(template.FuncMap{
		"load_value": ctx.loader.Serialize,
	}).Parse(tempKey)
	if ok != nil {
		return fmt.Errorf("key: %s not found(parse error)", key)
	}
	b := bytes.NewBuffer(nil)
	err := tpl.Execute(b, ctx.Value)
	if err != nil {
		return fmt.Errorf("load from template failed: err: %s data: %s", err.Error(), b.String())
	}

	data := b.String()
	ctx.cache[key] = data
	err = ctx.loader.Deserialize(data, result)
	if err != nil {
		return fmt.Errorf("Unmarshal error: %s, data: %s ", err.Error(), b.String())
	}
	return nil
}

func (ctx *DefaultProperties) ExecTemplate(r io.Reader) (io.Reader, error) {
	buf := bytes.NewBuffer(nil)

	_, err := io.Copy(buf, r)
	if err != nil {
		return nil, err
	}

	tpl, ok := template.New("").Option("missingkey=error").Parse(buf.String())
	if ok != nil {
		logf("parse error")
		return nil, ok
	}

	buf.Reset()
	err = tpl.Execute(buf, ctx)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

type JsonReader struct{}

func NewJsonReader() *JsonReader {
	return &JsonReader{}
}

type JsonLoader struct{}

func NewJsonLoader() *JsonLoader {
	return &JsonLoader{}
}

func (v *JsonReader) Read(r io.Reader) (*Value, error) {
	buf := bytes.NewBuffer(nil)

	_, err := io.Copy(buf, r)
	if err != nil {
		return nil, err
	}

	ret := Value{}
	logf("value: %s\n", buf.String())
	err = json.Unmarshal(buf.Bytes(), &ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func (v *JsonLoader) Serialize(o interface{}) (string, error) {
	b, err := json.Marshal(o)
	return string(b), err
}

func (v *JsonLoader) Deserialize(value string, result interface{}) error {
	//t := reflect.TypeOf(result)
	//if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.String {
	//	v := reflect.ValueOf(result)
	//	v = v.Elem()
	//	v.SetString(value)
	//	return nil
	//}
	return json.Unmarshal([]byte(value), result)
}
