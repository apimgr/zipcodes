package server

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

// handleSwaggerUI serves the Swagger UI for API documentation with site theme
func (s *Server) handleSwaggerUI(w http.ResponseWriter, r *http.Request) {
	tmpl := `<!DOCTYPE html>
<html lang="en" data-theme="dark">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>API Documentation - Zipcode Lookup</title>
    <link rel="stylesheet" href="/static/css/main.css">
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.10.0/swagger-ui.css">
    <style>
        body { margin: 0; padding: 0; display: flex; flex-direction: column; min-height: 100vh; }
        #swagger-container { flex: 1; }
        .swagger-ui { background: var(--bg-primary); }
        .swagger-ui .topbar { display: none; }
        .swagger-ui .info { color: var(--text-primary); }
        .swagger-ui .scheme-container { background: var(--bg-secondary); }
        .swagger-ui .opblock { background: var(--bg-secondary); border-color: var(--border-color); }
        .swagger-ui .opblock-tag { color: var(--text-primary); border-color: var(--border-color); }
        .swagger-ui .opblock-summary { background: var(--bg-tertiary); }
        .swagger-ui .opblock-description { color: var(--text-secondary); }
        .swagger-ui table thead tr td, .swagger-ui table thead tr th { color: var(--text-primary); border-color: var(--border-color); }
        .swagger-ui .parameter__name { color: var(--accent-primary); }
        .swagger-ui .response-col_status { color: var(--accent-success); }
        .swagger-ui input, .swagger-ui select, .swagger-ui textarea { background: var(--bg-tertiary); color: var(--text-primary); border-color: var(--border-color); }
        .swagger-ui .btn { background: var(--accent-primary); color: white; }
    </style>
</head>
<body data-theme="dark">
    <header id="main-header">
        <div class="header-container">
            <div class="header-left">
                <a class="logo" href="/">📮 Zipcode Lookup</a>
            </div>
            <nav id="main-nav" class="header-center">
                <a href="/">Search</a>
                <a href="/server/swagger" class="active">API Docs</a>
                <a href="/graphql">GraphQL</a>
            </nav>
            <div class="header-right">
                <button id="theme-toggle" class="btn-icon" aria-label="Toggle theme">🌙</button>
            </div>
        </div>
    </header>

    <div id="swagger-container">
        <div id="swagger-ui"></div>
    </div>

    <script src="/static/js/main.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.0/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.0/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                url: "/api/v1/server/openapi.json",
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout"
            });
            window.ui = ui;
        };
    </script>
</body>
</html>`

	t, err := template.New("swagger").Parse(tmpl)
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.Execute(w, nil); err != nil {
		log.Printf("render documentation template: %v", err)
	}
}

