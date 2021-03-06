// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package fig

import (
	"io"
)

type Serializer interface {
	Serialize(o interface{}) (string, error)
}

type Deserializer interface {
	Deserialize(v string, result interface{}) error
}

type ValueReader interface {
	Read(r io.Reader) (*Value, error)
}

type ValueLoader interface {
	Serializer
	Deserializer
}

type Properties interface {
	// 配置ValueReader
	SetValueReader(r ValueReader)

	// 读取value
	ReadValue(r io.Reader) error

	// 设置ValueLoader值提取器
	SetValueLoader(l ValueLoader)

	// param: key属性名称
	// param: defaultValue: 默认值
	// return: 属性值，如不存在返回默认值
	Get(key string, defaultValue string) string

	// param: key属性名称
	// param: result: 填充对象指针
	// return: 正常返回nil,否则返回错误
	GetValue(key string, result interface{}) error
}
