package settings

import (
	"fmt"
	"strconv"
)

func ToMap(opts ...Option) map[string]any {
	result := make(map[string]any)
	for _, opt := range opts {
		result[opt.Name] = opt.Value
	}
	return result
}

func GetOption(name string, opts ...Option) (Option, bool) {
	for _, opt := range opts {
		if opt.Name == name {
			return opt, true
		}
	}
	return Option{}, false
}

func FilterOptions(names []string, opts ...Option) []Option {
	contains := func(name string) bool {
		for _, val := range names {
			if val == name {
				return true
			}
		}
		return false
	}

	result := make([]Option, 0)
	for _, opt := range opts {
		if contains(opt.Name) {
			result = append(result, opt)
		}
	}
	return result
}

type Option struct {
	Name  string
	Value interface{}
}

func (opt Option) String() string {
	return fmt.Sprintf("%s: %v", opt.Name, opt.Value)
}

func (opt Option) StringValue() string {
	return fmt.Sprintf("%v", opt.Value)
}

func (opt Option) IntValue() int {
	switch v := opt.Value.(type) {
	case int:
		return v
	case int8:
		return int(v)
	case int16:
		return int(v)
	case int32:
		return int(v)
	case int64:
		return int(v)
	case uint:
		return int(v)
	case uint8:
		return int(v)
	case uint16:
		return int(v)
	case uint32:
		return int(v)
	case uint64:
		return int(v)
	case string:
		icast, _ := strconv.Atoi(v)
		return icast
	}
	return 0
}

func (opt Option) FloatValue() float64 {
	if f, ok := opt.Value.(float64); ok {
		return f
	}
	return 0.0
}

func (opt Option) BoolValue() bool {
	switch 	v := opt.Value.(type) {
	case bool:
		return v
	case string:
		bval, _ := strconv.ParseBool(v)
		return bval
	}
	return false
}
