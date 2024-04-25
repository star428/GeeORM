# GeeORM

## day0-ORM基础

> 对象关系映射（Object realtional Mapping）是通过使用描述对象和数据库之间映射的元数据，将面向对象语言程序中的对象自动持久化到关系数据库中。

映射关系如下：

| 数据库                  | 面向对象编程语言          |
| ----------------------- | ------------------------- |
| 表（table）             | 类（class，struct）       |
| 记录/行（record/row）   | 类的具体实现/值（object） |
| 字段/列（field/column） | 对象属性（attribute）     |

具体例子如下：

```sql
create table user (name text, age integer);
insert into user(name, age) values ("tom", 18);
select * from user;
```

假如我们使用ORM框架，我们可以这样写：

```go
type user struct{
	name string
	age int
}

orm.createTable(&user{}) // 根据结构体名和结构体中的属性创建表
orm.save(&user{"tom", 18}) // 插入一个记录

var users []user

orm.Find(&users) // 从数据库中读取所有行到slice中
```

ORM框架就是一个对象和数据库之间的桥梁，可以避免写繁琐的sql语言，通过操作具体的对象，就可以完成对关系型数据库的操作。

我们首先关注第一行的问题：

我们要通过任意类型(struct)的指针，来获取struct的名称和属性的名称。我们使用反射来获取：

```go
typ := reflect.Indirect(reflect.ValueOf(&Account{})).Type() // 获取struct指针的名字

for i := 0; i < typ.NumField(); i++{
	field := typ.field(i)
	name := field.Name // 获取每个属性的名字
}
```

我们还需要关注这些问题：

1. 不同关系型数据库的sql语句是有区别的，ORM框架如何避免这些问题？
2. 对象的字段发生改变，数据库表的结构能够自动更新（变列名），能否支持自动迁移（migrate）？
3. 数据库支持的功能很多，例如事务(transaction)，ORM框架能否实现这些（事务的ACID）？

## day1-database/sql基础

* SQLite的基本操作（一般的sql语句）
* Go标准库database/sql连接并操作SQLite数据库并做简单封装

```go
package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // 做初始化操作并不导入
)

func main() {
	// 连接到数据库，第一个参数为驱动名称，第二个参数为数据库路径，没有会新建一个
	db, _ := sql.Open("sqlite3", "../database/gee.db") 

	defer func() {
		_ = db.Close()
	}()

	// 执行语句，查询语句不会返回结果，通常查询使用Query()，QueryRow()
	// 第一个参数是sql语句，第二个是占位符？对应的值
	db.Exec("DROP TABLE IF EXISTS User;")
	db.Exec("CREATE TABLE User(Name text);")
	r, err := db.Exec("INSERT INTO User(`Name`) values(?), (?)", "Tom", "Sam")
	if err == nil {
		affected, _ := r.RowsAffected()
		println(affected)
	}

	// 只返回一行
	row := db.QueryRow("SELECT Name FROM User LIMIT 1")
	var name string

	// 接收一个或多个指针作为参数（比如一个slice指针组a []*int，使用a...即可）
	if err := row.Scan(&name); err == nil {
		println(name)
	}
}

```

### 实现一个简单的log库

> 不使用原生log库的原因是log标准库没有日志分级，不打印文件和行号，这意味着我们很难快速知道哪里出了错误

这个log库具有以下特性

* 支持日志分级（info，error，disabled）
* 不同日志显示不同颜色
* 显示打印日志代码所对应的文件名和行号

```go
errorLog := log.New(os.Stdout, "\033[31m[error]\033[0m ", log.LstdFlags|log.Lshortfile)
// 创建新的log，后面就可以使用Println/Printf方法
```

### 核心session部分

可以理解为创建会话（一次会话连接到关闭的过程中可能有多个事务）

session只需要存储三个变量：

```go
// Session is a struct that contains a database pointer, sql language, and sql variables
type Session struct {
	// db is a database pointer
	db *sql.DB
	// sql is sql language such as "select * from user where id = ?"
	sql strings.Builder
	// sqlVars is the variables in the sql language, such as 1, 2, 3
	sqlVars []any
}
```

接下来封装 `Exec()`执行, `Query()`查询, `QueryRow()`单行查询三个方法

目的是：

* 统一打印日志（包括执行的SQL语句和错误日志）
* 每个操作执行结束后清空 `(s *Session).sql`和 `(s *Session).sqlVars`，这样session可以复用，开启一次会话，执行多个sql语句。

### 核心结构Engine

功能：

* 连接到数据库并载入driver（`func NewEngine`）
* 关闭与数据库的连接（`func Close`）
* 创建一个新会话来执行操作