// handleOpenAPISpec serves the OpenAPI specification JSON
func (s *Server) handleOpenAPISpec(w http.ResponseWriter, r *http.Request) {
	spec := map[string]interface{}{
		"openapi": "3.0.0",
		"info": map[string]interface{}{
			"title":       "Zipcode Lookup API",
			"description": "Public US ZIP lookup API using an embedded 33,100-record geographic dataset",
			"version":     "1.0.0",
			"contact": map[string]string{
				"name": "Zipcode Lookup API",
				"url":  "https://github.com/apimgr/zipcodes",
			},
			"license": map[string]string{
				"name": "MIT",
				"url":  "https://opensource.org/licenses/MIT",
			},
		},
		"servers": []map[string]string{
			{"url": "/api/v1", "description": "API v1"},
		},
		"tags": []map[string]string{
			{"name": "zipcodes", "description": "Zipcode data endpoints"},
			{"name": "geoip", "description": "GeoIP location endpoints"},
		},
		"paths": map[string]interface{}{
			"/zipcodes.json": map[string]interface{}{
				"get": map[string]interface{}{
					"tags":        []string{"zipcodes"},
					"summary":     "Download complete dataset",
					"description": "Get the embedded ZIP dataset as JSON (33,100 source records)",
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Successful response",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"type": "object",
										"properties": map[string]interface{}{
											"ok":   map[string]string{"type": "boolean"},
											"data": map[string]interface{}{"type": "array", "items": map[string]string{"$ref": "#/components/schemas/Zipcode"}},
										},
									},
								},
							},
						},
					},
				},
			},
			"/me": map[string]interface{}{
				"get": map[string]interface{}{
					"tags":        []string{"server"},
					"summary":     "Look up the caller IP",
					"description": "Returns the request IP and GeoIP fields when a public address and database are available.",
					"responses":   map[string]interface{}{"200": map[string]interface{}{"description": "Caller location or IP-only result"}},
				},
			},
			"/zipcodes/nearby": map[string]interface{}{
				"get": map[string]interface{}{
					"tags":        []string{"zipcodes"},
					"summary":     "Search ZIP codes within a radius",
					"description": "Supply zip or a lat/lng pair and a radius in miles. The radius is bounded by server configuration and never exceeds 100 miles.",
					"parameters": []map[string]interface{}{
						{"name": "zip", "in": "query", "schema": map[string]string{"type": "string", "pattern": "^[0-9]{5}$"}},
						{"name": "lat", "in": "query", "schema": map[string]interface{}{"type": "number", "minimum": -90, "maximum": 90}},
						{"name": "lng", "in": "query", "schema": map[string]interface{}{"type": "number", "minimum": -180, "maximum": 180}},
						{"name": "radius", "in": "query", "required": true, "schema": map[string]interface{}{"type": "number", "exclusiveMinimum": 0, "maximum": 100}},
					},
					"responses": map[string]interface{}{"200": map[string]interface{}{"description": "Nearest ZIP records within the requested radius"}, "400": map[string]interface{}{"description": "Invalid radius or origin"}},
				},
			},
			"/zipcodes/search": map[string]interface{}{
				"get": map[string]interface{}{
					"tags":        []string{"zipcodes"},
					"summary":     "Search zipcodes",
					"description": "Search zipcodes by code, city, state, or prefix. List results are paginated with page and limit.",
					"parameters": []map[string]interface{}{
						{
							"name":        "q",
							"in":          "query",
							"description": "Search query (zipcode, city, state, or prefix)",
							"required":    true,
							"schema":      map[string]string{"type": "string"},
							"examples": map[string]interface{}{
								"zipcode": map[string]string{
									"value":   "94102",
									"summary": "Search by zipcode",
								},
								"city": map[string]string{
									"value":   "San Francisco",
									"summary": "Search by city",
								},
								"state": map[string]string{
									"value":   "CA",
									"summary": "Search by state",
								},
								"cityState": map[string]string{
									"value":   "New York, NY",
									"summary": "Search by city and state",
								},
							},
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Successful response",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": "#/components/schemas/SearchResponse",
									},
								},
							},
						},
						"400": map[string]interface{}{
							"description": "Bad request",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": "#/components/schemas/ErrorResponse",
									},
								},
							},
						},
					},
				},
			},
			"/zipcodes/{code}": map[string]interface{}{
				"get": map[string]interface{}{
					"tags":        []string{"zipcodes"},
					"summary":     "Get zipcode details",
					"description": "Get detailed information for a specific zipcode",
					"parameters": []map[string]interface{}{
						{
							"name":        "code",
							"in":          "path",
							"description": "5-digit zipcode",
							"required":    true,
							"schema":      map[string]string{"type": "string"},
							"example":     "94102",
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Successful response",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": "#/components/schemas/ZipcodeResponse",
									},
								},
							},
						},
						"404": map[string]interface{}{
							"description": "Zipcode not found",
						},
					},
				},
			},
			"/zipcodes/{code}.txt": map[string]interface{}{
				"get": map[string]interface{}{
					"tags":        []string{"zipcodes"},
					"summary":     "Get zipcode as plain text",
					"description": "Get zipcode information in plain text format",
					"parameters": []map[string]interface{}{
						{
							"name":        "code",
							"in":          "path",
							"description": "5-digit zipcode",
							"required":    true,
							"schema":      map[string]string{"type": "string"},
							"example":     "94102",
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Successful response",
							"content": map[string]interface{}{
								"text/plain": map[string]interface{}{
									"schema": map[string]string{"type": "string"},
								},
							},
						},
					},
				},
			},
			"/zipcodes/city/{city}": map[string]interface{}{
				"get": map[string]interface{}{
					"tags":        []string{"zipcodes"},
					"summary":     "Get zipcodes by city",
					"description": "Get paginated ZIP records for a specific city (default 250 records per page).",
					"parameters": []map[string]interface{}{
						{
							"name":        "city",
							"in":          "path",
							"description": "City name",
							"required":    true,
							"schema":      map[string]string{"type": "string"},
							"example":     "San Francisco",
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Successful response",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": "#/components/schemas/SearchResponse",
									},
								},
							},
						},
					},
				},
			},
			"/zipcodes/state/{state}": map[string]interface{}{
				"get": map[string]interface{}{
					"tags":        []string{"zipcodes"},
					"summary":     "Get zipcodes by state",
					"description": "Get paginated ZIP records for a specific state (default 250 records per page).",
					"parameters": []map[string]interface{}{
						{
							"name":        "state",
							"in":          "path",
							"description": "State code (2 letters)",
							"required":    true,
							"schema":      map[string]string{"type": "string"},
							"example":     "CA",
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Successful response",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": "#/components/schemas/SearchResponse",
									},
								},
							},
						},
					},
				},
			},
			"/zipcodes/autocomplete": map[string]interface{}{
				"get": map[string]interface{}{
					"tags":        []string{"zipcodes"},
					"summary":     "Autocomplete suggestions",
					"description": "Get autocomplete suggestions for zipcode search",
					"parameters": []map[string]interface{}{
						{
							"name":        "q",
							"in":          "query",
							"description": "Search query",
							"required":    true,
							"schema":      map[string]string{"type": "string"},
						},
						{
							"name":        "limit",
							"in":          "query",
							"description": "Maximum number of suggestions (1-50, default: 10)",
							"schema":      map[string]string{"type": "integer"},
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Successful response",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"type": "object",
										"properties": map[string]interface{}{
											"ok": map[string]string{"type": "boolean"},
											"suggestions": map[string]interface{}{
												"type": "array",
												"items": map[string]string{
													"type": "string",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"/zipcodes/stats": map[string]interface{}{
				"get": map[string]interface{}{
					"tags":        []string{"zipcodes"},
					"summary":     "Get database statistics",
					"description": "Get statistics about the zipcode database",
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Successful response",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"type": "object",
										"properties": map[string]interface{}{
											"ok": map[string]string{"type": "boolean"},
											"data": map[string]interface{}{
												"type": "object",
												"properties": map[string]interface{}{
													"total_zipcodes": map[string]string{"type": "integer"},
													"total_cities":   map[string]string{"type": "integer"},
													"total_states":   map[string]string{"type": "integer"},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"/geoip": map[string]interface{}{
				"get": map[string]interface{}{
					"tags":        []string{"geoip"},
					"summary":     "Lookup request IP",
					"description": "Get geolocation information for the request IP address",
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Successful response",
						},
					},
				},
			},
			"/geoip.txt": map[string]interface{}{
				"get": map[string]interface{}{
					"tags":        []string{"geoip"},
					"summary":     "Lookup request IP (text)",
					"description": "Get geolocation information for the request IP as plain text",
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Successful response",
							"content": map[string]interface{}{
								"text/plain": map[string]interface{}{
									"schema": map[string]string{"type": "string"},
								},
							},
						},
					},
				},
			},
		},
		"components": map[string]interface{}{
			"schemas": map[string]interface{}{
				"Zipcode": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"zip":        map[string]string{"type": "string", "pattern": "^[0-9]{5}$"},
						"zip_code":   map[string]string{"type": "integer", "description": "Legacy numeric ZIP field"},
						"state_full": map[string]string{"type": "string"},
						"timezone":   map[string]string{"type": "string"},
						"area_codes": map[string]interface{}{"type": "array", "items": map[string]string{"type": "string"}},
						"city":       map[string]string{"type": "string", "description": "City name"},
						"state":      map[string]string{"type": "string", "description": "State abbreviation"},
						"county":     map[string]string{"type": "string", "description": "County name"},
						"latitude":   map[string]interface{}{"type": "number", "format": "float", "nullable": true},
						"longitude":  map[string]interface{}{"type": "number", "format": "float", "nullable": true},
					},
				},
				"ZipcodeResponse": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"ok": map[string]string{"type": "boolean"},
						"data": map[string]interface{}{
							"$ref": "#/components/schemas/Zipcode",
						},
					},
				},
				"SearchResponse": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"ok": map[string]string{"type": "boolean"},
						"data": map[string]interface{}{
							"type":  "array",
							"items": map[string]string{"$ref": "#/components/schemas/Zipcode"},
						},
						"pagination": map[string]interface{}{"$ref": "#/components/schemas/Pagination"},
					},
				},
				"Pagination": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"page":  map[string]string{"type": "integer"},
						"limit": map[string]string{"type": "integer"},
						"total": map[string]string{"type": "integer"},
						"pages": map[string]string{"type": "integer"},
					},
				},
				"ErrorResponse": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"ok":      map[string]string{"type": "boolean"},
						"error":   map[string]string{"type": "string"},
						"message": map[string]string{"type": "string"},
						"details": map[string]interface{}{"type": "object", "additionalProperties": true},
					},
				},
			},
		},
	}

	paths := spec["paths"].(map[string]interface{})
	pageParameters := []map[string]interface{}{
		{"name": "page", "in": "query", "description": "Page number (default 1)", "schema": map[string]interface{}{"type": "integer", "minimum": 1, "default": 1}},
		{"name": "limit", "in": "query", "description": "Records per page (default 250, maximum 250)", "schema": map[string]interface{}{"type": "integer", "minimum": 1, "maximum": 250, "default": 250}},
	}
	for _, path := range []string{"/zipcodes/search", "/zipcodes/nearby", "/zipcodes/city/{city}", "/zipcodes/state/{state}"} {
		pathItem, ok := paths[path].(map[string]interface{})
		if !ok {
			continue
		}
		operation, ok := pathItem["get"].(map[string]interface{})
		if !ok {
			continue
		}
		parameters, _ := operation["parameters"].([]map[string]interface{})
		operation["parameters"] = append(parameters, pageParameters...)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(spec)
}

