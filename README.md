
# go-core-kit

toolkit for golang projects

`go get github.com/mizuki1412/go-core-kit@v1.4.5`

# init
本库使用的入口，以及配置参数信息相关的绑定函数

# class
通用的一些类的封装和定义

# library
通用的工具库

# service
通用的服务库

# service-third
针对第三方服务接口的封装

# pc
应用于pc端，web+go的模式，go作为基座的一些封装。

# mod
公用的业务模块

# tool-local
本地使用的一些小工具

# 其他

## 可替换的函数

# TODO

- [ ] https://github.com/gin-contrib/sessions/pull/148 等官方更新
- [ ] https://github.com/unrolled/render 模板渲染

## 旧项目升级注意

- ctx.SessionGetUser().(*model2.User); 
- interface{}->any
- 配置文件：aliyun
- dao