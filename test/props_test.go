// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package test

import (
	"github.com/xfali/fig"
	"os"
	"strings"
	"testing"
)

var test_config_str = `
{
  "Env": "dev",
  "LogResponse": true,
  "LogRequestBody": true,
  "LogLevel": 1,
  "LogInnerLevel": 1,
  "LogClient": true,
  "ServerPort": 8080,

  "DataSources": {
    "default": {
      "DriverName": "{{.Env.CONTEXT_TEST_ENV}}",
      "DriverInfo": "root:123@tcp(localhost:3306)/test?charset=utf8",
      "MaxConn": 1000,
      "MaxIdleConn": 500,
      "ConnMaxLifetime": 1000
    }
  }
}
`

var test_yaml_str = `
  Env: "dev"
  LogResponse: true
  LogRequestBody: true
  LogLevel: 1
  LogInnerLevel: 1
  LogClient: true
  ServerPort: 8080
  Value:
    float: 1.5

  DataSources: 
    default: 
      DriverName: "{{.Env.CONTEXT_TEST_ENV}}"
      DriverInfo: "root:123@tcp(localhost:3306)/test?charset=utf8"
      MaxConn: 1000
      MaxIdleConn: 500
      ConnMaxLifetime: 1000
`

func init() {
	os.Setenv("CONTEXT_TEST_ENV", "ONLY FOR TEST")
}

func TestFile(t *testing.T) {
	t.Run("json", func(t *testing.T) {
		config, err := fig.LoadFile("config.json", fig.NewJsonReader())
		if err != nil {
			t.Fatal(err)
		}

		v := config.Get("LogResponse", "")
		if v == "" {
			t.Fatal("LogResponse not found")
		}
		t.Log("env value:", v)

		v = config.Get("DataSources.default.DriverName", "")
		if v == "" {
			t.Fatal("DataSources.default.DriverName not found")
		}
		t.Log("DataSources.default.DriverName value:", v)
	})

	t.Run("yaml", func(t *testing.T) {
		f, err := os.Open("config.yaml")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		config := fig.New()
		config.SetValueReader(&fig.YamlReader{})
		err = config.LoadValue(f)
		if err != nil {
			t.Fatal(err)
		}

		v := config.Get("LogResponse", "")
		if v == "" {
			t.Fatal("LogResponse not found")
		}
		t.Log("env value:", v)

		v = config.Get("DataSources.default.DriverName", "")
		if v == "" {
			t.Fatal("DataSources.default.DriverName not found")
		}
		t.Log("DataSources.default.DriverName value:", v)
	})
}

func TestContext(t *testing.T) {
	t.Run("json", func(t *testing.T) {
		config := fig.New(fig.SetValueReader(&fig.JsonReader{}))
		err := config.LoadValue(strings.NewReader(test_config_str))
		if err != nil {
			t.Fatal(err)
		}

		v := config.Get("LogResponse", "")
		if v == "" {
			t.Fatal("LogResponse not found")
		}
		t.Log("env value:", v)

		v = config.Get("DataSources.default.DriverName", "")
		if v == "" {
			t.Fatal("DataSources.default.DriverName not found")
		}
		t.Log("DataSources.default.DriverName value:", v)
	})

	t.Run("yaml", func(t *testing.T) {
		config := fig.New()
		config.SetValueReader(&fig.YamlReader{})
		err := config.LoadValue(strings.NewReader(test_yaml_str))
		if err != nil {
			t.Fatal(err)
		}

		v := config.Get("LogResponse", "")
		if v == "" {
			t.Fatal("LogResponse not found")
		}
		t.Log("env value:", v)

		v = config.Get("DataSources.default.DriverName", "")
		if v == "" {
			t.Fatal("DataSources.default.DriverName not found")
		}
		t.Log("DataSources.default.DriverName value:", v)
	})
}

func TestContextGetValue(t *testing.T) {
	config := fig.New()
	err := config.LoadValue(strings.NewReader(test_config_str))
	if err != nil {
		t.Fatal(err)
	}

	v := ""
	err = config.GetValue("Env", &v)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(v)

	p := ""
	err = config.GetValue("ServerPort", &p)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(p)
	}

	port := 0
	err = config.GetValue("ServerPort", &port)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(port)
	}
}

func TestFromCache(t *testing.T) {
	config := fig.New()
	err := config.LoadValue(strings.NewReader(test_config_str))
	if err != nil {
		t.Fatal(err)
	}

	v := config.Get("LogResponse", "")
	if v == "" {
		t.Fatal("LogResponse not found")
	}
	t.Log("LogResponse value:", v)

	v = config.Get("LogResponse", "")
	if v == "" {
		t.Fatal("LogResponse not found")
	}
	t.Log("LogResponse value:", v)
}

func BenchmarkGet(b *testing.B) {
	b.Run("json", func(b *testing.B) {
		config := fig.New(fig.SetValueReader(&fig.JsonReader{}))
		err := config.LoadValue(strings.NewReader(test_config_str))
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
		err := config.LoadValue(strings.NewReader(test_yaml_str))
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
