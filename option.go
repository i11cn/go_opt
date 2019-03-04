package option

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
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

	CommandParser struct {
		options map[string]Item
		types   []reflect.Type
	}
)

func Parse(cmd ...[]string) *CommandParser {
	ret := &CommandParser{}
	ret.parse(cmd...)
	return ret
}

func (cp *CommandParser) parse(cmd ...[]string) {
	cp.options = make(map[string]Item)
	use := os.Args
	if len(cmd) > 0 {
		use = cmd[0]
	}
	for _, arg := range use {
		fmt.Println(arg)
	}
}

func (cp *CommandParser) Bind(t ...reflect.Type) {
}

func (cp *CommandParser) Usage() string {
	ret := &bytes.Buffer{}
	return ret.String()
}

func (cp *CommandParser) Get(out interface{}) error {
	return nil
}

func (cp *CommandParser) proc_type(t reflect.Type) error {
	return nil
}