// handleGraphQLPlayground serves the GraphQL Playground with site theme
func (s *Server) handleGraphQLPlayground(w http.ResponseWriter, r *http.Request) {
	tmpl := `<!DOCTYPE html>
<html lang="en" data-theme="dark">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>GraphQL Playground - Zipcode Lookup</title>
    <link rel="stylesheet" href="/static/css/main.css">
    <style>
        body { margin: 0; padding: 0; display: flex; flex-direction: column; min-height: 100vh; }
        #graphql-container { flex: 1; display: flex; flex-direction: column; }
        #root { flex: 1; }
    </style>
</head>
<body data-theme="dark">
    <header id="main-header">
        <div class="header-container">
            <div class="header-left">
                <a class="logo" href="/">📮 Zipcode Lookup</a>
            </div>
            <nav id="main-nav" class="header-center">
                <a href="/">Search</a>
                <a href="/server/swagger">API Docs</a>
                <a href="/graphql" class="active">GraphQL</a>
            </nav>
            <div class="header-right">
                <button id="theme-toggle" class="btn-icon" aria-label="Toggle theme">🌙</button>
            </div>
        </div>
    </header>

    <div id="graphql-container">
        <div id="root"></div>
    </div>

    <link rel="stylesheet" href="https://unpkg.com/graphql-playground-react@1.7.28/build/static/css/index.css">
    <script src="/static/js/main.js"></script>
    <script src="https://unpkg.com/graphql-playground-react@1.7.28/build/static/js/middleware.js"></script>
    <script>
        window.addEventListener('load', function (event) {
            GraphQLPlayground.init(document.getElementById('root'), {
                endpoint: '/api/v1/graphql',
                settings: {
                    'editor.theme': 'dark',
                    'editor.cursorShape': 'line',
                    'theme': 'dark'
                },
                tabs: [
                    {
                        endpoint: '/api/v1/graphql',
                        query: '# Welcome to Zipcode Lookup GraphQL API\n# Press the Play button to run a query\n\nquery GetZipcode {\n  zipcode(code: "94102") {\n    zipcode\n    city\n    state\n    county\n    coordinates {\n      latitude\n      longitude\n    }\n  }\n}\n\nquery SearchByCity {\n  search(city: "San Francisco") {\n    zipcode\n    city\n    state\n  }\n}\n\nquery SearchByState {\n  search(state: "CA") {\n    zipcode\n    city\n    state\n    county\n  }\n}\n\nquery Stats {\n  stats {\n    totalZipcodes\n    totalCities\n    totalStates\n  }\n}'
                    }
                ]
            });
        });
    </script>
</body>
</html>`

	t, err := template.New("graphql").Parse(tmpl)
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.Execute(w, nil); err != nil {
		log.Printf("render documentation template: %v", err)
	}
}

// handleGraphQL handles GraphQL queries
func (s *Server) handleGraphQL(w http.ResponseWriter, r *http.Request) {
	s.serveGraphQL(w, r)
}
