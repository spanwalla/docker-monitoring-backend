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
        "/api/v1/pingers/login": {
            "post": {
                "description": "Login",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "pinger"
                ],
                "summary": "Login",
                "parameters": [
                    {
                        "description": "input",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_controller_http_v1.loginInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/internal_controller_http_v1.pingerRoutes"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    }
                }
            }
        },
        "/api/v1/pingers/register": {
            "post": {
                "description": "Register pinger",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "pinger"
                ],
                "summary": "Register",
                "parameters": [
                    {
                        "description": "input",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_controller_http_v1.registerInput"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/internal_controller_http_v1.pingerRoutes"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    }
                }
            }
        },
        "/api/v1/reports": {
            "get": {
                "description": "Get latest report by every pinger ever exists in database",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "reports"
                ],
                "summary": "Get reports",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/github_com_spanwalla_docker-monitoring-backend_internal_service.ReportOutput"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "JWT": []
                    }
                ],
                "description": "Store pinger's report to database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "reports"
                ],
                "summary": "Store report",
                "parameters": [
                    {
                        "description": "input",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_controller_http_v1.storeReportInput"
                        }
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "echo.HTTPError": {
            "type": "object",
            "properties": {
                "message": {}
            }
        },
        "github_com_spanwalla_docker-monitoring-backend_internal_service.PingResult": {
            "type": "object",
            "required": [
                "id",
                "ip",
                "latency_ms",
                "state",
                "status",
                "timestamp"
            ],
            "properties": {
                "id": {
                    "type": "string"
                },
                "ip": {
                    "type": "string"
                },
                "latency_ms": {
                    "type": "integer"
                },
                "state": {
                    "type": "string",
                    "enum": [
                        "created",
                        "restarting",
                        "running",
                        "removing",
                        "paused",
                        "exited",
                        "dead"
                    ]
                },
                "status": {
                    "type": "string"
                },
                "timestamp": {
                    "type": "string"
                }
            }
        },
        "github_com_spanwalla_docker-monitoring-backend_internal_service.ReportOutput": {
            "type": "object",
            "properties": {
                "content": {
                    "$ref": "#/definitions/pgtype.JSONB"
                },
                "created_at": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "pinger_name": {
                    "type": "string"
                }
            }
        },
        "internal_controller_http_v1.loginInput": {
            "type": "object",
            "required": [
                "name",
                "password"
            ],
            "properties": {
                "name": {
                    "type": "string",
                    "maxLength": 32,
                    "minLength": 4
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "internal_controller_http_v1.pingerRoutes": {
            "type": "object"
        },
        "internal_controller_http_v1.registerInput": {
            "type": "object",
            "required": [
                "name",
                "password"
            ],
            "properties": {
                "name": {
                    "type": "string",
                    "maxLength": 32,
                    "minLength": 4
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "internal_controller_http_v1.storeReportInput": {
            "type": "object",
            "required": [
                "report"
            ],
            "properties": {
                "report": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/github_com_spanwalla_docker-monitoring-backend_internal_service.PingResult"
                    }
                }
            }
        },
        "pgtype.JSONB": {
            "type": "object",
            "properties": {
                "bytes": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "status": {
                    "$ref": "#/definitions/pgtype.Status"
                }
            }
        },
        "pgtype.Status": {
            "type": "integer",
            "enum": [
                0,
                1,
                2
            ],
            "x-enum-varnames": [
                "Undefined",
                "Null",
                "Present"
            ]
        }
    },
    "securityDefinitions": {
        "JWT": {
            "description": "JSON Web Token",
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
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Docker Monitoring Service",
	Description:      "This is a service for storing and showing docker container's reports.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
