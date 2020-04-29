package swagger

type SwaggerPath struct {
	Path string
}

func (swagger *SwaggerPath) Param(param interface{}) *SwaggerPath {
	// todo
	return swagger
}
func (swagger *SwaggerPath) Description(val string) *SwaggerPath {
	// todo
	return swagger
}
func (swagger *SwaggerPath) Summary(val string) *SwaggerPath {
	// todo
	return swagger
}

type Doc struct{}

func (s *Doc) ReadDoc() string {
	return `
 {
    "schemes": [],
    "swagger": "2.0",
    "info": {
        "description": "This is a sample server Petstore server.",
        "title": "My Swagger API2",
        "contact": {},
        "license": {},
        "version": "1.0"
    },
    "host": "127.0.0.1:10000",
    "basePath": "/server",
    "paths": {
        "/rest/user/loginByUsername": {
            "post": {
                "description": "登录223",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Show a account",
                "operationId": "user.loginByUsername",
                "parameters": [
                    {
                        "type": "string",
                        "description": "username",
                        "name": "username",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "pwd",
                        "name": "pwd",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "schema",
                        "name": "schema",
                        "in": "path"
                    }
                ]
            }
        }
    }
}`
}

/**
note: https://github.com/swaggo/swag
直接通过swag库的代码方式实现

标准：https://swagger.io/specification/v2/
*/
