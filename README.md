
# go-core-kit

toolkit for golang projects

`go get github.com/mizuki1412/go-core-kit@v1.5.5`

详细说明见 doc/go-core-kit-doc.md

# 其他

## 可替换的函数

# TODO

- [ ] 测试当 redis 没配置时，session 是否正常; 用jwt替代
- [ ] class中有些类型是写死pg的，对其他数据库不友好： array，jsonarray
- [ ] class 和 model 中的接受函数测试下，可能在 scan 时有问题
- [ ] class.time scan test
- [ ] page()

  later：
- [ ] https://github.com/unrolled/render 模板渲染
- [ ] swagger response 的定义
- [ ] mqtt subscribe 中如果执行太久，会重复执行subscribe？ 暂时用go fun处理
- [ ] sql base mapper: 增加多数据库适配
- [ ] 改进：关于子查询的优化。where in 等
- [ ] 性能：每次 dao 都会重新解析 model

# 1.0 升级 2.0 指南

- sqlkit 重构，参考 `doc/goland-live-templates.md`，重新生成dao模板代码
- cmd 重构，改用新的 cli 包
- class 基础类重构，推荐用 NewXX() 或 NXX() 新建
- class.Decimal 指针改为值类型
- dao 函数中带 args 参数的，都改用[]any，一致性