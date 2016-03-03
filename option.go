package option

import (
	"github.com/i11cn/go_json"
)

/*
用来测试的数据:

{
  "config": {"path": "./logs/", "format":"%T %L %N : %M"},
  "server": [2,
    "/var/www/html",
    {"host":"192.168.1.10", "port":10000},
    {"host":"192.168.1.11", "port":10000, "relay":{"host":"192.168.10.10", "port":20000}}
  ]
}
*/

type (
	Options json.Json
)

func NewOptions() *Options {
	return (*Options)(json.NewJson())
}

func (o *Options) Get(path ...string) *Options {
	return (*Options)((*json.Json)(o).Get(path...))
}

func (o *Options) Exist(path ...string) bool {
	return o.Get(path...) != nil
}

func (o *Options) String() (string, bool) {
	return (*json.Json)(o).String()
}

func (o *Options) Int() (ret int, ok bool) {
	return (*json.Json)(o).Int()
}

func (o *Options) UInt() (ret uint, ok bool) {
	return (*json.Json)(o).UInt()
}

func (o *Options) Int64() (ret int64, ok bool) {
	return (*json.Json)(o).Int64()
}

func (o *Options) UInt64() (ret uint64, ok bool) {
	return (*json.Json)(o).UInt64()
}

func (o *Options) Float() (ret float32, ok bool) {
	return (*json.Json)(o).Float()
}

func (o *Options) Float64() (ret float64, ok bool) {
	return (*json.Json)(o).Float64()
}

func (o *Options) Bool() (ret bool, ok bool) {
	return (*json.Json)(o).Bool()
}

func (o *Options) ParseJsonFile(path string) error {
	return nil
}

func (o *Options) Set(key string, value interface{}) *Options {
	(*json.Json)(o).Set(key, value)
	return o
}

func (o *Options) Replace(key string, value interface{}) *Options {
	return o
}

func (o *Options) Append(key string, value interface{}) *Options {
	return o
}

func (o *Options) SetPath(value interface{}, path ...string) *Options {
	return o
}

func (o *Options) ReplacePath(value interface{}, path ...string) *Options {
	return o
}

func (o *Options) AppendPath(value interface{}, path ...string) *Options {
	return o
}

func (o *Options) Merge(o2 *Options) *Options {
	(*json.Json)(o).Merge((*json.Json)(o2))
	return o
}

func (o *Options) ParseJson(js string) (err error) {
	var j *json.Json
	if j, err = json.FromString(js); err == nil {
		(*json.Json)(o).Merge(j)
	}
	return
}

func (o *Options) ParseIniFile(path string) error {
	return nil
}

func (o *Options) ParseCommand() *Options {
	return o
}
