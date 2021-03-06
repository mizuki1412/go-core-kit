
## action_init
```
func Init(router *router.Router) {
	tag := "$tname$"
	r := router.Group("/rest/$tag$")
	r.Use(middleware.AuthUsernameAndPwd())
	{
		r.Post("/$name$", $name$).Swagger.Tag(tag).Summary("$summary$").Param($name$Params{})
	}
}

type $name$Params struct{
    Phone    string `description:"手机号" default:"" trim:"true"`
	Pwd      string `validate:"required"`
}
func $name$(ctx *context.Context){
    params := $name$Params{}
	ctx.BindForm(&params)
	
    ctx.JsonSuccess(nil)
}
```

## action
```
type $name$Params struct{
    Phone    string `description:"手机号" default:"" trim:"true"`
	Pwd      string `validate:"required"`
}
func $name$(ctx *context.Context){
    params := $name$Params{}
	ctx.BindForm(&params)
	
    ctx.JsonSuccess(nil)
}
```

## bean_extend
```
func (th *$struct$) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	th.Id = cast.ToInt32(value)
	return nil
}
func (th $struct$) Value() (driver.Value, error) {
    // todo 注意返回值类型
	return int64(th.Id), nil
}
```

## bean_extend_list
bean list的sort/filter/find功能
```
type $name$List []*$name$
func (l $name$List) Len() int           { return len(l) }
func (l $name$List) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l $name$List) Less(i, j int) bool { return l[i].Id.String < l[j].Id.String }
func (l $name$List) Filter(fun func(ele *$name$) bool) $name$List {
	arr:=make($name$List, 0, len(l))
	for _,e:=range l{
		if fun(e) {
			arr = append(arr, e)
		}
	}
	return arr
}
func (l $name$List) Find(fun func(ele *$name$) bool) *$name$ {
	for _, e := range l {
		if fun(e) {
			return e
		}
	}
	return nil
}
func (l $name$List) MapReduce(fun func(ele *$name$) interface{}) []interface{} {
	var results []interface{}
	for _, e := range l {
		results = append(results, fun(e))
	}
	return results
}
```

## dao
```
/// auto template
type Dao struct {
	sqlkit.Dao
}
const (
	ResultDefault byte = iota
	ResultNone
)
func New(schema string, tx ...*sqlkit.Dao) *Dao{
	dao:=&Dao{}
	dao.NewHelper(schema,tx...)
	return dao
}
func (dao *Dao) cascade(obj *$bean$) {
	switch dao.ResultType {
	case ResultDefault:
		// todo 注意校验nil
		// todo 如果没有级联，此函数删除
    case ResultNone:
		// todo 将外联的置为nil
	}
}
func (dao *Dao) scan(sql string, args []interface{}) []*$bean$ {
	rows := dao.Query(sql, args...)
	list := make([]*$bean$,0,5)
	defer rows.Close()
	for rows.Next() {
		m := &$bean${}
		err := rows.StructScan(m)
		if err != nil {
			panic(exception.New(err.Error()))
		}
		list = append(list, m)
	}
	for i := range list{
		dao.cascade(list[i])
	}
	return list
}
func (dao *Dao) scanOne(sql string, args []interface{}) *$bean$ {
	rows := dao.Query(sql, args...)
	defer rows.Close()
	for rows.Next() {
		m := $bean${}
		err := rows.StructScan(&m)
		if err != nil {
			panic(exception.New(err.Error()))
		}
		dao.cascade(&m)
		return &m
	}
	return nil
}
////
```

## dao_no_cascade
```
/// auto template
type Dao struct {
	sqlkit.Dao
}
func New(schema string, tx ...*sqlkit.Dao) *Dao{
	dao:=&Dao{}
	dao.NewHelper(schema,tx...)
	return dao
}
func (dao *Dao) scan(sql string, args []interface{}) []*$bean$ {
	rows := dao.Query(sql, args...)
	list := make([]*$bean$,0,5)
	defer rows.Close()
	for rows.Next() {
		m := &$bean${}
		err := rows.StructScan(m)
		if err != nil {
			panic(exception.New(err.Error()))
		}
		list = append(list, m)
	}
	return list
}
func (dao *Dao) scanOne(sql string, args []interface{}) *$bean$ {
	rows := dao.Query(sql, args...)
	defer rows.Close()
	for rows.Next() {
		m := $bean${}
		err := rows.StructScan(&m)
		if err != nil {
			panic(exception.New(err.Error()))
		}
		return &m
	}
	return nil
}
////
```

## dao_demo
```
func (dao *Dao) FindById(id int32) *$bean$ {
	sql, args := sqlkit.Builder().Select("*").From(dao.GetTableD("$name$")).Where("id=?",id).MustSql()
	return dao.scanOne(sql, args)
}

type ListParam struct {
	IdList []int32
}

func (dao *Dao) List(param ListParam) []*$bean$ {
	builder := sqlkit.Builder().Select("*").From(dao.GetTableD("$name$")).Where("off=false").OrderBy("id")
	sql, args := builder.MustSql()
	return dao.scan(sql, args)
}
```

## exception
```
panic(exception.New("$msg$"))
```

## recover
```
defer func() {
    if err := recover(); err != nil {
        var msg string
        if e, ok := err.(exception.Exception); ok {
            msg = e.Msg
            // 带代码位置信息
            logkit.Error(e.Error())
        } else {
            msg = cast.ToString(err)
            logkit.Error(msg)
        }
    }
}()
```