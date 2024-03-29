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
		config := fig.New(fig.SetValueReader(&fig.JsonReader{}))
		config.SetValueReader(&fig.YamlReader{})
		err := config.ReadValue(strings.NewReader(test_yaml_str))
		if err != nil {
			t.Fatal(err)
		}

		v := fig.GetBool(config)("LogResponse", false)
		if v == false {
			t.Fatal("LogResponse not found")
		}
		t.Log("LogResponse value:", v)

		port := fig.GetInt(config)("ServerPort", -1)
		if port == -1 {
			t.Fatal("ServerPort not found")
		}
		t.Log("port value:", port)

		name := fig.GetString(config)("DataSources.default.DriverName", "")
		if name == "" {
			t.Fatal("DataSources.default.DriverName not found")
		}
		t.Log("DataSources.default.DriverName value:", name)

		floatValue := fig.GetFloat32(config)("Value.float", 0)
		if floatValue == 0 {
			t.Fatal("Value.float not found")
		}
		t.Log("Value.float value:", floatValue)
	})
}

type TestStruct struct {
	dummy1            int
	Port              int  `fig:"ServerPort"`
	LogResponse       bool `fig:"LogResponse"`
	dummy2            int
	FloatValue        float32 `fig:"Value.float"`
	FloatHaveEnvValue float32 `fig:"Value.floatHaveEnv"`
	FloatEnvValue     float32 `fig:"Value.floatEnv"`
	DriverName        string  `fig:"DataSources.default.DriverName"`
	DriverNameGet1    string  `fig:"DataSources.default.DriverNameGet1"`
	DriverNameGet3    string  `fig:"DataSources.default.DriverNameGet3"`
	dummy3            int
}

func TestFill(t *testing.T) {
	config := fig.New()
	config.SetValueReader(&fig.YamlReader{})
	err := config.ReadValue(strings.NewReader(test_yaml_str))
	if err != nil {
		t.Fatal(err)
	}
	test(config, t)

	config = fig.New()
	config.SetValueReader(&fig.JsonReader{})
	config.SetValueLoader(&fig.JsonLoader{})
	err = config.ReadValue(strings.NewReader(test_config_str))
	if err != nil {
		t.Fatal(err)
	}
	test(config, t)
}

func test(config fig.Properties, t *testing.T) {
	test := TestStruct{}
	ret := fig.Fill(config, &test)

	if ret != nil {
		t.Fatal(ret)
	} else {
		if test.Port != 8080 {
			t.Fatal("expect Port 8080 got: ", test.Port)
		}
		if test.LogResponse != true {
			t.Fatal("expect LogResponse true got: ", test.LogResponse)
		}
		if test.FloatValue != 1.5 {
			t.Fatal("expect FloatValue 1.5 got: ", test.FloatValue)
		}
		if test.FloatHaveEnvValue != 1.1 {
			t.Fatal("expect FloatHaveEnvValue 1.1 got: ", test.FloatEnvValue)
		}
		if test.FloatEnvValue != 1.7 {
			t.Fatal("expect FloatEnvValue 1.7 got: ", test.FloatEnvValue)
		}
		if test.DriverName != "ONLY FOR TEST" {
			t.Fatal("expect DriverName ONLY FOR TEST got: ", test.DriverName)
		}
		if test.DriverNameGet1 != "ONLY FOR TEST" {
			t.Fatal("expect DriverName ONLY FOR TEST got: ", test.DriverName)
		}
		if test.DriverNameGet3 != "func3_return" {
			t.Fatal("expect DriverName func3_return got: ", test.DriverName)
		}
		if test.dummy1 != 0 || test.dummy2 != 0 || test.dummy3 != 0 {
			t.Fatal("dummy must be 0")
		}
		t.Log(test)
	}
}

type TestStruct2 struct {
	x           string `figPx:"DataSources.default"`
	MaxIdleConn *int
	DvrName     string `fig:"DriverName"`
	conn        int    `fig:"MaxConn"`
	dummy3      int
}

func (t *TestStruct2) test() {

}

func (t *TestStruct2) Test2() {

}

