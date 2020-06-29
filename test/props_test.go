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

var test_ctx_str = `
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

func TestContext(t *testing.T) {
    t.Run("json", func(t *testing.T) {
        ctx := fig.New()
        err := ctx.LoadValue(strings.NewReader(test_ctx_str))
        if err != nil {
            t.Fatal(err)
        }

        v := ctx.Get("LogResponse", "")
        if v == "" {
            t.Fatal("LogResponse not found")
        }
        t.Log("env value:", v)

        v = ctx.Get("DataSources.default.DriverName", "")
        if v == "" {
            t.Fatal("DataSources.default.DriverName not found")
        }
        t.Log("DataSources.default.DriverName value:", v)
    })

    t.Run("yaml", func(t *testing.T) {
        ctx := fig.New()
        ctx.SetValueReader(&fig.YamlValue{})
        err := ctx.LoadValue(strings.NewReader(test_yaml_str))
        if err != nil {
            t.Fatal(err)
        }

        v := ctx.Get("LogResponse", "")
        if v == "" {
            t.Fatal("LogResponse not found")
        }
        t.Log("env value:", v)

        v = ctx.Get("DataSources.default.DriverName", "")
        if v == "" {
            t.Fatal("DataSources.default.DriverName not found")
        }
        t.Log("DataSources.default.DriverName value:", v)
    })
}

func TestContextGetValue(t *testing.T) {
    ctx := fig.New()
    err := ctx.LoadValue(strings.NewReader(test_ctx_str))
    if err != nil {
        t.Fatal(err)
    }

    v := ""
    err = ctx.GetValue("Env", &v)
    if err != nil {
        t.Fatal(err)
    }
    t.Log(v)

    p := ""
    err = ctx.GetValue("ServerPort", &p)
    if err != nil {
        t.Log(err)
    } else {
        t.Fatal(p)
    }

    port := 0
    err = ctx.GetValue("ServerPort", &port)
    if err != nil {
        t.Fatal(err)
    } else {
        t.Log(port)
    }
}

func TestFromCache(t *testing.T) {
    ctx := fig.New()
    err := ctx.LoadValue(strings.NewReader(test_ctx_str))
    if err != nil {
        t.Fatal(err)
    }

    v := ctx.Get("LogResponse", "")
    if v == "" {
        t.Fatal("LogResponse not found")
    }
    t.Log("LogResponse value:", v)

    v = ctx.Get("LogResponse", "")
    if v == "" {
        t.Fatal("LogResponse not found")
    }
    t.Log("LogResponse value:", v)
}

func BenchmarkGet(b *testing.B) {
    b.Run("json", func(b *testing.B) {
        ctx := fig.New()
        err := ctx.LoadValue(strings.NewReader(test_ctx_str))
        if err != nil {
            b.Fatal(err)
        }
        for i := 0; i < b.N; i++ {
            v := ctx.Get("LogResponse", "")
            if v == "" {
                b.Fatal("LogResponse not found")
            }
        }
    })

    b.Run("yaml", func(b *testing.B) {
        ctx := fig.New(fig.SetValueReader(&fig.YamlValue{}))
        err := ctx.LoadValue(strings.NewReader(test_yaml_str))
        if err != nil {
            b.Fatal(err)
        }
        for i := 0; i < b.N; i++ {
            v := ctx.Get("LogResponse", "")
            if v == "" {
                b.Fatal("LogResponse not found")
            }
        }
    })
}
