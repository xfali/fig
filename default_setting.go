// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package fig

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/xfali/goutils/log"
	"github.com/xfali/goutils/prop"
	"io"
	"strings"
	"sync"
	"text/template"
)

type Value = map[string]interface{}

type Opt func(ctx *DefaultSetting) error

type DefaultSetting struct {
	Value *Value
	Env   map[string]string

	reader ValueReader

	cache map[string]interface{}
	lock  sync.RWMutex
}

var Default *DefaultSetting = NewSetting()

func NewSetting(opts ...Opt) *DefaultSetting {
	ret := &DefaultSetting{
		Value: nil,
		Env:   prop.GetEnvs(),

		reader: &JsonValue{},

		cache: map[string]interface{}{},
	}

	for _, opt := range opts {
		err := opt(ret)
		if err != nil {
			log.Error("opt err! : %s\n", err.Error())
			return nil
		}
	}

	return ret
}

func SetValueReader(r ValueReader) Opt {
	return func(ctx *DefaultSetting) error {
		ctx.SetValueReader(r)
		return nil
	}
}

func SetValue(r io.Reader) Opt {
	return func(ctx *DefaultSetting) error {
		return ctx.LoadValue(r)
	}
}

func (ctx *DefaultSetting) SetValueReader(r ValueReader) {
	ctx.reader = r
}

func (ctx *DefaultSetting) RefreshEnv() {
	ctx.lock.Lock()
	defer ctx.lock.Unlock()
	ctx.cache = map[string]interface{}{}

	ctx.Env = prop.GetEnvs()
}

func (ctx *DefaultSetting) LoadValue(r io.Reader) error {
	ctx.lock.Lock()
	defer ctx.lock.Unlock()
	ctx.cache = map[string]interface{}{}

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

// Env.ENVNAME
// Value.A.B.C
func (ctx *DefaultSetting) Get(key string) (string, bool) {
	if key == "" {
		return "", false
	}

	ctx.lock.Lock()
	defer ctx.lock.Unlock()

	if v, ok := ctx.cache[key]; ok {
		if ret, ok := v.(string); ok {
			return ret, true
		}
	}

	tempKey := "{{ ." + key + "}}"
	tpl, ok := template.New("").Option("missingkey=error").Parse(tempKey)
	if ok != nil {
		log.Info("key: %s not found(parse error)")
		return "", false
	}
	b := strings.Builder{}
	err := tpl.Execute(&b, ctx)
	if err != nil {
		return "", false
	}

	ret := b.String()
	ctx.cache[key] = ret
	return ret, true
}

// Value.A.B.C
// 依赖于ValueReader的序列化和反序列化方式
func (ctx *DefaultSetting) GetValue(key string, result interface{}) error {
	if key == "" {
		return fmt.Errorf("key is empty")
	}

	ctx.lock.Lock()
	defer ctx.lock.Unlock()

	if v, ok := ctx.cache[key]; ok {
		if ret, ok := v.(string); ok {
			err := ctx.reader.Deserialize(ret, result)
			if err != nil {
				return fmt.Errorf("Unmarshal from cache error: %s, data: %s ", err.Error(), ret)
			}
			return nil
		}
	}

	tempKey := "{{ load_value ." + key + "}}"
	tpl, ok := template.New("").Option("missingkey=error").Funcs(template.FuncMap{
		"load_value": ctx.reader.Serialize,
	}).Parse(tempKey)
	if ok != nil {
		return fmt.Errorf("key: %s not found(parse error)")
	}
	b := bytes.NewBuffer(nil)
	err := tpl.Execute(b, ctx)
	if err != nil {
		return fmt.Errorf("load from template failed: %s", b.String())
	}

	data := b.String()
	ctx.cache[key] = data
	err = ctx.reader.Deserialize(data, result)
	if err != nil {
		return fmt.Errorf("Unmarshal error: %s, data: %s ", err.Error(), b.String())
	}
	return nil
}

func (ctx *DefaultSetting) ExecTemplate(r io.Reader) (io.Reader, error) {
	buf := bytes.NewBuffer(nil)

	_, err := io.Copy(buf, r)
	if err != nil {
		return nil, err
	}

	tpl, ok := template.New("").Option("missingkey=error").Parse(buf.String())
	if ok != nil {
		log.Info("key: %s not found(parse error)")
		return nil, ok
	}

	buf.Reset()
	err = tpl.Execute(buf, ctx)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

type JsonValue struct{}

func (v *JsonValue) Read(r io.Reader) (*Value, error) {
	buf := bytes.NewBuffer(nil)

	_, err := io.Copy(buf, r)
	if err != nil {
		return nil, err
	}

	ret := Value{}
	log.Debug("value: %s\n", buf.String())
	err = json.Unmarshal(buf.Bytes(), &ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func (v *JsonValue) Serialize(o interface{}) (string, error) {
	b, err := json.Marshal(o)
	return string(b), err
}

func (v *JsonValue) Deserialize(value string, result interface{}) error {
	return json.Unmarshal([]byte(value), result)
}
