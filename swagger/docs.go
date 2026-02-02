// Package swagger - minimal stub so import works. Run: swag init -g cmd/server/main.go -o swagger
package swagger

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": ["http", "https"],
    "swagger": "2.0",
    "info": {
        "title": "Go E-Commerce API",
        "version": "1.0",
        "description": "A mini e-commerce API built with Go, PostgreSQL, and sqlc"
    },
    "host": "localhost:8080",
    "basePath": "/",
    "paths": {},
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header",
            "description": "Type Bearer followed by space and JWT token"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info (run swag init to regenerate with full docs)
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{"http", "https"},
	Title:            "Go E-Commerce API",
	Description:      "A mini e-commerce API built with Go, PostgreSQL, and sqlc",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
