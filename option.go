package option

import (
	"encoding/json"
)

type (
	Options struct {
		Name  string
		Value interface{}
	}
)

func NewOptions() *Options {
	return &Options{}
}

func (o *Options) Dump() string {
	data, err := json.Marshal(o)
	if err != nil {
		return ""
	} else {
		return string(data)
	}
}

func (o *Options) MustInt(path ...string) bool {
	return false
}

func (o *Options) MustIntWithRange(min, max int64, path ...string) bool {
	return false
}

func (o *Options) MustFloat(path ...string) bool {
	return false
}

func (o *Options) MustFloatWithRange(min, max float64, path ...string) bool {
	return false
}

func (o *Options) MustBool(path ...string) bool {
	return false
}

func (o *Options) Get(path ...string) string {
	return ""
}

func (o *Options) GetInt(path ...string) int64 {
	return 0
}

func (o *Options) GetFloat(path ...string) float64 {
	return 0.0
}

func (o *Options) ParseJsonFile(path string) error {
	return nil
}

func (o *Options) ParseJson(js string) error {
	return nil
}

func (o *Options) ParseCommand() {
}
