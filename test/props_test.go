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
	"time"
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
  "Value": {
    "float": 1.5
  },

  "DataSources": {
    "default": {
      "DriverName": "{{.Env.CONTEXT_TEST_ENV}}",
	  "DriverNameGet0": "{{ env "CONTEXT_TEST_ENV" }}",
      "DriverNameGet1": "{{ env "CONTEXT_TEST_ENV" "func1_return" }}",
      "DriverNameGet2": "{{ env ".Env.CONTEXT_TEST_ENV" "func2_return" }}",
      "DriverNameGet3": "{{ env "NOT_EXIST" "func3_return" }}",
	  "DriverNameGet4": "{{ env "NOT_EXIST" }}",
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
      DriverNameGet0: "{{ env "CONTEXT_TEST_ENV" }}"
      DriverNameGet1: "{{ env "CONTEXT_TEST_ENV" "func1_return" }}"
      DriverNameGet2: "{{ env ".Env.CONTEXT_TEST_ENV" "func2_return" }}"
      DriverNameGet3: "{{ env "NOT_EXIST" "func3_return" }}"
      DriverNameGet4: "{{ env "NOT_EXIST" }}"
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
		v = config.Get("DataSources.default.DriverNameGet1", "")
		if v == "" {
			t.Fatal("DataSources.default.DriverNameGet1 not found")
		}
		if v != "ONLY FOR TEST" {
			t.Fatal("not match")
		}
		v = config.Get("DataSources.default.DriverNameGet0", "")
		if v == "" {
			t.Fatal("DataSources.default.DriverNameGet0 not found")
		}
		if v != "ONLY FOR TEST" {
			t.Fatal("not match")
		}
		v = config.Get("DataSources.default.DriverNameGet1", "")
		if v == "" {
			t.Fatal("DataSources.default.DriverNameGet1 not found")
		}
		if v != "ONLY FOR TEST" {
			t.Fatal("not match")
		}
		v = config.Get("DataSources.default.DriverNameGet2", "")
		if v == "" {
			t.Fatal("DataSources.default.DriverNameGet2 not found")
		}
		if v != "ONLY FOR TEST" {
			t.Fatal("not match")
		}
		v = config.Get("DataSources.default.DriverNameGet3", "")
		if v == "" {
			t.Fatal("DataSources.default.DriverNameGet3 not found")
		}
		if v != "func3_return" {
			t.Fatal("not match")
		}
		v = config.Get("DataSources.default.DriverNameGet4", "")
		if v != "" {
			t.Fatal("DataSources.default.DriverNameGet3 must not found")
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
		v = config.Get("DataSources.default.DriverNameGet0", "")
		if v == "" {
			t.Fatal("DataSources.default.DriverNameGet0 not found")
		}
		if v != "ONLY FOR TEST" {
			t.Fatal("not match")
		}
		v = config.Get("DataSources.default.DriverNameGet1", "")
		if v == "" {
			t.Fatal("DataSources.default.DriverNameGet1 not found")
		}
		if v != "ONLY FOR TEST" {
			t.Fatal("not match")
		}
		v = config.Get("DataSources.default.DriverNameGet2", "")
		if v == "" {
			t.Fatal("DataSources.default.DriverNameGet2 not found")
		}
		if v != "ONLY FOR TEST" {
			t.Fatal("not match")
		}
		v = config.Get("DataSources.default.DriverNameGet3", "")
		if v == "" {
			t.Fatal("DataSources.default.DriverNameGet3 not found")
		}
		if v != "func3_return" {
			t.Fatal("not match")
		}
		v = config.Get("DataSources.default.DriverNameGet4", "")
		if v != "" {
			t.Fatal("DataSources.default.DriverNameGet3 must not found")
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
		config.SetValueReader(fig.NewYamlReader())
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
				t.Fatal("expect have value but get: ", ret["default"].DriverInfo, "yaml not support")
			}

			if ret["default"].MaxIdleConn != 500 {
				t.Fatal("expect 500 but get: ", ret["default"].MaxIdleConn, "yaml not support")
			}
			t.Log(ret["default"])
		}

		ret2 := map[string]database{}
		err = config.GetValue("DataSources", &ret2)
		if err != nil {
			t.Fatal(err, "ret [ default ] is nil")
		} else {
			if ret2["default"].DriverInfo != "root:123@tcp(localhost:3306)/test?charset=utf8" {
				t.Fatal("expect have value but get: ", ret2["default"].DriverInfo, "yaml not support")
			}

			if ret2["default"].MaxIdleConn != 500 {
				t.Fatal("expect 500 but get: ", ret2["default"].MaxIdleConn, "yaml not support")
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

func TestMerge(t *testing.T) {
	s1 := fig.NewSettableProperties()
	s1.SetValueReader(fig.NewYamlReader())
	s1.SetValueLoader(fig.NewYamlLoader())
	err := s1.ReadValue(strings.NewReader(test_yaml_str))
	if err != nil {
		t.Fatal(err)
	}

	s2 := fig.NewSettableProperties()

	s1.Set("a", 1)
	s1.Set("b", "2")
	s1.Set("c", 3.14)

	s2.Set("a", -1)
	s2.Set("b", "-2")
	s2.Set("c", -3.14)
	s2.Set("d", time.Now())

	m := fig.MergeProperties(s1, s2)

	v := s1.Get("LogResponse", "")
	if v == "" {
		t.Fatal("LogResponse not found")
	} else {
		t.Log(v)
	}

	ret := m.Get("a", "")
	if ret != "1" {
		t.Fatal("expect 1 but get ", ret)
	} else {
		t.Log(ret)
	}

	ret = m.Get("b", "")
	if ret != "2" {
		t.Fatal("expect 2 but get ", ret)
	} else {
		t.Log(ret)
	}

	ret = m.Get("c", "")
	if ret != "3.14" {
		t.Fatal("expect 3.14 but get ", ret)
	} else {
		t.Log(ret)
	}

	ret = m.Get("d", "")
	if ret == "" {
		t.Fatal("expect time but get ", ret)
	} else {
		t.Log(ret)
	}
}
