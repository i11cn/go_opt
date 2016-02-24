package option

import (
	"bytes"
	"fmt"
	"github.com/bitly/go-simplejson"
)

func (o *Options) make_path(key ...string) *Options {
	ret := o
	for _, k := range key {
		if c, ok := ret.child[k]; ok {
			ret = c
		} else {
			use := NewOptions()
			use.name = k
			ret.child[k] = use
			ret = use
		}
	}
	return ret
}

func (o *Options) marshal_json(buf *bytes.Buffer) error {
	if o.value == nil && len(o.child) == 0 {
		return nil
	}
	buf.WriteString(fmt.Sprintf("\"%s\":", o.name))
	if len(o.child) > 0 {
		buf.WriteString("{")
		var comma bool = false
		for _, c := range o.child {
			if comma {
				buf.WriteString(", ")
			}
			c.marshal_json(buf)
			comma = true
		}
		buf.WriteString("}")
	} else {
		switch v := o.value.(type) {
		case string:
			buf.WriteString(fmt.Sprintf("\"%s\"", v))
		default:
			buf.WriteString(fmt.Sprint(o.value))
		}
	}
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
