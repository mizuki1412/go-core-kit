package tag

var DBField = Tag{Name: "db"}

var DBPk = Tag{Name: "pk", Value: "true"}

var DBTable = Tag{Name: "table"}

var DBPkAuto = Tag{Name: "auto", Value: "true"}

// DBColumnLogicDel 逻辑删除字段
var DBColumnLogicDel = Tag{Name: "logicDel", Value: "true"}

// DecimalPrecision decimal精度
var DecimalPrecision = Tag{Name: "precision"}