func TestFillEx(t *testing.T) {
	t.Run("json", func(t *testing.T) {
		config := fig.New()
		config.SetValueReader(&fig.JsonReader{})
		err := config.ReadValue(strings.NewReader(test_config_str))
		if err != nil {
			t.Fatal(err)
		}

		test := TestStruct2{}
		ret := fig.FillEx(config, &test, true)

		if ret != nil {
			t.Fatal(ret)
		} else {
			if *test.MaxIdleConn != 500 {
				t.Fatal("expect MaxIdleConn 500 got: ", test.MaxIdleConn)
			}
			if test.conn != 0 {
				t.Fatal("expect conn 0 got: ", test.conn)
			}
			if test.DvrName != "ONLY FOR TEST" {
				t.Fatal("expect DriverName ONLY FOR TEST got: ", test.DvrName)
			}
			if test.dummy3 != 0 {
				t.Fatal("dummy must be 0")
			}
			t.Log(test)
		}
	})

	t.Run("yaml", func(t *testing.T) {
		config := fig.New()
		config.SetValueReader(&fig.YamlReader{})
		err := config.ReadValue(strings.NewReader(test_yaml_str))
		if err != nil {
			t.Fatal(err)
		}

		test := TestStruct2{}
		ret := fig.FillEx(config, &test, true)

		if ret != nil {
			t.Fatal(ret)
		} else {
			if *test.MaxIdleConn != 500 {
				t.Fatal("expect MaxIdleConn 500 got: ", test.MaxIdleConn)
			}
			if test.conn != 0 {
				t.Fatal("expect conn 0 got: ", test.conn)
			}
			if test.DvrName != "ONLY FOR TEST" {
				t.Fatal("expect DriverName ONLY FOR TEST got: ", test.DvrName)
			}
			if test.dummy3 != 0 {
				t.Fatal("dummy must be 0")
			}
			t.Log(test)
		}
	})
}

type TestStruct3 struct {
	x           string `valuePx:"DataSources.default"`
	MaxIdleConn *int
	DvrName     string `value:"DriverName"`
	conn        int    `value:"MaxConn"`
	dummy3      int
}

type TestStruct4 struct {
	x           string `valuePx:"DataSources.default"`
	x2          string `figPx:"DataSources.default"`
	MaxIdleConn *int
	DvrName     string `fig:"DriverName"`
	conn        int    `value:"MaxConn"`
	dummy3      int
	DefaultInt  int  `value:"NoValue,default=10"`
	DefaultBool bool `value:"NoValue,default=true"`
}

