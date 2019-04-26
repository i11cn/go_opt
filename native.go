package option

import (
	"reflect"
	"strings"
)

type (
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

func make_setter(f reflect.StructField, info *flag_info) func(reflect.Value) error {
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

func valid_type(t reflect.Type) bool {
	switch t.String() {
	case "string", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "float32", "float64", "bool":
		return true
	case "*string", "*int", "*int8", "*int16", "*int32", "*int64", "*uint", "*uint8", "*uint16", "*uint32", "*uint64", "*float32", "*float64", "*bool":
		return true
	}
	return false
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
