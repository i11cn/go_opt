package option

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type (
	Options struct {
		value interface{}
		items map[string]*Options
	}

	opt_dump map[string]interface{}
)

func NewOptions() *Options {
	return &Options{}
}

func (o *Options) Dump() string {
	d := o.dump()
	data, err := json.Marshal(d)
	if err != nil {
		return ""
	} else {
		return string(data)
	}
}

func (o *Options) Test() {
	o.items = make(map[string]*Options)
	o.items["name"] = &Options{15.5, nil}
	o.items["test"] = &Options{nil, make(map[string]*Options)}
	o.items["test"].items["aaa"] = &Options{true, nil}
}

func (o *Options) Get(path ...string) *Options {
	if o == nil || len(path) == 0 {
		return o
	}
	cur := o
	var ok bool
	for _, p := range path {
		if cur.items == nil {
			return nil
		}
		if cur, ok = cur.items[p]; !ok {
			return nil
		}
	}
	return cur
}

func (o *Options) Exist(path ...string) bool {
	return o.Get(path...) != nil
}

func (o *Options) MustInt() bool {
	if o.value == nil {
		return false
	}
	switch v := o.value.(type) {
	case int64:
		return true
	case string:
		use, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			o.value = use
			return true
		} else {
			return false
		}
	case float64:
		o.value = int64(v)
		return true
	default:
		return false
	}
}

func (o *Options) MustIntWithRange(min, max int64) bool {
	if !o.MustInt() {
		return false
	}
	if v, ok := o.value.(int64); ok {
		return v >= min && v <= max
	} else {
		return false
	}
}

func (o *Options) MustFloat() bool {
	if o.value == nil {
		return false
	}
	switch v := o.value.(type) {
	case float64:
		return true
	case string:
		use, err := strconv.ParseFloat(v, 64)
		if err == nil {
			o.value = use
			return true
		} else {
			return false
		}
	case int64:
		o.value = float64(v)
		return true
	default:
		return false
	}
}

func (o *Options) MustFloatWithRange(min, max float64) bool {
	if !o.MustFloat() {
		return false
	}
	if v, ok := o.value.(float64); ok {
		return v >= min && v <= max
	} else {
		return false
	}
}

func (o *Options) MustBool() bool {
	if o.value == nil {
		return false
	}
	switch v := o.value.(type) {
	case bool:
		return true
	case string:
		use, err := strconv.ParseBool(v)
		if err == nil {
			o.value = use
			return true
		} else {
			return false
		}
	default:
		return false
	}
}

func (o *Options) String() string {
	if ret, ok := o.value.(string); ok {
		return ret
	} else {
		return fmt.Sprint(o.value)
	}
}

func (o *Options) Int() int64 {
	if ret, ok := o.value.(int64); ok {
		return ret
	} else {
		return 0
	}
}

func (o *Options) Float() float64 {
	if ret, ok := o.value.(float64); ok {
		return ret
	} else {
		return 0.0
	}
}

func (o *Options) Bool() bool {
	if ret, ok := o.value.(bool); ok {
		return ret
	} else {
		return false
	}
}

func (o *Options) ParseJsonFile(path string) error {
	return nil
}

func (o *Options) ParseJson(js string) error {
	return o.parse_json_string(js)
}

func (o *Options) ParseCommand() {
}
