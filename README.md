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

**Enjoy 😊**
