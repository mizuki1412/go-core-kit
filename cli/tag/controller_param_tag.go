package tag

// Validate 是否必填
var Validate = Tag{Name: "validate"}

const ValidateRequired = "required"

// Comment 描述
var Comment = Tag{Name: "comment"}

// Default 默认值
var Default = Tag{Name: "default"}

// Trim 去空格：request时处理
var Trim = Tag{Name: "trim", Value: "true"}
