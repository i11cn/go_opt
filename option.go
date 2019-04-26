package option

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type (
	Item string

	CommandParser struct {
		types   []reflect.Type
		flags   map[string]*flag_info
		setters map[string]func(reflect.Value) error
		parsed  bool
	}
)

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
				return fmt.Errorf("命令行参数 %s 的类型不正确，期望类型是 %s", arg, cur.flag_type.String())
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
			return fmt.Errorf("未知的命令行参数 %s", arg)
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
			return fmt.Errorf(buf.String())
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
		return fmt.Errorf("只能获取数据到结构指针，类型 %s 不能获取数据", t.String())
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("只能获取数据到结构类型中，类型 %s 不能获取数据", t.String())
	}
	v := reflect.ValueOf(out).Elem()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		name := t.String() + "." + f.Name
		if setter, exist := cp.setters[name]; exist {
			setter(v)
		} else {
			return fmt.Errorf("没有对应结构 %s 中 %s 字段的获取方法", t.Name(), f.Name)
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
			return fmt.Errorf("不支持获取数据到类型 %s 中，数据只能保存到指针类型中", out_type.Type().Name())
		}
		out_type = out_type.Elem()
		out_type.Set(*v)
	} else {
		return fmt.Errorf("没有绑定命令行参数 %s", flag)
	}
	return nil
}

func (cp *CommandParser) add_setter(t reflect.Type, f reflect.StructField, info *flag_info) {
	name := t.String() + "." + f.Name
	if _, exist := cp.setters[name]; exist {
		return
	}
	cp.setters[name] = make_setter(f, info)
}

func (cp *CommandParser) proc_field(t reflect.Type, f reflect.StructField, tag1, tag2 string) error {
	if !valid_type(f.Type) {
		return fmt.Errorf("目前还不支持所指定的类型： %s", f.Type.String())
	}
	info, err := new_tag_info(f.Type, tag1, tag2)
	if err != nil {
		return nil
	}
	info.flags = tag_parser_names(f.Name, tag1)
	for _, n := range info.flags {
		if i, exist := cp.flags[n]; exist {
			if !i.same_as_def(tag1) {
				return fmt.Errorf("命令行参数 %s 重复定义，之前定义为类型 %s ，新定义的类型是 %s", n, i.flag_type.String(), f.Type.String())
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
		return fmt.Errorf("抱歉，只支持Struct")
	}
	if cp.flags == nil {
		cp.flags = make(map[string]*flag_info)
	}
	for _, tmp := range cp.types {
		if t == tmp {
			return nil
		}
	}
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
