# service

## sqlkit

### 注意

- **注意 commit: 如果事务中第一句是select语句，commit将会出错, 错误提示 parse C 等。**
- 在rows遍历时，注意close，特别是有级联查询存在时，如果不close将会占用连接。
- bean struct中如果没有db标签，则不会被通用接口insert/update
- sqlx的`missing destination name sth in sth`，是查询出来的字段和类字段不符，在select中限定字段即可。
- update set时：`Set("extend",squirrel.Expr("'{}'::jsonb"))` or `Set("extend","{}")`
- class.mapString在插入数据库时将用jsonb格式，并且不是完全替换，而是merge的方式(```coalesce(extend, '{}'::jsonb) || '$param'::jsonb```)。如果要删除其中的key，需要设置key为null。 merge时只会merge顶层的keys。

###
数据库驱动：
- postgres: _ "github.com/lib/pq"
- mysql: _ "github.com/go-sql-driver/mysql"

### demo
```go
func (dao *Dao) UpdateConfirm(id int64){
    sql, args, err := dao.Builder().Update(meta.GetTableName(dao.Schema)).Set("extend",squirrel.Expr("jsonb_set(extend, '{confirm}','true',true)")).Where("id=?",id).ToSql()
    if err != nil {
        panic(exception.New(err.Error()))
    }
    dao.Exec(sql, args...)
}

func (dao *Dao) List(dTypes []string) []model.AlarmMsg {
	// todo
	builder := dao.Builder().Select(meta.GetColumnsWithPrefix("msg")).From(dao.GetTableD("alarm_msg msg") + dao.GetTableD("device_type_info info")).Where("msg.deviceType=info.id").OrderBy("msg.deviceType, msg.id")
	if dTypes!=nil && len(dTypes) > 0 {
		flag, arg := pghelper.GenUnnestString(dTypes)
		builder = builder.Where("msg.deviceType in "+flag, arg)
	}
	sql, args := builder.MustSql()
	return dao.scan(sql, args)
}

func (dao *Dao) ListId(dType []string) []string {
	builder := dao.Builder().Select("id").From(meta.GetTableName(dao.Schema)).Where("off=?", false).OrderBy("id")
	if dType!=nil && len(dType)>0{
		builder = pghelper.WhereUnnestInt(builder,"id in ", dType)
	}
	sql, args := builder.MustSql()
	rows := dao.Query(sql, args...)
	list := make([]string, 0, 5)
	defer rows.Close()
	for rows.Next() {
		ret, err := rows.SliceScan()
		if err != nil {
			panic(exception.New(err.Error()))
		}
		list = append(list, ret[0].(string))
	}
	return list
}
```
