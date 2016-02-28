package option

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
)

type (
	Options2     map[string]interface{}
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

func (a []interface{}) dump(indent string) {
	for _, v := range a {
		switch use := v.(type) {
		case map[string]interface{}:
			fmt.Printf("%s(map[string]interface{})  =\r\n", indent)
			dump_map(indent+"  ", &use)
		case []interface{}:
			fmt.Print("%s([]interface{}) =\r\n", indent)
			dump_array(indent+"  ", use)
		default:
			fmt.Printf("%s(%T) = %v\r\n", indent, use, use)
		}
	}
}

func dump_array(indent string, a []interface{}) {
	for _, v := range a {
		switch use := v.(type) {
		case map[string]interface{}:
			fmt.Printf("%s(map[string]interface{})  =\r\n", indent)
			dump_map(indent+"  ", &use)
		case []interface{}:
			fmt.Print("%s([]interface{}) =\r\n", indent)
			dump_array(indent+"  ", use)
		default:
			fmt.Printf("%s(%T) = %v\r\n", indent, use, use)
		}
	}
}

func test_json() {
	str := `{
  "config": {"path": "./logs/", "format":"%T %L %N : %M"},
  "server": [2,
    "/var/www/html",
	{"host":"192.168.1.10", "port":10000, "enable":false},
	{"host":"192.168.1.11", "port":10000, "ha":null, "relay":{"host":"192.168.10.10", "port":20000}}
  ]
}`
	j := make(map[string]interface{})
	if err := json.Unmarshal([]byte(str), &j); err != nil {
		fmt.Println("转换出错: ", err.Error())
		return
	}
	fmt.Println(j)
	fmt.Println("====================================")
	fmt.Println(str)
	fmt.Println("====================================")
	dump_map("", &j)
}

func dump_map(indent string, m *map[string]interface{}) {
	for k, v := range *m {
		switch use := v.(type) {
		case map[string]interface{}:
			fmt.Printf("%s(map[string]interface{}) - %s =\r\n", indent, k)
			dump_map(indent+"  ", &use)
		case []interface{}:
			fmt.Printf("%s([]interface{}) - %s =", indent, k)
			dump_array(indent+"  ", use)
		default:
			fmt.Printf("%s(%T) - %s = %v\r\n", indent, use, k, use)
		}
	}
}

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
