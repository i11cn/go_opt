package option

import (
	"bytes"
	"errors"
	"os"
	"reflect"
	"strings"
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
		flag_type   reflect.Type
		flags       []string
		is_switch   bool
		is_option   bool
		default_str string
		default_val *reflect.Value
		usage       string
		value       *reflect.Value
		parsed      bool
	}
	CommandParser struct {
		types   []reflect.Type
		flags   map[string]*flag_info
		setters map[string]func(reflect.Value) error
		parsed  bool
	}
)

func tag_parser_names(name, tag string) []string {
	ret := make([]string, 0, 5)
	parts := strings.Split(tag, ",")
	for _, p := range parts {
		if p == "switch" || p == "option" {
			break
		}
		if !strings.HasPrefix(p, "-") {
			if len(p) == 1 {
				p = "-" + p
			} else {
				p = "--" + p
			}
		}
		ret = append(ret, p)
	}
	if len(ret) == 0 {
		ret = append(ret, "--"+name)
	}
	return ret
}

func tag_parser_is_switch(tag string) bool {
	parts := strings.Split(tag, ",")
	for _, p := range parts {
		if p == "switch" {
			return true
		}
	}
	return false
}

func tag_parser_has_default(tag string) (opt bool, dft string) {
	parts := strings.Split(tag, ",")
	for _, p := range parts {
		switch p {
		case "switch", "option":
			opt = true
		default:
			if opt {
				dft = p
				return
			}
		}
	}
	return false, ""
}

func new_tag_info(t reflect.Type, tag1, tag2 string) (info *flag_info, err error) {
	info = &flag_info{}
	info.flag_type = t
	info.usage = tag2
	info.is_switch = tag_parser_is_switch(tag1)
	info.is_option, info.default_str = tag_parser_has_default(tag1)
	if info.is_option && len(info.default_str) > 0 {
		u := StringConverter(info.default_str)
		info.default_val, err = u.ToType(t)
		if err != nil {
			info.default_val = nil
		}
	}
	return
}

func (i1 *flag_info) same_as(i2 *flag_info) bool {
	return i1.is_switch == i2.is_switch && i1.is_option == i2.is_option && i1.default_str == i2.default_str
}

func (i *flag_info) same_as_def(tag string) bool {
	if i.is_switch != tag_parser_is_switch(tag) {
		return false
	}
	o, d := tag_parser_has_default(tag)
	return i.is_option == o && i.default_str == d
}

func NewParser(t ...reflect.Type) (*CommandParser, error) {
	ret := &CommandParser{}
	ret.setters = make(map[string]func(reflect.Value) error)
	return ret, ret.Bind(t...)
}

func (cp *CommandParser) Parse(cmd ...[]string) error {
	use := os.Args
	if len(cmd) > 0 {
		use = cmd[0]
	}
	var cur *flag_info = nil
	for _, arg := range use {
		if cur != nil {
			u := StringConverter(arg)
			if v, err := u.ToType(cur.flag_type); err != nil {
				return errors.New("命令行参数 " + arg + " 的类型不正确，期望类型是 " + cur.flag_type.String())
			} else {
				cur.value = v
				cur.parsed = true
				cur = nil
			}
			continue
		}
		if use, exist := cp.flags[arg]; exist {
			cur = use
			if cur.parsed {
				// TODO: 此处可以根据需要检查命令行参数是否可以重复
			}
			if cur.is_switch {
				v := reflect.ValueOf(true)
				cur.value = &v
			}
		} else {
			return errors.New("未知的命令行参数 " + arg)
		}
	}
	for _, i := range cp.flags {
		if i.value == nil && i.default_val == nil {
			buf := &bytes.Buffer{}
			buf.WriteString("缺少必须的命令行参数: ")
			buf.WriteString(i.flags[0])
			if len(i.flags) > 1 {
				buf.WriteString("(")
				buf.WriteString(strings.Join(i.flags[1:], ", "))
				buf.WriteString(")")
			}
			return errors.New(buf.String())
		}
	}
	cp.parsed = true
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
	t := reflect.TypeOf(out)
	if t.Kind() != reflect.Ptr {
		return errors.New("只能获取数据到结构指针，类型 " + t.String() + " 不能获取数据")
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return errors.New("只能获取数据到结构类型中，类型 " + t.String() + " 不能获取数据")
	}
	v := reflect.ValueOf(out).Elem()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		name := t.String() + "." + f.Name
		if setter, exist := cp.setters[name]; exist {
			setter(v)
		} else {
			return errors.New("没有对应结构 " + t.Name() + " 中 " + f.Name + " 字段的获取方法")
		}
	}
	return nil
}

func (cp *CommandParser) GetFlag(flag string, out interface{}) error {
	if info, exist := cp.flags[flag]; exist {
		v := info.value
		if v == nil {
			v = info.default_val
		}
		out_type := reflect.ValueOf(out)
		if out_type.Kind() != reflect.Ptr {
			return errors.New("不支持获取数据到类型 " + out_type.Type().Name() + " 中，数据只能保存到指针类型中")
		}
		out_type = out_type.Elem()
		out_type.Set(*v)
	} else {
		return errors.New("没有绑定命令行参数 " + flag)
	}
	return nil
}

func (cp *CommandParser) make_setter(f reflect.StructField, info *flag_info) func(reflect.Value) error {
	ret := func(out reflect.Value) error {
		v := info.value
		if v == nil {
			v = info.default_val
		}
		fld := out.FieldByName(f.Name)
		fld.Set(*v)
		return nil
	}
	return ret
}

func (cp *CommandParser) add_setter(t reflect.Type, f reflect.StructField, info *flag_info) {
	name := t.String() + "." + f.Name
	if _, exist := cp.setters[name]; exist {
		return
	}
	cp.setters[name] = cp.make_setter(f, info)
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
	info, err := new_tag_info(f.Type, tag1, tag2)
	if err != nil {
		return nil
	}
	info.flags = tag_parser_names(f.Name, tag1)
	for _, n := range info.flags {
		if i, exist := cp.flags[n]; exist {
			if !i.same_as_def(tag1) {
				return errors.New("命令行参数 " + n + " 重复定义，之前定义为类型 " + i.flag_type.String() + " ，新定义的类型是 " + f.Type.String())
			}
		} else {
			cp.flags[n] = info
		}
	}
	cp.add_setter(t, f, info)
	return nil
}

func (cp *CommandParser) proc_type(t reflect.Type) error {
	struct_type := t
	if struct_type.Kind() == reflect.Ptr {
		struct_type = struct_type.Elem()
	}
	if struct_type.Kind() != reflect.Struct {
		return errors.New("抱歉，只支持Struct")
	}
	if cp.flags == nil {
		cp.flags = make(map[string]*flag_info)
	}
	for _, tmp := range cp.types {
		if t == tmp {
			return nil
		}
	}
	// 此处处理t的Tags
	for i := 0; i < struct_type.NumField(); i++ {
		f := struct_type.Field(i)
		if tag, exist := f.Tag.Lookup("cmd"); exist {
			var usage string
			if usage, exist = f.Tag.Lookup("usage"); !exist {
				usage = ""
			}
			if err := cp.proc_field(struct_type, f, tag, usage); err != nil {
				return err
			}
		}
	}
	cp.types = append(cp.types, struct_type)
	return nil
}
