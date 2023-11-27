# goland-live-templates

用于golang的模板代码示例

## cli_init

```
cli.RootCMD(&cobra.Command{
    Use: "main",
    Run: func(cmd *cobra.Command, args []string) {
        restkit.AddActions(user.All()...)
        restkit.AddActions(download.Init)
        _ = restkit.Run()
    },
})
cli.AddChildCMD(&cobra.Command{
    Use: "test",
    Run: func(cmd *cobra.Command, args []string) {
       
    },
})
cli.Execute()
```

## action_init
```
func Init(router *router.Router) {
	tag := "$tag$:$tname$"
	r := router.Group("/rest/$tag$")
	r.Use(middleware.AuthUsernameAndPwd())
	{
		r.Post("/$name$", $name$).Openapi.Tag(tag).Summary("$summary$").ReqParam($name$Params{})
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

## action_init_full

```
func Init(router *router.Router) {
	tag := "$tag$:$tname$"
	r := router.Group("/rest/$tag$")
	r.Use(middleware.AuthUsernameAndPwd())
	{
		r.Post("/update", update).Openapi.Tag(tag).Summary("增加和修改").ReqParam(updateParams{})
		r.Post("/del", del).Openapi.Tag(tag).Summary("删除").ReqParam(delParams{})
		r.Post("/list", list).Openapi.Tag(tag).Summary("列表").ReqParam(listParams{})
		r.Post("/detail", detail).Openapi.Tag(tag).Summary("详情").ReqParam(detailParams{})
	}
}

type updateParams struct {
	//Phone    string `description:"手机号" default:"" trim:"true"`
	//Pwd      string `validate:"required"`
}

func update(ctx *context.Context) {
	params := updateParams{}
	ctx.BindForm(&params)

	ctx.JsonSuccess(nil)
}

type delParams struct {
	Id int32 `validate:"required"`
}

func del(ctx *context.Context) {
	params := delParams{}
	ctx.BindForm(&params)

	ctx.JsonSuccess(nil)
}

type listParams struct{}

func list(ctx *context.Context) {
	params := listParams{}
	ctx.BindForm(&params)

	ctx.JsonSuccess(nil)
}

type detailParams struct {
	Id int32 `validate:"required"`
}

func detail(ctx *context.Context) {
	params := detailParams{}
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
func (th *$struct$) Scan(value any) error {
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
func (l $name$List) MapReduce(fun func(ele *$name$) any) []any {
	var results []any
	for _, e := range l {
		results = append(results, fun(e))
	}
	return results
}
```

## dao_new
```
type Dao struct {
	sqlkit.Dao[$name$]
}

const (
	ResultDefault byte = iota
	ResultNone
)

func New(cascadeType byte, ds ...*sqlkit.DataSource) Dao {
	dao := Dao{sqlkit.New[$name$](ds...)}
	dao.Cascade = func(obj *$name$) {
		switch cascadeType {
		case ResultDefault:
		case ResultNone:
		}
	}
	return dao
}
```

## dao_new_no_cascade
```
type Dao struct {
	sqlkit.Dao[$name$]
}

func New(ds ...*sqlkit.DataSource) Dao {
	return Dao{sqlkit.New[$name$](ds...)}
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