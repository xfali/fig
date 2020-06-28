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

func init() {
    os.Setenv("CONTEXT_TEST_ENV", "ONLY FOR TEST")
}

func TestContext(t *testing.T) {
    ctx := fig.NewSetting()
    err := ctx.LoadValue(strings.NewReader(test_ctx_str))
    if err != nil {
        t.Fatal(err)
    }

    v, ok := ctx.Get("Value.LogResponse")
    if !ok {
        t.Fatal("Value.LogResponse not found")
    }
    t.Log("env value:", v)

    v, ok = ctx.Get("Value.DataSources.default.DriverName")
    if !ok {
        t.Fatal("Value.DataSources.default.DriverName not found")
    }
    t.Log("Value.DataSources.default.DriverName value:", v)
}

func TestContextGetValue(t *testing.T) {
    ctx := fig.NewSetting()
    err := ctx.LoadValue(strings.NewReader(test_ctx_str))
    if err != nil {
        t.Fatal(err)
    }

    v := ""
    err = ctx.GetValue("Value.Env", &v)
    if err != nil {
        t.Fatal(err)
    }
    t.Log(v)

    p := ""
    err = ctx.GetValue("Value.ServerPort", &p)
    if err != nil {
        t.Log(err)
    } else {
        t.Fatal(p)
    }

    port := 0
    err = ctx.GetValue("Value.ServerPort", &port)
    if err != nil {
        t.Fatal(err)
    } else {
        t.Log(port)
    }
}

func TestFromCache(t *testing.T) {
    ctx := fig.NewSetting()
    err := ctx.LoadValue(strings.NewReader(test_ctx_str))
    if err != nil {
        t.Fatal(err)
    }

    v, ok := ctx.Get("Value.LogResponse")
    if !ok {
        t.Fatal("Value.LogResponse not found")
    }
    t.Log("env value:", v)

    v, ok = ctx.Get("Value.LogResponse")
    if !ok {
        t.Fatal("Value.LogResponse not found")
    }
    t.Log("env value:", v)
}

func BenchmarkGet(b *testing.B) {
    ctx := fig.NewSetting()
    err := ctx.LoadValue(strings.NewReader(test_ctx_str))
    if err != nil {
        b.Fatal(err)
    }
    for i := 0; i < b.N; i++ {
        v, ok := ctx.Get("Value.LogResponse")
        if v == "" || !ok {
            b.Fatal("Value.LogResponse not found")
        }
    }
}
