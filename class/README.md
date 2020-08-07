# class
## usage

- 定义一些通用类，用于值可能为null的场景。
- 实现数据库的读写转化(Value, Scan)和json的转化接口(MarshalJSON, UnmarshalJSON)。
- 一些常用的函数，如Set等等

## 主要函数

- MarshalJSON, UnmarshalJSON ：用于json序列化
- Scan, Value： 用于sql，Value用于sql传参时驱动调用的。
- isValid：用于sql处理
- Set：值设置