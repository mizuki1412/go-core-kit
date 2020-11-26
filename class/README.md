# class

通用的抽象的类定义，可用于本核心库下，弥补原生类型的不足。

## usage

- 定义一些通用类，用于值可能为null的场景。
- 实现数据库的读写转化(Value, Scan)和json的转化接口(MarshalJSON, UnmarshalJSON)。
- 一些常用的类自带函数，如Set等等

## 主要函数

- MarshalJSON, UnmarshalJSON ：用于json序列化
- Scan, Value： 用于sql，Value用于sql传参时驱动调用的。
- isValid：用于sql处理
- Set：值设置

## 类库

### 基本类型

Bool, Decimal, Float64, Int32, Int64, String, Time

### Map

提供了一些map常用的函数接口。

- MapString：对postgres的jsonb格式做了适配。
- MapStringSync：提供了线程安全的MapString

### 数组类型

- ArrInt: 针对postgres.array的int数组封装，提供ToInt32Slice方法。
- ArrString: 针对postgres.array的string分装。
- MapStringArr：针对postgres jsonb的分装，array形式的jsonb

### queue

队列

### file

http上传文件流 