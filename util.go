// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package fig

func GetString(props Properties) func(key string, defaultValue string) string {
	return func(key string, defaultValue string) string {
		return props.Get(key, defaultValue)
	}
}

func GetBool(props Properties) func(key string, defaultValue bool) bool {
	return func(key string, defaultValue bool) bool {
		var v bool
		err := props.GetValue(key, &v)
		if err != nil {
			return defaultValue
		} else {
			return v
		}
	}
}

func GetInt(props Properties) func(key string, defaultValue int) int {
	return func(key string, defaultValue int) int {
		var v int
		err := props.GetValue(key, &v)
		if err != nil {
			return defaultValue
		} else {
			return v
		}
	}
}

func GetUint(props Properties) func(key string, defaultValue uint) uint {
	return func(key string, defaultValue uint) uint {
		var v uint
		err := props.GetValue(key, &v)
		if err != nil {
			return defaultValue
		} else {
			return v
		}
	}
}

func GetInt64(props Properties) func(key string, defaultValue int64) int64 {
	return func(key string, defaultValue int64) int64 {
		var v int64
		err := props.GetValue(key, &v)
		if err != nil {
			return defaultValue
		} else {
			return v
		}
	}
}

func GetUint64(props Properties) func(key string, defaultValue uint64) uint64 {
	return func(key string, defaultValue uint64) uint64 {
		var v uint64
		err := props.GetValue(key, &v)
		if err != nil {
			return defaultValue
		} else {
			return v
		}
	}
}

func GetFloat32(props Properties) func(key string, defaultValue float32) float32 {
	return func(key string, defaultValue float32) float32 {
		var v float32
		err := props.GetValue(key, &v)
		if err != nil {
			return defaultValue
		} else {
			return v
		}
	}
}

func GetFloat64(props Properties) func(key string, defaultValue float64) float64 {
	return func(key string, defaultValue float64) float64 {
		var v float64
		err := props.GetValue(key, &v)
		if err != nil {
			return defaultValue
		} else {
			return v
		}
	}
}
