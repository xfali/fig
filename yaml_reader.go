// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package fig

import (
	"bytes"
	"github.com/xfali/goutils/log"
	"gopkg.in/yaml.v2"
	"io"
)

type YamlReader struct{}

func (v *YamlReader) Read(r io.Reader) (*Value, error) {
	buf := bytes.NewBuffer(nil)

	_, err := io.Copy(buf, r)
	if err != nil {
		return nil, err
	}

	ret := Value{}
	log.Debug("value: %s\n", buf.String())
	err = yaml.Unmarshal(buf.Bytes(), &ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func (v *YamlReader) Serialize(o interface{}) (string, error) {
	b, err := yaml.Marshal(o)
	return string(b), err
}

func (v *YamlReader) Deserialize(value string, result interface{}) error {
	return yaml.Unmarshal([]byte(value), result)
}
