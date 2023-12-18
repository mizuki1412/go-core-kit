package openapi

type ApiDocV3 struct {
	Openapi    string                                       `json:"openapi,omitempty"`
	Info       *ApiDocV3Info                                `json:"info,omitempty"`
	Paths      map[string]map[string]*ApiDocV3PathOperation `json:"paths,omitempty"` // path:method:info, 这里只处理 method 对应的 operation
	Components *ApiDocV3ComponentObj                        `json:"components,omitempty"`
	Servers    []string                                     `json:"servers,omitempty"` // 目前只要填写 url
}

type ApiDocV3Info struct {
	Title       string               `json:"title,omitempty"`
	Description string               `json:"description,omitempty"`
	Contact     *ApiDocV3InfoContact `json:"contact,omitempty"`
	License     *ApiDocV3InfoLicense `json:"license,omitempty"`
	Version     string               `json:"version,omitempty"`
}

type ApiDocV3InfoContact struct {
	Name  string `json:"name,omitempty"`
	Url   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

type ApiDocV3InfoLicense struct {
	Name string `json:"name,omitempty"`
	Url  string `json:"url,omitempty"`
}

type ApiDocV3PathOperation struct {
	Tags        []string                    `json:"tags,omitempty"` // 一般就一个
	Summary     string                      `json:"summary,omitempty"`
	Description string                      `json:"description,omitempty"`
	OperationId string                      `json:"operationId,omitempty"`
	Deprecated  bool                        `json:"deprecated,omitempty"`
	Parameters  []*ApiDocV3ReqParam         `json:"parameters,omitempty"`
	RequestBody *ApiDocV3ReqBody            `json:"requestBody,omitempty"`
	Responses   map[string]*ApiDocV3ResBody `json:"responses,omitempty"` // key 为 default 或 具体的http-code
}

type ApiDocV3ReqParam struct {
	Name        string          `json:"name,omitempty"` // If in is "header" and the name field is "Accept", "Content-Type" or "Authorization"
	In          string          `json:"in,omitempty"`   // "query", "header", "path" or "cookie".
	Description string          `json:"description,omitempty"`
	Required    bool            `json:"required,omitempty"`
	Deprecated  bool            `json:"deprecated,omitempty"`
	Schema      *ApiDocV3Schema `json:"schema,omitempty"`
	// allowEmptyValue
	// style
}

type ApiDocV3ReqBody struct {
	Description string                            `json:"description,omitempty"`
	Required    bool                              `json:"required,omitempty"`
	Content     map[string]*ApiDocV3SchemaWrapper `json:"content,omitempty"` // key 为 media type：application/json，multipart/form-data
}

type ApiDocV3ResBody struct {
	Description string                            `json:"description,omitempty"`
	Content     map[string]*ApiDocV3SchemaWrapper `json:"content,omitempty"` // key 为 media type：application/json，
	// header
}

type ApiDocV3SchemaWrapper struct {
	Schema *ApiDocV3Schema `json:"schema,omitempty"`
}

type ApiDocV3Schema struct {
	Ref              string                     `json:"$ref,omitempty"`   // ref 和其他的字段不共存 eg:#/components/schemas/LoginReq
	Type             string                     `json:"type,omitempty"`   // integer, string, number, boolean, object, array
	Format           string                     `json:"format,omitempty"` // int32/int64, float/double, byte/binary/date/date-time/password/或者任意需要客户端解析-email/uuid/...,
	Default          any                        `json:"default,omitempty"`
	Title            string                     `json:"title,omitempty"`
	Properties       map[string]*ApiDocV3Schema `json:"properties,omitempty"` // type=object, key=属性名
	Items            *ApiDocV3Schema            `json:"items,omitempty"`      // type=array 时
	Description      string                     `json:"description,omitempty"`
	Required         []string                   `json:"required,omitempty"` // properties中的key
	ContentMediaType string                     `json:"contentMediaType,omitempty"`
}

type ApiDocV3ComponentObj struct {
	Schemas map[string]*ApiDocV3Schema `json:"schemas,omitempty"` // key = 名称
	// SecuritySchemes
}
