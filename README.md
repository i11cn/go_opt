# go_opt

### HOW TO
- 预习:
  - option包使用了struct的tag来定义命令行参数
  - tag的格式如下: **cmd:"[flag名称...][,switch|option,默认参数]"**
  - 其中各个字段的含义解释如下：
    - flag名称，是出现在命令行中，解释为本flag的字符串，可以不带前缀的短横'-'
    - flag名称可以有多个，如果不写flag名称，则自动解析为当前结构中字段的名称
    - tag中的名称不带'-'时，按照单字符增加单个'-'，多字符增加双'-'的规则
    - 如果要指定不同的'-'前缀，可以明确写在tag中
    - switch表示该flag是个开关，即之后不会跟flag的参数，如果不指定为switch，命令行中该flag之后的参数，会作为该flag的参数
    - switch同时隐含了option，同时要求flag类型为bool，不得为其他类型
    - option表示该flag是可选的，不一定会出现在命令行中，但同时必须在option之后增加默认参数
  - 还有一个tag: **usage:"flag说明文字"**
    - usage的使用会在usage完善的时候用到，即可以通过Parser的Usage方法，输出一段命令行参数说明

- 创建Parser:

**option.NewParser**

示例：

```
opt, _ := option.NewParser()
```

- 绑定

**CommandParser.Bind**

示例:

```
type flag_value struct {
    Host string `cmd:"H,option,localhost"`
    Port int    `cmd:"P"`
}
t := flag_value{}
opt.Bind(reflect.TypeOf(t))
```

**可以把NewParser和Bind结合起来**

示例:

```
type flag_value struct {
    Host string `cmd:"H,option,localhost"`
    Port int    `cmd:"port,P"`
}
t := flag_value{}
opt, _ := option.NewParser(reflect.TypeOf(t))
```

**CommandParser.Parse**

示例:

```
opt.Parse()
```

或者自己输入数组来解析:

```
opt.Parse([]string{"-H", "192.168.1.1", "--port", "12345"})
```

**CommandParser.Get**

示例:

```
type flag_value struct {
    Host string `cmd:"H,option,localhost"`
    Port int    `cmd:"P"`
}
t := flag_value{}
opt := option.NewParser(reflect.TypeOf(t))
opt.Get(&t)
```

### 注意

每一步都会有error的返回值，应该检查该返回值，确保正确之后进行下一步操作
- NewParser如果不带参数，不会产生错误，如果有参数，等同于调用了一次Bind，会有错误返回
- Bind如果tag解析出错，tag设置出错等，都会返回错误
- Parser如果命令行参数输入的不正确，会返回错误
- Get获取时的数据类型不符，会返回错误

### TODO
1. 增加命令行参数，非flag的解析
2. 增加Usage的显示
3. 增加对重复flag的处理，处理为数组格式
4. 增加更多的配置选项，比如是否接手重复flag等
5. 增加flag的顺序检查
6. 增加struct之外的单个flag的设置和获取方法
7. 此时准备发布成1.0.0
