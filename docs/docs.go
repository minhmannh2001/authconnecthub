// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Nguyen Minh Manh",
            "email": "nguyenminhmannh2001@gmail.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/": {
            "get": {
                "description": "This endpoint renders the index.html page with potential toast notification settings based on query parameters and validation.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "text/html"
                ],
                "tags": [
                    "index"
                ],
                "summary": "Get Index Page",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Message to display in the toast notification",
                        "name": "toast-message",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Type of the toast notification (e.g., success, error)",
                        "name": "toast-type",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Hash value used for validation (optional)",
                        "name": "hash-value",
                        "in": "query"
                    }
                ],
                "responses": {}
            }
        },
        "/private": {
            "get": {
                "security": [
                    {
                        "JWT": []
                    }
                ],
                "description": "This endpoint is accessible only to authorized users and returns a greeting message.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "private"
                ],
                "summary": "Access a private resource",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/v1/auth/login": {
            "get": {
                "description": "This endpoint renders the login page and displays a toast notification if provided query parameters are valid.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "text/html"
                ],
                "tags": [
                    "Authen"
                ],
                "summary": "Login Page",
                "parameters": [
                    {
                        "type": "string",
                        "description": "The message to display in the toast notification.",
                        "name": "toast-message",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The type of the toast notification (e.g., success, error).",
                        "name": "toast-type",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "A hash value used for validation.",
                        "name": "hash-value",
                        "in": "query"
                    }
                ],
                "responses": {}
            }
        },
        "/v1/auth/logout": {
            "get": {
                "security": [
                    {
                        "JWT": []
                    }
                ],
                "description": "Logs out the currently authenticated user and redirects to the home page with a success toast notification.",
                "tags": [
                    "Authen"
                ],
                "summary": "Logout User",
                "responses": {}
            }
        }
    },
    "securityDefinitions": {
        "JWT": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/v1",
	Schemes:          []string{},
	Title:            "AuthConnect Hub",
	Description:      "A centralized authentication hub for my home applications in Go using Gin framework.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
