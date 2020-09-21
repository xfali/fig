// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package test

import (
	"github.com/xfali/fig"
	"strings"
	"testing"
)

func BenchmarkGet(b *testing.B) {
	b.Run("json", func(b *testing.B) {
		config := fig.New(fig.SetValueReader(&fig.JsonReader{}))
		err := config.ReadValue(strings.NewReader(test_config_str))
		if err != nil {
			b.Fatal(err)
		}
		for i := 0; i < b.N; i++ {
			v := config.Get("LogResponse", "")
			if v == "" {
				b.Fatal("LogResponse not found")
			}
		}
	})

	b.Run("yaml", func(b *testing.B) {
		config := fig.New(fig.SetValueReader(&fig.YamlReader{}))
		err := config.ReadValue(strings.NewReader(test_yaml_str))
		if err != nil {
			b.Fatal(err)
		}
		for i := 0; i < b.N; i++ {
			v := config.Get("LogResponse", "")
			if v == "" {
				b.Fatal("LogResponse not found")
			}
		}
	})
}
