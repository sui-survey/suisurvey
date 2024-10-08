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
        "/api/v1/create-form": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "CreateForm"
                ],
                "summary": "CreateForm API",
                "parameters": [
                    {
                        "description": "data",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.CreateFormDto"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/result.Result"
                        }
                    }
                }
            }
        },
        "/api/v1/form/{blobId}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "GetForm"
                ],
                "summary": "GetForm API",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Blob ID",
                        "name": "blobId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/result.Result"
                        }
                    },
                    "400": {
                        "description": "Invalid blob ID",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "enums.ItermType": {
            "type": "integer",
            "enum": [
                1,
                2,
                3,
                4,
                5
            ],
            "x-enum-varnames": [
                "TEXT",
                "TEXTAREA",
                "CHECKBOX",
                "RADIO",
                "SELECT"
            ]
        },
        "model.CreateFormDto": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "itemList": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Item"
                    }
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "model.Item": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "type": {
                    "$ref": "#/definitions/enums.ItermType"
                },
                "value": {
                    "type": "string"
                }
            }
        },
        "result.Result": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {},
                "message": {
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
	Title:            "Survey",
	Description:      "backend survey for walrus",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
