package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql" //千万不要忘记导入 默认会执行init初始化一些操作的
	"github.com/go-xorm/xorm"
	"xorm.io/core"
)

func main() {
	engine, err := xorm.NewEngine("mysql", "root:root@/elmcms?charset=utf8")
	if err != nil {
		panic(err.Error())
	}
	//设置名称映射规则
	//第一种：驼峰式   表明为：user_table  字段为： user_id  user_name ....
	engine.SetMapper(core.SnakeMapper{})
	engine.Sync2(new(UserTable))

	//第二种：表明为：studenttable 字段为：Id StudentName StudentAge StudentSex .....
	engine.SetMapper(core.SameMapper{})
	engine.Sync2(new(StudentTable))

	//第三种： 表名为：person_table 字段为：id person_name person_age .....
	engine.SetMapper(core.GonicMapper{})
	engine.Sync2(new(PersonTable))

	//注意：engine.SetMapper(core.SameMapper{})  这种形式创建的表示无法判断出数据是否为空或者表是否存在的！

	//判断一个表当中内容是否为空
	personEmpty, err := engine.IsTableEmpty(new(PersonTable))
	if err != nil {
		panic(err.Error())
	}
	if personEmpty {
		fmt.Println("人员表是空的!")
	} else {
		fmt.Println("人员表不为空!")
	}

	//判断表结构是否存在
	studentExist, err := engine.IsTableExist(new(StudentTable))
	if err != nil {
		panic(err.Error())
	}
	if studentExist {
		fmt.Println("学生表存在!")
	} else {
		fmt.Println("学生表不存在!")
	}

}

//用户表
type UserTable struct {
	UserId   int64  `xorm:"pk autoincr"` //用户id  主键
	UserName string `xorm:"varchar(32)"` //用户名称
	UserAge  int64  `xorm:"default 1"`   //用户年龄
	UserSex  int64  `xorm:"default 0"`   //用户性别
}

//学生表
type StudentTable struct {
	Id          int64  `xorm:"pk autoincr"` //主键 自增
	StudentName string `xorm:"varchar(24)"` //
	StudentAge  int    `xorm:"int default 0"`
	StudentSex  int    `xorm:"index"` //sex为索引
}

//人类表
type PersonTable struct {
	Id         int64     `xorm:"pk autoincr"`   //主键自增
	PersonName string    `xorm:"varchar(24)"`   //可变字符
	PersonAge  int       `xorm:"int default 0"` //默认值
	PersonSex  int       `xorm:"notnull"`       //不能为空
	City       CityTable `xorm:"-"`             //不映射该字段 那就不会在数据库里面创建该字段
}

type CityTable struct {
	CityName      string
	CityLongitude float32
	CityLatitude  float32
}

/*
xorm中对数据类型有自己的定义，具体的Tag规则如下，另Tag中的关键字均不区分大小写：

| name  | 当前field对应的字段的名称  |
| ---   | --- | --- |
| pk |  是否是Primary Key |
| name  | 当前field对应的字段的名称 |
| pk    | 是否是Primary Key       |
| autoincr | 是否是自增 |
| [not ]null 或 notnull | 是否可以为空 |
| unique | 是否是唯一 |
| index | 是否是索引 |
| extends | 应用于一个匿名成员结构体或者非匿名成员结构体之上
| - | 这个Field将不进行字段映射 |
| -> | Field将只写入到数据库而不从数据库读取 |
| <- | Field将只从数据库读取，而不写入到数据库 |
| created | Field将在Insert时自动赋值为当前时间 |
| updated | Field将在Insert或Update时自动赋值为当前时间 |
|deleted | Field将在Delete时设置为当前时间，并且当前记录不删除 |
| version | Field将会在insert时默认为1，每次更新自动加1 |
| default 0或default(0) | 设置默认值，紧跟的内容如果是Varchar等需要加上单引号 |
| json | 表示内容将先转成Json格式 |

*/

/*
### 字段映射规则
除了上述表名的映射规则和使用Tag对字段进行设置以外，基础的Go语言结构体数据类型也会对应到数据库表中的字段中，具体的一些数据类型对应规则如下：

| Go语言数据类型 | xorm 中的类型 |
| -------------| -------------|
| implemented Conversion | Text |
| int, int8, int16, int32, uint, uint8, uint16, uint32 | Int |
| int64, uint64 | BigInt |
| float32 | Float |
| float64 | Double |
| complex64, complex128 | Varchar(64) |
| []uint8 | Blob |
| array, slice, map except []uint8 | Text |
| bool | Bool |
| string | Varchar(255) |
| time.Time | DateTime |
| cascade struct | BigInt |
| struct | Text |
| Others | Text |

*/
