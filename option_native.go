package option

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
)

type (
	option_map   map[string]interface{}
	option_array []interface{}

	json_value interface {
		MarshalJSON() ([]byte, error)
		UnmarshalJSON([]byte) error
	}
	json_string string
	json_number struct {
		value interface{}
	}
	json_bool bool
	json_null interface{}

	json_array  []json_value
	json_object struct {
		name  string
		value json_value
	}
)

func (j *json_string) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString("\"")
	buf.WriteString(*(*string)(j))
	buf.WriteString("\"")
	return buf.Bytes(), nil
}

func (j *json_string) UnmarshalJSON(d []byte) error {
	return nil
}

func (j *json_number) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString("")
	switch v := j.value.(type) {
	case float32, float64, uint, uint32, uint64, int, int32, int64:
		buf.WriteString(fmt.Sprint(v))
	default:
		return nil, errors.New("not a number format")
	}
	return buf.Bytes(), nil
}

func (j *json_number) UnmarshalJSON(d []byte) error {
	return nil
}

func (o *Options) make_path(key ...string) *Options {
	ret := o
	return ret
}

func (o *Options) get_child_by_key(key string, create bool) *Options {
	return nil
}

func (o *Options) get_value() (ret interface{}, exist bool) {
	return nil, false
}

func (o *Options) marshal_json(buf *bytes.Buffer) error {
	return nil
}

func (o *Options) parse_json_object(j *simplejson.Json) error {
	return nil
}

func (o *Options) parse_json_string(js string) error {
	return nil
}

func (o *Options) parse_ini_file() error {
	return nil
}

func (o *Options) parse_command_line() error {
	return nil
}
