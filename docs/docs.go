// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/apps": {
            "post": {
                "tags": [
                    "Apps"
                ],
                "summary": "Creates an IskayPet Application",
                "parameters": [
                    {
                        "description": "Body params",
                        "name": "createAppModel",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.CreateAppModel"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/apps/groups": {
            "get": {
                "description": "Needed for create a project in a specific group",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Groups"
                ],
                "summary": "Get all groups from GitLab",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.AppGroupModel"
                            }
                        }
                    }
                }
            }
        },
        "/apps/search": {
            "get": {
                "tags": [
                    "Apps"
                ],
                "summary": "Get relevant info for an IskayPet app",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Application name",
                        "name": "app_name",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/apps/types": {
            "get": {
                "description": "For example: app_name: users-api, key: token, value: hash",
                "tags": [
                    "Apps"
                ],
                "summary": "Get all application types (backend, frontend, etc.)",
                "responses": {}
            }
        },
        "/apps/{appId}/secrets": {
            "post": {
                "tags": [
                    "Secrets"
                ],
                "summary": "Creates secret for application",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "App ID",
                        "name": "appId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Body params",
                        "name": "createAppSecretModel",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.CreateAppSecretModel"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/ping": {
            "get": {
                "description": "Health",
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Check"
                ],
                "summary": "Check if the instance is healthy or unhealthy",
                "responses": {
                    "200": {
                        "description": "pong",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.AppGroupModel": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "model.CreateAppModel": {
            "type": "object",
            "properties": {
                "app_type_id": {
                    "type": "integer"
                },
                "group_id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "model.CreateAppSecretModel": {
            "type": "object",
            "properties": {
                "key": {
                    "type": "string"
                },
                "value": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Pets API",
	Description:      "Create apps, services and infrastructure.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
