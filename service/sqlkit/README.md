# service

## sqlkit

- **注意 commit: 如果事务中第一句是select语句，commit将会出错, 错误提示 parse C 等。**
- 在rows遍历时，注意close，特别是有级联查询存在时，如果不close将会占用连接。
- bean struct中如果没有db标签，则不会被通用接口insert/update
- sqlx的`missing destination name sth in sth`，是查询出来的字段和类字段不符，在select中限定字段即可。
- update set时：`Set("extend",squirrel.Expr("'{}'::jsonb"))` or `Set("extend","{}")`

### demo
```go
func (dao *Dao) UpdateConfirm(id int64){
	sql, args, err := sqlkit.Builder().Update(dao.GetTableD("device_alarm_record")).Set("extend",squirrel.Expr("jsonb_set(extend, '{confirm}','true',true)")).Where("id=?",id).ToSql()
	if err != nil {
		panic(exception.New(err.Error()))
	}
	dao.Exec(sql, args...)
}

func (dao *Dao) List(dTypes []string) []model.AlarmMsg {
	builder := sqlkit.Builder().Select("msg.*").From(dao.GetTableD("alarm_msg msg") + dao.GetTableD("device_type_info info")).Where("msg.deviceType=info.id").OrderBy("msg.deviceType, msg.id")
	if dTypes!=nil && len(dTypes) > 0 {
		flag, arg := pghelper.GenUnnestString(dTypes)
		builder = builder.Where("msg.deviceType in "+flag, arg)
	}
	sql, args := builder.MustSql()
	return dao.scan(sql, args)
}

func (dao *Dao) ListId(dType []string) []string {
	builder := sqlkit.Builder().Select("id").From(dao.GetTableD("device")).Where("off=?", false).OrderBy("id")
	if dType!=nil && len(dType)>0{
		flag,arg := pghelper.GenUnnestString(dType)
		builder = builder.Where("type in "+flag, arg...)
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
