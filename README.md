# 改进说明（improvement）

原来的sprot只支持指针类型的变量，要通过sproto.Int(ptr *int)类似的方法赋值访问，非常麻烦。

我猜测云风大可能是想在nil的时候节省一些字节的开销。

这个库增强了Encode和Decode的能力，现在可以直接处理int string bool的值类型。

同时，使用相同位置及配置的值类型及对应指针，在tag一致的情况下，且指针均不为nil的情况下，编码结果是相等的，例如：

``` golang
type WithPtr struct{
	ID *int `sproto:"integer,0,name=ID"`
}

type WithVal struct{
	ID int `sproto:"integer,0,name=ID"`
}

```

如果WithPtr.ID=nil，则使用WithVal对编码进行Decode之后，WithVal.ID会等于0。

> 个人倾向尽可能不要使用nil作为值，作为一个用于跨平台的编码，nil容易在不同平台上产生不同的解析结果，极易产生歧义。

> 同样，个人不建议使用值类型的Struct，所以不支持了。

更多的实现效果请参考encode_test.go中的例子。

# gosproto
[sproto](https://github.com/cloudwu/sproto)'s encoder and decoder in golang.

# type map
sproto type      | golang type
---------------- | -------------------------------------------------
string           | \*string, []byte, string
integer          | \*int8, \*uint8, \*int16, \*uint16, \*int32, \*uint32, \*int64, \*uint64, \*int, \*uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64, int, uint
boolean          | \*bool, bool
object           | \*struct
array of string  | []string
array of integer | []int8, []uint8, []int16, []uint16, []int32, []uint32, []int64, []uint64, []int, []uint
array of boolean | []bool
array of object  | []\*struct

# schema
You can define go struct corresponding to sproto schema directly as examples in all test cases.
Or use [sprotodump](https://github.com/lvzixun/sprotodump) to change sproto schema to go file.

# test
```
go test github.com/xjdrew/gosproto
```

# benchmark
```
$ go test -bench . github.com/xjdrew/gosproto
PASS
BenchmarkEncode-4         300000      4122 ns/op
BenchmarkDecode-4         300000      5417 ns/op
BenchmarkEncodePacked-4122  300000      4712 ns/op
BenchmarkDecodePacked-4712  200000      6267 ns/op
ok      github.com/xjdrew/gosproto5.752s
```
