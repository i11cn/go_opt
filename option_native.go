package option

import (
	"github.com/bitly/go-simplejson"
)

func (o *Options) dump() opt_dump {
	ret := opt_dump{}
	if o.value != nil {
		ret["."] = o.value
	}
	for k, v := range o.items {
		ret[k] = v.dump()
	}
	return ret
}

func (o *Options) parse_json_object(j *simplejson.Json) error {
	if s, err := j.String(); err == nil {
		o.value = s
	}
	cur := j.GetIndex(0)
	for idx := 1; cur != nil; {
		cur = j.GetIndex(idx)
		idx++
	}
	return nil
}
func (o *Options) parse_json_string(js string) error {
	j, err := simplejson.NewJson([]byte(js))
	if err != nil {
		return err
	}
	return o.parse_json_object(j)
	return nil
}

func (o *Options) parse_command_line() error {
	return nil
}
