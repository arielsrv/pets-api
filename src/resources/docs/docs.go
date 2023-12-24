// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
        "/apps/{appID}/secrets": {
            "post": {
                "description": "Get snippet key, conflict if secret already exist.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "secrets"
                ],
                "summary": "Creates the secret",
                "parameters": [
                    {
                        "description": "Body params",
                        "name": "createAppSecretModel",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.CreateSecretRequest"
                        }
                    },
                    {
                        "type": "integer",
                        "description": "Pet ID",
                        "name": "appID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.CreateSecretResponse"
                            }
                        }
                    },
                    "404": {
                        "description": "App not found",
                        "schema": {
                            "$ref": "#/definitions/server.Error"
                        }
                    },
                    "409": {
                        "description": "Key already exist",
                        "schema": {
                            "$ref": "#/definitions/server.Error"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/server.Error"
                        }
                    }
                }
            }
        },
        "/ping": {
            "get": {
                "description": "Health",
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "health"
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
        "model.CreateSecretRequest": {
            "type": "object",
            "properties": {
                "key": {
                    "type": "string"
                },
                "value": {
                    "type": "string"
                }
            }
        },
        "model.CreateSecretResponse": {
            "type": "object",
            "properties": {
                "key": {
                    "type": "string"
                },
                "original_key": {
                    "type": "string"
                },
                "snippet_url": {
                    "type": "string"
                }
            }
        },
        "server.Error": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "status_code": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Pets API.",
	Description:      "Backend for Pets Clients.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
