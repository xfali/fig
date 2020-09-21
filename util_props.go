// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package fig

import "io"

func MergeProperties(props ...Properties) Properties {
	return &mergedProperties{
		props: props,
	}
}

type mergedProperties struct {
	props []Properties
}

// 配置ValueReader
func (p *mergedProperties) SetValueReader(r ValueReader) {
	panic("MergedProperties cannot reset ValueReader")
}

// 读取value
func (p *mergedProperties) ReadValue(r io.Reader) error {
	panic("MergedProperties cannot ReadValue")
}

// 设置ValueLoader值提取器
func (p *mergedProperties) SetValueLoader(l ValueLoader) {
	panic("MergedProperties cannot reset ValueLoader")
}

// param: key属性名称
// param: defaultValue: 默认值
// return: 属性值，如不存在返回默认值
func (p *mergedProperties) Get(key string, defaultValue string) string {
	l := len(p.props)
	if l > 0 {
		for i := 0; i < l-1; i++ {
			value := p.props[i].Get(key, "")
			if value != "" {
				return value
			}
		}
		return p.props[l-1].Get(key, defaultValue)
	}
	return defaultValue
}

// param: key属性名称
// param: result: 填充对象指针
// return: 正常返回nil,否则返回错误
func (p *mergedProperties) GetValue(key string, result interface{}) (err error) {
	for i := range p.props {
		err = p.props[i].GetValue(key, result)
		if err == nil {
			return nil
		}
	}
	return
}

type SettableProperties struct {
	DefaultProperties
}

func NewSettableProperties(opts ...Opt) *SettableProperties {
	ret := &SettableProperties{
		DefaultProperties: *New(opts...),
	}
	ret.Value = &Value{}
	return ret
}

func (p *SettableProperties) Set(key string, value interface{}) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	(*p.Value)[key] = value

	return nil
}

func (p *SettableProperties) Delete(key string) {
	p.lock.Lock()
	defer p.lock.Unlock()

	delete(*p.Value, key)
}
