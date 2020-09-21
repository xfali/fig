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
		config, err := fig.LoadFile("config.json", fig.NewJsonReader(), fig.NewJsonLoader())
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
		config.SetValueReader(fig.NewYamlReader())
		config.SetValueLoader(fig.NewYamlLoader())
		err = config.ReadValue(f)
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

func TestFileAll(t *testing.T) {
	t.Run("json", func(t *testing.T) {
		config, err := fig.LoadFile("config.json", fig.NewJsonReader(), fig.NewJsonLoader())
		if err != nil {
			t.Fatal(err)
		}

		ret := map[string]interface{}{}
		err = config.GetValue("", &ret)
		if err != nil {
			t.Fatal(err)
		}
		if ret["ServerPort"].(float64) != 8080 {
			t.Fatal("expect ServerPort but get ", ret["ServerPort"])
		}
		t.Log("value:", ret)

		str := config.Get("", "")
		if str == "" {
			t.Fatal("str is empty")
		}
		t.Log(str)
	})

	t.Run("yaml", func(t *testing.T) {
		f, err := os.Open("config.yaml")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		config := fig.New()
		config.SetValueReader(fig.NewYamlReader())
		config.SetValueLoader(fig.NewYamlLoader())
		err = config.ReadValue(f)
		if err != nil {
			t.Fatal(err)
		}

		ret := map[string]interface{}{}
		err = config.GetValue("", &ret)
		if err != nil {
			t.Fatal(err)
		}
		if ret["ServerPort"].(float64) != 8080 {
			t.Fatal("expect ServerPort but get ", ret["ServerPort"])
		}
		t.Log("value:", ret)

		str := config.Get("", "")
		if str == "" {
			t.Fatal("str is empty")
		}
		t.Log(str)
	})
}

func TestContext(t *testing.T) {
	t.Run("json", func(t *testing.T) {
		config := fig.New(fig.SetValueReader(&fig.JsonReader{}), fig.SetValueLoader(fig.NewJsonLoader()))
		err := config.ReadValue(strings.NewReader(test_config_str))
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
		config.SetValueLoader(fig.NewYamlLoader())
		err := config.ReadValue(strings.NewReader(test_yaml_str))
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
	type database struct {
		DriverName      string
		DriverInfo      string
		MaxConn         int
		MaxIdleConn     int
		ConnMaxLifetime int
	}

	t.Run("json", func(t *testing.T) {
		config := fig.New()
		config.SetValueReader(fig.NewJsonReader())
		config.SetValueLoader(fig.NewJsonLoader())
		err := config.ReadValue(strings.NewReader(test_config_str))
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

		// json cannot unmarshal
		if err != nil {
			t.Log(err)
		} else {
			t.Log(p)
		}

		port := 0
		err = config.GetValue("ServerPort", &port)
		if err != nil || port != 8080 {
			t.Fatal(err, port)
		} else {
			t.Log(port)
		}

		ret := map[string]*database{}
		err = config.GetValue("DataSources", &ret)
		if err != nil || ret["default"] == nil {
			t.Fatal(err, "ret [ default ] is nil")
		} else {
			if ret["default"].DriverInfo != "root:123@tcp(localhost:3306)/test?charset=utf8" {
				t.Fatal("expect have value but get: ", ret["default"].DriverInfo)
			}

			if ret["default"].MaxIdleConn != 500 {
				t.Fatal("expect 500 but get: ", ret["default"].MaxIdleConn)
			}
			t.Log(ret["default"])
		}

		ret2 := map[string]database{}
		err = config.GetValue("DataSources", &ret2)
		if err != nil {
			t.Fatal(err, "ret [ default ] is nil")
		} else {
			if ret2["default"].DriverInfo != "root:123@tcp(localhost:3306)/test?charset=utf8" {
				t.Fatal("expect have value but get: ", ret2["default"].DriverInfo)
			}

			if ret2["default"].MaxIdleConn != 500 {
				t.Fatal("expect 500 but get: ", ret2["default"].MaxIdleConn)
			}
			t.Log(ret2["default"])
		}
	})

	t.Run("yaml", func(t *testing.T) {
		config := fig.New()
		config.SetValueReader(fig.NewYamlReader())
		config.SetValueLoader(fig.NewYamlLoader())
		err := config.ReadValue(strings.NewReader(test_yaml_str))
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
		if err != nil || p != "8080" {
			t.Fatal(err)
		} else {
			t.Log(p)
		}

		port := 0
		err = config.GetValue("ServerPort", &port)
		if err != nil || port != 8080 {
			t.Fatal(err, port)
		} else {
			t.Log(port)
		}

		ret := map[string]*database{}
		err = config.GetValue("DataSources", &ret)
		if err != nil || ret["default"] == nil {
			t.Fatal(err, "ret [ default ] is nil")
		} else {
			if ret["default"].DriverInfo != "root:123@tcp(localhost:3306)/test?charset=utf8" {
				t.Log("expect have value but get: ", ret["default"].DriverInfo, "yaml not support")
			}

			if ret["default"].MaxIdleConn != 500 {
				t.Log("expect 500 but get: ", ret["default"].MaxIdleConn, "yaml not support")
			}
			t.Log(ret["default"])
		}

		ret2 := map[string]database{}
		err = config.GetValue("DataSources", &ret2)
		if err != nil {
			t.Fatal(err, "ret [ default ] is nil")
		} else {
			if ret2["default"].DriverInfo != "root:123@tcp(localhost:3306)/test?charset=utf8" {
				t.Log("expect have value but get: ", ret2["default"].DriverInfo, "yaml not support")
			}

			if ret2["default"].MaxIdleConn != 500 {
				t.Log("expect 500 but get: ", ret2["default"].MaxIdleConn, "yaml not support")
			}
			t.Log(ret2["default"])
		}
	})

	t.Run("yaml reader and json loader", func(t *testing.T) {
		config := fig.New()
		config.SetValueReader(fig.NewYamlReader())
		config.SetValueLoader(fig.NewJsonLoader())
		err := config.ReadValue(strings.NewReader(test_yaml_str))
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

		// json cannot unmarshal
		if err != nil {
			t.Log(err)
		} else {
			t.Log(p)
		}

		port := 0
		err = config.GetValue("ServerPort", &port)
		if err != nil || port != 8080 {
			t.Fatal(err, port)
		} else {
			t.Log(port)
		}

		ret := map[string]*database{}
		err = config.GetValue("DataSources", &ret)
		if err != nil || ret["default"] == nil {
			t.Fatal(err, "ret [ default ] is nil")
		} else {
			if ret["default"].DriverInfo != "root:123@tcp(localhost:3306)/test?charset=utf8" {
				t.Fatal("expect have value but get: ", ret["default"].DriverInfo)
			}

			if ret["default"].MaxIdleConn != 500 {
				t.Fatal("expect 500 but get: ", ret["default"].MaxIdleConn)
			}
			t.Log(ret["default"])
		}

		ret2 := map[string]database{}
		err = config.GetValue("DataSources", &ret2)
		if err != nil {
			t.Fatal(err, "ret [ default ] is nil")
		} else {
			if ret2["default"].DriverInfo != "root:123@tcp(localhost:3306)/test?charset=utf8" {
				t.Fatal("expect have value but get: ", ret2["default"].DriverInfo)
			}

			if ret2["default"].MaxIdleConn != 500 {
				t.Fatal("expect 500 but get: ", ret2["default"].MaxIdleConn)
			}
			t.Log(ret2["default"])
		}
	})
}

func TestFromCache(t *testing.T) {
	config := fig.New()
	err := config.ReadValue(strings.NewReader(test_config_str))
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
