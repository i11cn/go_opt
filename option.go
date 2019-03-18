package option

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/sanity-io/litter"
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
	Item string

	flag_info struct {
		def_type    reflect.Type
		flag_type   reflect.Type
		is_switch   bool
		default_val *reflect.Value
		usage       string
		value       *reflect.Value
		setter      func(interface{}, string) error
	}
	CommandParser struct {
		options map[string]Item
		types   []reflect.Type
		flags   map[string]flag_info
		parsed  bool
	}
)

func NewParser(t ...reflect.Type) (*CommandParser, error) {
	ret := &CommandParser{}
	return ret, ret.Bind(t...)
}

func (cp *CommandParser) Parse(cmd ...[]string) error {
	cp.options = make(map[string]Item)
	use := os.Args
	if len(cmd) > 0 {
		use = cmd[0]
	}
	for _, arg := range use {
		fmt.Println(arg)
	}
	return nil
}

func (cp *CommandParser) Bind(ts ...reflect.Type) error {
	for _, t := range ts {
		if err := cp.proc_type(t); err != nil {
			return err
		}
	}
	return nil
}

func (cp *CommandParser) Usage() string {
	ret := &bytes.Buffer{}
	return ret.String()
}

func (cp *CommandParser) Get(out interface{}) error {
	return nil
}

func (cp *CommandParser) parse_tag(t reflect.Type, name, tag1, tag2 string) (info flag_info, names []string, err error) {
	info = flag_info{}
	parts := strings.Split(tag1, ",")
	option_flag := false
	for _, p := range parts {
		switch p {
		case "switch":
			info.is_switch = true
			fallthrough
		case "option":
			option_flag = true
		default:
			if option_flag {
				u := StringConverter(p)
				info.default_val, err = u.ToType(t)
				if err != nil {
					info.default_val = nil
				}
			} else {
				if !strings.HasPrefix(p, "-") {
					if len(p) == 1 {
						p = "-" + p
					} else {
						p = "--" + p
					}
				}
				if names == nil {
					names = make([]string, 0, 5)
				}
				names = append(names, p)
			}
		}
	}
	if len(tag2) > 0 {
		info.usage = tag2
	}
	if names == nil {
		names = make([]string, 0, 1)
		names = append(names, "--"+name)
	}
	return
}

func (cp *CommandParser) make_setter(t reflect.Type, f reflect.StructField, d, v *reflect.Value) func(interface{}) {
	return func(out interface{}) {
		out_val := reflect.ValueOf(out)
		out_type := out_val.Type()
		if out_type.Kind() == reflect.Ptr {
			out_val = out_val.Elem()
			out_type = out_type.Elem()
		}
		if t != out_type {
			return
		}
		// val := out_val.FieldByIndex(f.Index)
	}
}

func (cp *CommandParser) valid_type(t reflect.Type) bool {
	switch t.String() {
	case "string", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "float32", "float64", "bool":
		return true
	case "*string", "*int", "*int8", "*int16", "*int32", "*int64", "*uint", "*uint8", "*uint16", "*uint32", "*uint64", "*float32", "*float64", "*bool":
		return true
	}
	return false
}

func (cp *CommandParser) proc_field(t reflect.Type, f reflect.StructField, tag1, tag2 string) error {
	if !cp.valid_type(f.Type) {
		return errors.New("目前还不支持所指定的类型： " + f.Type.String())
	}
	info, names, err := cp.parse_tag(f.Type, f.Name, tag1, tag2)
	if err != nil {
		return err
	}
	fmt.Println(info)
	if info.default_val != nil {
		fmt.Println(info.default_val)
	}
	litter.Dump(names)
	return nil
}

func (cp *CommandParser) parse_arg_names(name, tag string) []string {
	ret := make([]string, 0, 5)
	parts := strings.Split(tag, ",")
	for _, p := range parts {
		switch p {
		case "switch", "option":
			break
		default:
			if !strings.HasPrefix(p, "-") {
				if len(p) == 1 {
					p = "-" + p
				} else {
					p = "--" + p
				}
			}
			ret = append(ret, p)
		}
	}
	if len(ret) == 0 {
		ret = append(ret, "--"+name)
	}
	return ret
}

func (cp *CommandParser) proc_type(t reflect.Type) error {
	struct_type := t
	if struct_type.Kind() == reflect.Ptr {
		struct_type = struct_type.Elem()
	}
	if struct_type.Kind() != reflect.Struct {
		return errors.New("抱歉，只支持Struct")
	}
	for _, tmp := range cp.types {
		if t == tmp {
			return nil
		}
	}
	// 此处处理t的Tags
	for i := 0; i < struct_type.NumField(); i++ {
		f := struct_type.Field(i)
		fmt.Println(f.Name, f.Tag)
		if tag, exist := f.Tag.Lookup("cmd"); exist {
			var usage string
			if usage, exist = f.Tag.Lookup("usage"); !exist {
				usage = ""
			}
			if err := cp.proc_field(struct_type, f, tag, usage); err != nil {
				fmt.Println(err.Error())
			}
		}
	}
	cp.types = append(cp.types, struct_type)
	return nil
}
