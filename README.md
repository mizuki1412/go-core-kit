
# go-core-kit

toolkit for golang projects

`go get github.com/mizuki1412/go-core-kit@v2.0.0`

详细说明见 doc/go-core-kit-doc.md

# TODO

  later：
- [ ] 接口加密：rsa+aes对接口加密(前端生成AES密钥，用rsa私钥将AES密钥加密，放在header.encript 传给到后端，后端用rsa公钥解密后获取到AES密钥，解密数据流。
  对参数值进行加密，同时aes密钥中增加时间变量)
- [ ] https://github.com/unrolled/render 模板渲染
- [ ] mqtt subscribe 中如果执行太久，会重复执行subscribe？ 暂时用go fun处理
- [ ] sql base mapper: 增加多数据库适配
- [ ] 改进：关于子查询的优化。where in 等
- [x] 性能：每次 dao 都会重新解析 model
- [ ] 重构：mod user

# 1.0 升级 2.0 指南

- cmd 重构，改用新的 cli 包
- class 基础类重构，推荐用 NewXX() 或 NXX() 新建
- class.Decimal 指针改为值类型
- class.time 用回默认的 nullTime，观察 scan 的时区是否有问题
- model定义时sql的标签注意：table、logicDel
- sqlkit 重构，参考 `doc/goland-live-templates.md`，重新生成dao模板代码
- dao 函数中带 args 参数的，都改用[]any，一致性
- dao 采用链式操作 (参考userdao)，提供了一些基础的封装函数
- dao 的OrderBy注意，一个字段一个
- dao 级联时注意是否忽略删除标记获取，因为默认是取未删除的
- dao resultType去掉，在new时设置，不再动态指定
- rest 取消 session，全面改用 jwt，见 jwtkit 说明; 也保留了cookie
- rest 默认返回值改变：code=0 表示 ok，code=401 表示未认证（也反映到 httpcode 中）
- rest swagger接口定义方式改变, 配置方式改为functional options
- rest authup改为authjwt
- logkit 基于slog
- 配置参数修改：openapi