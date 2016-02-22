package option

import (
	"fmt"
)

type (
	Options struct {
		value interface{}
		child map[string]*Options
	}
)

func NewOptions() *Options {
	return &Options{}
}

func (o *Options) Marshal() (ret string, err error) {
	return
}

func (o *Options) Get(path ...string) *Options {
	if o == nil || len(path) == 0 {
		return o
	}
	return nil
}

func (o *Options) Exist(path ...string) bool {
	return o.Get(path...) != nil
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
	ret = 0
	ok = false
	switch v := o.value.(type) {
	case int:
		ret = int(v)
		ok = true
	case uint:
		ret = int(v)
		ok = true
	case int32:
		ret = int(v)
		ok = true
	case uint32:
		ret = int(v)
		ok = true
	case int64:
		ret = int(v)
		ok = true
	case uint64:
		ret = int(v)
		ok = true
	}
	return
}

func (o *Options) UInt() (ret uint, ok bool) {
	return
}

func (o *Options) Int64() (ret int, ok bool) {
	return
}

func (o *Options) UInt64() (ret uint, ok bool) {
	return
}

func (o *Options) Float() (ret float64, ok bool) {
	return
}

func (o *Options) Bool() (ret bool, ok bool) {
	return
}

func (o *Options) ParseJsonFile(path string) error {
	return nil
}

func (o *Options) Set(o2 *Options) {
}

func (o *Options) ParseJson(js string) error {
	return o.parse_json_string(js)
}

func (o *Options) ParseIniFile(path string) error {
	return nil
}

func (o *Options) ParseCommand() {
}
