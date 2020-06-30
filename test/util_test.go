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

func TestUtil(t *testing.T) {
	t.Run("yaml", func(t *testing.T) {
		ctx := fig.New()
		ctx.SetValueReader(&fig.YamlReader{})
		err := ctx.LoadValue(strings.NewReader(test_yaml_str))
		if err != nil {
			t.Fatal(err)
		}

		v := fig.GetBool(ctx)("LogResponse", false)
		if v == false {
			t.Fatal("LogResponse not found")
		}
		t.Log("LogResponse value:", v)

		port := fig.GetInt(ctx)("ServerPort", -1)
		if port == -1 {
			t.Fatal("ServerPort not found")
		}
		t.Log("port value:", port)

		name := fig.GetString(ctx)("DataSources.default.DriverName", "")
		if name == "" {
			t.Fatal("DataSources.Default.DriverName not found")
		}
		t.Log("DataSources.Default.DriverName value:", name)

	})
}
