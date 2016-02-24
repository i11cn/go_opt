package option

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"reflect"
	"strings"
)

type (
	Options struct {
		name  string
		value interface{}
		child map[string]*Options
	}
)

func NewOptions() *Options {
	return &Options{child: map[string]*Options{}}
}

func (o *Options) Marshal() (ret string, err error) {
	return
}

func (o *Options) Get(path ...string) *Options {
	if o == nil || len(path) == 0 {
		return o
	}
	current := o
	var exist bool
	for _, k := range path {
		if current, exist = current.child[k]; !exist {
			return nil
		}
	}
	return current
}

func (o *Options) Exist(path ...string) bool {
	return o.Get(path...) != nil
}

func (o *Options) Key() string {
	return o.name
}

func (o *Options) String() string {
	switch v := o.value.(type) {
	case string:
		return v
	default:
		return fmt.Sprint(v)
	}
}

func (o *Options) Int() (ret int, ok bool) {
	ok = true
	switch v := o.value.(type) {
	case float32, float64:
		ret = int(reflect.ValueOf(v).Float())
	case uint, uint32, uint64:
		ret = int(reflect.ValueOf(v).Uint())
	case int, int32, int64:
		ret = int(reflect.ValueOf(v).Int())
	default:
		ret = 0
		ok = false
	}
	return
}

func (o *Options) UInt() (ret uint, ok bool) {
	ok = true
	switch v := o.value.(type) {
	case float32, float64:
		ret = uint(reflect.ValueOf(v).Float())
	case uint, uint32, uint64:
		ret = uint(reflect.ValueOf(v).Uint())
	case int, int32, int64:
		ret = uint(reflect.ValueOf(v).Int())
	default:
		ret = 0
		ok = false
	}
	return
}

func (o *Options) Int64() (ret int64, ok bool) {
	ok = true
	switch v := o.value.(type) {
	case float32, float64:
		ret = int64(reflect.ValueOf(v).Float())
	case uint, uint32, uint64:
		ret = int64(reflect.ValueOf(v).Uint())
	case int, int32, int64:
		ret = int64(reflect.ValueOf(v).Int())
	default:
		ret = 0
		ok = false
	}
	return
}

func (o *Options) UInt64() (ret uint64, ok bool) {
	ok = true
	switch v := o.value.(type) {
	case float32, float64:
		ret = uint64(reflect.ValueOf(v).Float())
	case uint, uint32, uint64:
		ret = uint64(reflect.ValueOf(v).Uint())
	case int, int32, int64:
		ret = uint64(reflect.ValueOf(v).Int())
	default:
		ret = 0
		ok = false
	}
	return
}

func (o *Options) Float() (ret float32, ok bool) {
	ok = true
	switch v := o.value.(type) {
	case float32, float64:
		ret = float32(reflect.ValueOf(v).Float())
	case uint, uint32, uint64:
		ret = float32(reflect.ValueOf(v).Uint())
	case int, int32, int64:
		ret = float32(reflect.ValueOf(v).Int())
	default:
		ret = 0
		ok = false
	}
	return
}

func (o *Options) Float64() (ret float64, ok bool) {
	ok = true
	switch v := o.value.(type) {
	case float32, float64:
		ret = float64(reflect.ValueOf(v).Float())
	case uint, uint32, uint64:
		ret = float64(reflect.ValueOf(v).Uint())
	case int, int32, int64:
		ret = float64(reflect.ValueOf(v).Int())
	default:
		ret = 0
		ok = false
	}
	return
}

func (o *Options) Bool() (ret bool, ok bool) {
	ok = true
	switch v := o.value.(type) {
	case float32, float64:
		ret = reflect.ValueOf(v).Float() != 0
	case uint, uint32, uint64:
		ret = reflect.ValueOf(v).Uint() != 0
	case int, int32, int64:
		ret = reflect.ValueOf(v).Int() != 0
	case string:
		switch strings.ToUpper(v) {
		case "T", "Y", "TRUE", "YES", "1":
			ret = true
		case "F", "N", "FALSE", "NO", "0":
			ret = false
		default:
			ret = false
			ok = false
		}
	default:
		ret = false
		ok = false
	}
	return
}

func (o *Options) ParseJsonFile(path string) error {
	return nil
}

func (o *Options) Set(key string, value interface{}) *Options {
	o.name = key
	o.value = value
	return o
}

func (o *Options) SetChild(key string, value interface{}) *Options {
	if ov, ok := o.child[key]; ok {
		ov.name = key
		ov.value = value
	} else {
		ov := NewOptions()
		ov.name = key
		ov.value = value
		o.child[key] = ov
	}
	return o
}

func (o *Options) Merge(o2 *Options) *Options {
	if o2 != nil {
		o.value = o2.value
		for k, v := range o2.child {
			if ov, exist := o.child[k]; exist {
				ov.Merge(v)
			} else {
				ov := NewOptions()
				ov.name = k
				o.child[k] = ov
				ov.Merge(v)
			}
		}
	}
	return o
}

func (o *Options) SetJson(j *simplejson.Json) *Options {
	return o
}

func (o *Options) ParseJson(js string) error {
	return o.parse_json_string(js)
}

func (o *Options) ParseIniFile(path string) error {
	return nil
}

func (o *Options) ParseCommand() *Options {
	return o
}
