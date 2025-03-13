# 氖-Ne

[![golang](https://img.shields.io/badge/Language-Go-green.svg?style=flat)](https://golang.org)
[![pkg.go.dev](https://img.shields.io/badge/dev-reference-007d9c?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/noble-gase/ne)
[![MIT](http://img.shields.io/badge/license-MIT-brightgreen.svg)](http://opensource.org/licenses/MIT)

[氖-Ne] Go开发实用库

## 获取

```shell
go get -u github.com/noble-gase/ne
```

## 包含

- array - 数组的常用方法
- hashes - 封装便于使用
- crypts - 封装便于使用(支持 AES & RSA)
- validator - 支持汉化和自定义规则
- sqls - 基于 `sqlx` 的轻量SQLBuilder
- linklist - 一个并发安全的双向列表
- errgroup - 基于官方版本改良，支持并发协程数量控制
- values - 用于处理 `k-v` 格式化的场景，如：生成签名串等
- coord - 距离、方位角、经纬度与平面直角坐标系的相互转化
- timewheel - 简单实用的单层时间轮(支持一次性和多次重试任务)
- images - 图片处理，如：缩略图、裁切、标注等
- 基于 Redis 的分布式锁
- 基于泛型的无限菜单分类层级树
- 实用的辅助方法：IP、file、time、slice、string、version compare 等

> ⚠️ 注意：如需支持协程并发复用的 `errgroup` 和 `timewheel`，请使用 👉 [氙-Xe](https://github.com/noble-gase/xe)

## SQL Builder

> ⚠️ 目前支持的特性有限，复杂的SQL（如：子查询等）还需自己手写

```go
builder := sqls.New(*sqlx.DB, func(ctx context.Context, query string, args ...any) {
    fmt.Println(query, args)
})
```

### 👉 Query

```go
ctx := context.Background()

type User struct {
    ID     int    `db:"id"`
    Name   string `db:"name"`
    Age    int    `db:"age"`
    Phone  string `db:"phone,omitempty"`
}

var (
    record User
    records []User
)

builder.Wrap(
    sqls.Table("user"),
    sqls.Where("id = ?", 1),
).One(ctx, &record)
// SELECT * FROM user WHERE (id = ?)
// [1]

builder.Wrap(
    sqls.Table("user"),
    sqls.Where("name = ? AND age > ?", "og", 20),
).All(ctx, &records)
// SELECT * FROM user WHERE (name = ? AND age > ?)
// [og 20]

builder.Wrap(
    sqls.Table("user"),
    sqls.Where("name = ?", "og"),
    sqls.Where("age > ?", 20),
).All(ctx, &records)
// SELECT * FROM user WHERE (name = ?) AND (age > ?)
// [og 20]

builder.Wrap(
    sqls.Table("user"),
    sqls.WhereIn("age IN (?)", []int{20, 30}),
).All(ctx, &records)
// SELECT * FROM user WHERE (age IN (?, ?))
// [20 30]

builder.Wrap(
    sqls.Table("user"),
    sqls.Select("id", "name", "age"),
    sqls.Where("id = ?", 1),
).One(ctx, &record)
// SELECT id, name, age FROM user WHERE (id = ?)
// [1]

builder.Wrap(
    sqls.Table("user"),
    sqls.Distinct("name"),
    sqls.Where("id = ?", 1),
).One(ctx, &record)
// SELECT DISTINCT name FROM user WHERE (id = ?)
// [1]

builder.Wrap(
    sqls.Table("user"),
    sqls.LeftJoin("address", "user.id = address.user_id"),
    sqls.Where("user.id = ?", 1),
).One(ctx, &record)
// SELECT * FROM user LEFT JOIN address ON user.id = address.user_id WHERE (user.id = ?)
// [1]

builder.Wrap(
    sqls.Table("address"),
    sqls.Select("user_id", "COUNT(*) AS total"),
    sqls.GroupBy("user_id"),
    sqls.Having("user_id = ?", 1),
).All(ctx, &records)
// SELECT user_id, COUNT(*) AS total FROM address GROUP BY user_id HAVING (user_id = ?)
// [1]

builder.Wrap(
    sqls.Table("user"),
    sqls.Where("age > ?", 20),
    sqls.OrderBy("age ASC", "id DESC"),
    sqls.Offset(5),
    sqls.Limit(10),
).All(ctx, &records)
// SELECT * FROM user WHERE (age > ?) ORDER BY age ASC, id DESC LIMIT ? OFFSET ?
// [20, 10, 5]

wrap1 := builder.Wrap(
    sqls.Table("user_1"),
    sqls.Where("id = ?", 2),
)

builder.Wrap(
    sqls.Table("user_0"),
    sqls.Where("id = ?", 1),
    sqls.Union(wrap1),
).All(ctx, &records)
// (SELECT * FROM user_0 WHERE (id = ?)) UNION (SELECT * FROM user_1 WHERE (id = ?))
// [1, 2]

builder.Wrap(
    sqls.Table("user_0"),
    sqls.Where("id = ?", 1),
    sqls.UnionAll(wrap1),
).All(ctx, &records)
// (SELECT * FROM user_0 WHERE (id = ?)) UNION ALL (SELECT * FROM user_1 WHERE (id = ?))
// [1, 2]

builder.Wrap(
    sqls.Table("user_0"),
    sqls.WhereIn("age IN (?)", []int{10, 20}),
    sqls.Limit(5),
    sqls.Union(
        builder.Wrap(
            sqls.Table("user_1"),
            sqls.Where("age IN (?)", []int{30, 40}),
            sqls.Limit(5),
        ),
    ),
).All(ctx, &records)
// (SELECT * FROM user_0 WHERE (age IN (?, ?)) LIMIT ?) UNION (SELECT * FROM user_1 WHERE (age IN (?, ?)) LIMIT ?)
// [10, 20, 5, 30, 40, 5]
```

### 👉 Insert

```go
ctx := context.Background()

type User struct {
    ID     int64  `db:"-"`
    Name   string `db:"name"`
    Age    int    `db:"age"`
    Phone  string `db:"phone,omitempty"`
}

builder.Wrap(Table("user")).Insert(ctx, &User{
    Name: "og",
    Age:  29,
})
// INSERT INTO user (name, age) VALUES (?, ?)
// [og 29]

builder.Wrap(sqls.Table("user")).Insert(ctx, sqls.M{
    "name": "og",
    "age":  29,
})
// INSERT INTO user (name, age) VALUES (?, ?)
// [og 29]
```

### 👉 Batch Insert

```go
ctx := context.Background()

type User struct {
    ID     int64  `db:"-"`
    Name   string `db:"name"`
    Age    int    `db:"age"`
    Phone  string `db:"phone,omitempty"`
}

builder.Wrap(Table("user")).BatchInsert(ctx, []*User{
    {
        Name: "og",
        Age:  20,
    },
    {
        Name: "og",
        Age:  29,
    },
})
// INSERT INTO user (name, age) VALUES (?, ?), (?, ?)
// [og 20 og 29]

builder.Wrap(sqls.Table("user")).BatchInsert(ctx, []sqls.M{
    {
        "name": "og",
        "age":  20,
    },
    {
        "name": "og",
        "age":  29,
    },
})
// INSERT INTO user (name, age) VALUES (?, ?), (?, ?)
// [og 20 og 29]
```

### 👉 Update

```go
ctx := context.Background()

type User struct {
    Name   string `db:"name"`
    Age    int    `db:"age"`
    Phone  string `db:"phone,omitempty"`
}

builder.Wrap(
    sqls.Table("user"),
    sqls.Where("id = ?", 1),
).Update(ctx, &User{
    Name: "og",
    Age:  29,
})
// UPDATE user SET name = ?, age = ? WHERE (id = ?)
// [og 29 1]

builder.Wrap(
    sqls.Table("user"),
    sqls.Where("id = ?", 1),
).Update(ctx, sqls.M{
    "name": "og",
    "age":  29,
})
// UPDATE user SET name = ?, age = ? WHERE (id = ?)
// [og 29 1]

builder.Wrap(
    sqls.Table("product"),
    sqls.Where("id = ?", 1),
).Update(ctx, sqls.M{
    "price": sqls.SQLExpr("price * ? + ?", 2, 100),
})
// UPDATE product SET price = price * ? + ? WHERE (id = ?)
// [2 100 1]
```

### 👉 Delete

```go
ctx := context.Background()

builder.Wrap(
    sqls.Table("user"),
    sqls.Where("id = ?", 1),
).Delete(ctx)
// DELETE FROM user WHERE id = ?
// [1]

builder.Wrap(sqls.Table("user")).Truncate(ctx)
// TRUNCATE user
```

### 👉 Transaction

```go
builder.Transaction(context.Background(), func(ctx context.Context, tx sqls.TXBuilder) error {
    _, err := tx.Wrap(
        sqls.Table("address"),
        sqls.Where("user_id = ?", 1),
    ).Update(ctx, sqls.M{"default": 0})
    if err != nil {
        return err
    }

    _, err = tx.Wrap(
        sqls.Table("address"),
        sqls.Where("id = ?", 1),
    ).Update(ctx, sqls.M{"default": 1})

    return err
})
```

**Enjoy 😊**
