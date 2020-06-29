// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package fig

import (
	"io"
)

type ValueReader interface {
	Read(r io.Reader) (*Value, error)
	Serialize(o interface{}) (string, error)
	Deserialize(v string, result interface{}) error
}

type Properties interface {
	//配置ValueReader
	SetValueReader(r ValueReader)
	//刷新环境变量
	RefreshEnv()
	//读取setting value
	LoadValue(r io.Reader) error
	// Env.ENVNAME
	// Value.A.B.C
	Get(key string, defaultValue string) string
	// Value.A.B.C
	// 依赖于ValueReader的序列化和反序列化方式
	GetValue(key string, result interface{}) error
}
