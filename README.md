# 氖-Ne

[![golang](https://img.shields.io/badge/Language-Go-green.svg?style=flat)](https://golang.org)
[![pkg.go.dev](https://img.shields.io/badge/dev-reference-007d9c?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/noble-gase/ne)
[![MIT](http://img.shields.io/badge/license-MIT-brightgreen.svg)](http://opensource.org/licenses/MIT)

[氖-Ne] Go开发工具包

## 获取

```shell
go get -u github.com/noble-gase/ne
```

| 模块      | 说明                                                                                         |
| --------- | -------------------------------------------------------------------------------------------- |
| array     | 切片常用操作                                                                                 |
| conv      | 类型转换                                                                                     |
| coord     | 距离、方位角、经纬度与平面直角坐标系的相互转化                                               |
| crypts    | 封装 Crypto 常用方法，支持： `AES` 和 `RSA`                                                  |
| curd      | 基于 [`Jet`](https://github.com/go-jet/jet) 封装的 curd 操作                                 |
| helper    | 常用的辅助方法合集，包含：HTTP、IP、VersionCompare 等                                        |
| httpzip   | 远程获取 `ZIP` 压缩包中的文件内容                                                            |
| hush      | 封装 Hash 常用方法                                                                           |
| images    | 图片处理，如：缩略图、裁切、标注等                                                           |
| leveltree | 基于泛型的菜单和组织单位等分类层级树                                                         |
| mutex     | 基于 Redis 的分布式锁                                                                        |
| protos    | 实现 `url.Values` 和 `proto.Message` 的相互转换                                              |
| redix     | 基于 `singleflight` 封装 Redis 常用操作                                                      |
| retry     | 重试操作                                                                                     |
| sqls      | 包含DB初始化和事务等封装                                                                     |
| steps     | 分批次处理切片                                                                               |
| values    | 用于处理 `k-v` 格式化的场景，如：生成签名串等                                                |
| verify    | 验证器（基于 [`validator`](https://github.com/go-playground/validator)）支持汉化和自定义规则 |

**Enjoy 😊**
