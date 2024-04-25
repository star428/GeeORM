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
