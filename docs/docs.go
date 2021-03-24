// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "index"
                ],
                "summary": "Index",
                "responses": {
                    "200": {
                        "description": "Here is the neuron server",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/admin/download": {
            "post": {
                "description": "let the master download latest data",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Download",
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "$ref": "#/definitions/utils.WebResponse"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/utils.WebResponse"
                        }
                    },
                    "500": {
                        "description": "Internal error",
                        "schema": {
                            "$ref": "#/definitions/utils.WebResponse"
                        }
                    }
                }
            }
        },
        "/admin/report": {
            "post": {
                "description": "let the master report special date",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Report",
                "parameters": [
                    {
                        "type": "string",
                        "description": "The report date",
                        "name": "date",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "boolean",
                        "description": "Full report or not",
                        "name": "full",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "$ref": "#/definitions/utils.WebResponse"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/utils.WebResponse"
                        }
                    },
                    "500": {
                        "description": "Internal error",
                        "schema": {
                            "$ref": "#/definitions/utils.WebResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "utils.WebResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "error": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "",
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "The ARK ROBOT Server",
	Description: "The ARK ROBOT REST API",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