func TestFillExWithTagName(t *testing.T) {
	t.Run("json", func(t *testing.T) {
		config := fig.New()
		config.SetValueReader(fig.NewJsonReader())
		config.SetValueLoader(fig.NewJsonLoader())
		err := config.ReadValue(strings.NewReader(test_config_str))
		if err != nil {
			t.Fatal(err)
		}

		test := TestStruct3{}
		ret := fig.FillExWithTagName(config, &test, true, "valuePx", "value")

		if ret != nil {
			t.Fatal(ret)
		} else {
			if *test.MaxIdleConn != 500 {
				t.Fatal("expect MaxIdleConn 500 got: ", test.MaxIdleConn)
			}
			if test.conn != 0 {
				t.Fatal("expect conn 0 got: ", test.conn)
			}
			if test.DvrName != "ONLY FOR TEST" {
				t.Fatal("expect DriverName ONLY FOR TEST got: ", test.DvrName)
			}
			if test.dummy3 != 0 {
				t.Fatal("dummy must be 0")
			}
			t.Log(test)
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

		test := TestStruct3{}
		ret := fig.FillExWithTagName(config, &test, true, "valuePx", "value")

		if ret != nil {
			t.Fatal(ret)
		} else {
			if *test.MaxIdleConn != 500 {
				t.Fatal("expect MaxIdleConn 500 got: ", test.MaxIdleConn)
			}
			if test.conn != 0 {
				t.Fatal("expect conn 0 got: ", test.conn)
			}
			if test.DvrName != "ONLY FOR TEST" {
				t.Fatal("expect DriverName ONLY FOR TEST got: ", test.DvrName)
			}
			if test.dummy3 != 0 {
				t.Fatal("dummy must be 0")
			}
			t.Log(test)
		}
	})

	t.Run("yaml and json", func(t *testing.T) {
		config := fig.New()
		config.SetValueReader(fig.NewYamlReader())
		config.SetValueLoader(fig.NewJsonLoader())
		err := config.ReadValue(strings.NewReader(test_yaml_str))
		if err != nil {
			t.Fatal(err)
		}

		test := TestStruct3{}
		ret := fig.FillExWithTagName(config, &test, true, "valuePx", "value")

		if ret != nil {
			t.Fatal(ret)
		} else {
			if *test.MaxIdleConn != 500 {
				t.Fatal("expect MaxIdleConn 500 got: ", test.MaxIdleConn)
			}
			if test.conn != 0 {
				t.Fatal("expect conn 0 got: ", test.conn)
			}
			if test.DvrName != "ONLY FOR TEST" {
				t.Fatal("expect DriverName ONLY FOR TEST got: ", test.DvrName)
			}
			if test.dummy3 != 0 {
				t.Fatal("dummy must be 0")
			}
			t.Log(test)
		}
	})
}

func TestFillExWithTagNames(t *testing.T) {
	test := TestStruct4{}
	t.Run("json", func(t *testing.T) {
		config := fig.New()
		config.SetValueReader(fig.NewJsonReader())
		config.SetValueLoader(fig.NewJsonLoader())
		err := config.ReadValue(strings.NewReader(test_config_str))
		if err != nil {
			t.Fatal(err)
		}

		ret := fig.FillExWithTagNames(config, &test, true, []string{fig.TagPrefixName, "valuePx"}, []string{fig.TagName, "value"})

		if ret != nil {
			t.Log(ret)
		}

		if *test.MaxIdleConn != 500 {
			t.Fatal("expect MaxIdleConn 500 got: ", test.MaxIdleConn)
		}
		if test.conn != 0 {
			t.Fatal("expect conn 0 got: ", test.conn)
		}
		if test.DvrName != "ONLY FOR TEST" {
			t.Fatal("expect DriverName ONLY FOR TEST got: ", test.DvrName)
		}
		if test.dummy3 != 0 {
			t.Fatal("dummy must be 0")
		}
		if test.DefaultInt != 10 {
			t.Fatal("expect 10 but get ", test.DefaultInt)
		}
		if !test.DefaultBool {
			t.Fatal("expect true but get ", test.DefaultBool)
		}
		t.Log(test)
	})

	t.Run("yaml", func(t *testing.T) {
		config := fig.New()
		config.SetValueReader(fig.NewYamlReader())
		config.SetValueLoader(fig.NewYamlLoader())
		err := config.ReadValue(strings.NewReader(test_yaml_str))
		if err != nil {
			t.Fatal(err)
		}

		ret := fig.FillExWithTagNames(config, &test, true, []string{fig.TagPrefixName, "valuePx"}, []string{fig.TagName, "value"})

		if ret != nil {
			t.Log(ret)
		}

		if *test.MaxIdleConn != 500 {
			t.Fatal("expect MaxIdleConn 500 got: ", test.MaxIdleConn)
		}
		if test.conn != 0 {
			t.Fatal("expect conn 0 got: ", test.conn)
		}
		if test.DvrName != "ONLY FOR TEST" {
			t.Fatal("expect DriverName ONLY FOR TEST got: ", test.DvrName)
		}
		if test.dummy3 != 0 {
			t.Fatal("dummy must be 0")
		}
		if test.DefaultInt != 10 {
			t.Fatal("expect 10 but get ", test.DefaultInt)
		}
		if !test.DefaultBool {
			t.Fatal("expect true but get ", test.DefaultBool)
		}
		t.Log(test)
	})

	t.Run("yaml and json", func(t *testing.T) {
		config := fig.New()
		config.SetValueReader(fig.NewYamlReader())
		config.SetValueLoader(fig.NewJsonLoader())
		err := config.ReadValue(strings.NewReader(test_yaml_str))
		if err != nil {
			t.Fatal(err)
		}

		ret := fig.FillExWithTagNames(config, &test, true, []string{fig.TagPrefixName, "valuePx"}, []string{fig.TagName, "value"})

		if ret != nil {
			t.Log(ret)
		}

		if *test.MaxIdleConn != 500 {
			t.Fatal("expect MaxIdleConn 500 got: ", test.MaxIdleConn)
		}
		if test.conn != 0 {
			t.Fatal("expect conn 0 got: ", test.conn)
		}
		if test.DvrName != "ONLY FOR TEST" {
			t.Fatal("expect DriverName ONLY FOR TEST got: ", test.DvrName)
		}
		if test.dummy3 != 0 {
			t.Fatal("dummy must be 0")
		}
		if test.DefaultInt != 10 {
			t.Fatal("expect 10 but get ", test.DefaultInt)
		}
		if !test.DefaultBool {
			t.Fatal("expect true but get ", test.DefaultBool)
		}
		t.Log(test)
	})
}
