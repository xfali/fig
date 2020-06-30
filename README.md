# fig

fig是一个轻量化的配置读取工具

## 安装
```
go get github.com/xfali/fig
```

## 使用
### 加载配置内容
```
config := fig.New()
err := config.LoadValue(strings.NewReader(test_ctx_str))
if err != nil {
    b.Fatal(err)
}
```
或者
```
config, err := fig.LoadFile("config.json", fig.NewJsonReader())
if err != nil {
    t.Fatal(err)
}
```
### 通过key获取属性值（字符串）
```
v := config.Get("DataSources.default.DriverName", "")
```
### 通过key获得反序列化对象
```
port := 0
err = config.GetValue("ServerPort", &port)
```
## 工具方法
```
config, _ := fig.LoadFile("config.json", fig.NewJsonReader())

v := fig.GetBool(config)("LogResponse", false)

floatValue := fig.GetFloat32(config)("Value.float", 0)
```

## tag
fig提供直接填充struct的field的方法，使用tag:"fig"来标识属性名称：
```
type TestStruct struct {
	dummy1      int
	Port        int  `fig:"ServerPort"`
	LogResponse bool `fig:"LogResponse"`
	dummy2      int
	FloatValue  float32 `fig:"Value.float"`
	DriverName  string  `fig:"DataSources.default.DriverName"`
	dummy3      int
}
```
使用fig.Fill方法根据tag填充struct：
```
config, _ := fig.LoadFile("config.json", fig.NewJsonReader())
test := TestStruct{}
err := fig.Fill(config, &test)
t.log(test)
```