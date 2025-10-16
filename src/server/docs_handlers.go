package server

import (
	"encoding/json"
	"html/template"
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
                <a href="/openapi" class="active">API Docs</a>
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
                url: "/api/v1/openapi.json",
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
	t.Execute(w, nil)
}

// handleOpenAPISpec serves the OpenAPI specification JSON
func (s *Server) handleOpenAPISpec(w http.ResponseWriter, r *http.Request) {
	spec := map[string]interface{}{
		"openapi": "3.0.0",
		"info": map[string]interface{}{
			"title":       "Zipcode Lookup API",
			"description": "US Postal Code lookup and search API with 340,000+ zipcodes",
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
			{"name": "admin", "description": "Admin endpoints (authentication required)"},
		},
		"paths": map[string]interface{}{
			"/zipcodes.json": map[string]interface{}{
				"get": map[string]interface{}{
					"tags":        []string{"zipcodes"},
					"summary":     "Download complete dataset",
					"description": "Get the complete zipcodes dataset as JSON (340K+ records, 6.3MB)",
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Successful response",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"type": "array",
										"items": map[string]string{
											"$ref": "#/components/schemas/Zipcode",
										},
									},
								},
							},
						},
					},
				},
			},
			"/zipcode/search": map[string]interface{}{
				"get": map[string]interface{}{
					"tags":        []string{"zipcodes"},
					"summary":     "Search zipcodes",
					"description": "Search zipcodes by code, city, state, or prefix",
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
			"/zipcode/{code}": map[string]interface{}{
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
			"/zipcode/{code}.txt": map[string]interface{}{
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
			"/zipcode/city/{city}": map[string]interface{}{
				"get": map[string]interface{}{
					"tags":        []string{"zipcodes"},
					"summary":     "Get zipcodes by city",
					"description": "Get all zipcodes for a specific city",
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
			"/zipcode/state/{state}": map[string]interface{}{
				"get": map[string]interface{}{
					"tags":        []string{"zipcodes"},
					"summary":     "Get zipcodes by state",
					"description": "Get all zipcodes for a specific state",
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
			"/zipcode/autocomplete": map[string]interface{}{
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
											"success": map[string]string{"type": "boolean"},
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
			"/zipcode/stats": map[string]interface{}{
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
											"success": map[string]string{"type": "boolean"},
											"data": map[string]interface{}{
												"type": "object",
												"properties": map[string]interface{}{
													"total_zipcodes": map[string]string{"type": "integer"},
													"total_cities":   map[string]string{"type": "integer"},
													"total_states":   map[string]string{"type": "integer"},
													"total_counties": map[string]string{"type": "integer"},
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
			"/health": map[string]interface{}{
				"get": map[string]interface{}{
					"tags":        []string{"zipcodes"},
					"summary":     "Health check",
					"description": "Check API health status",
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Service is healthy",
						},
					},
				},
			},
		},
		"components": map[string]interface{}{
			"securitySchemes": map[string]interface{}{
				"bearerAuth": map[string]string{
					"type":   "http",
					"scheme": "bearer",
				},
			},
			"schemas": map[string]interface{}{
				"Zipcode": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"zipcode":   map[string]string{"type": "integer", "description": "5-digit zipcode"},
						"city":      map[string]string{"type": "string", "description": "City name"},
						"state":     map[string]string{"type": "string", "description": "State abbreviation"},
						"county":    map[string]string{"type": "string", "description": "County name"},
						"latitude":  map[string]string{"type": "string", "description": "Latitude coordinate"},
						"longitude": map[string]string{"type": "string", "description": "Longitude coordinate"},
					},
				},
				"ZipcodeResponse": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"success": map[string]string{"type": "boolean"},
						"data": map[string]interface{}{
							"$ref": "#/components/schemas/Zipcode",
						},
						"timestamp": map[string]string{"type": "string", "format": "date-time"},
					},
				},
				"SearchResponse": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"success": map[string]string{"type": "boolean"},
						"count":   map[string]string{"type": "integer"},
						"data": map[string]interface{}{
							"type": "array",
							"items": map[string]string{
								"$ref": "#/components/schemas/Zipcode",
							},
						},
						"timestamp": map[string]string{"type": "string", "format": "date-time"},
					},
				},
				"ErrorResponse": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"success": map[string]string{"type": "boolean"},
						"error": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"code":    map[string]string{"type": "string"},
								"message": map[string]string{"type": "string"},
								"field":   map[string]string{"type": "string"},
							},
						},
						"timestamp": map[string]string{"type": "string", "format": "date-time"},
					},
				},
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(spec)
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
                <a href="/openapi">API Docs</a>
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
                        query: '# Welcome to Zipcode Lookup GraphQL API\n# Press the Play button to run a query\n\nquery GetZipcode {\n  zipcode(code: "94102") {\n    zipcode\n    city\n    state\n    county\n    coordinates {\n      latitude\n      longitude\n    }\n  }\n}\n\nquery SearchByCity {\n  search(city: "San Francisco") {\n    zipcode\n    city\n    state\n  }\n}\n\nquery SearchByState {\n  search(state: "CA") {\n    zipcode\n    city\n    state\n    county\n  }\n}\n\nquery Stats {\n  stats {\n    totalZipcodes\n    totalCities\n    totalStates\n    totalCounties\n  }\n}'
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
	t.Execute(w, nil)
}

// handleGraphQL handles GraphQL queries
func (s *Server) handleGraphQL(w http.ResponseWriter, r *http.Request) {
	// For now, return a simple message
	// Full GraphQL implementation would go here
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":    "GraphQL endpoint - Full implementation coming soon",
		"playground": "/graphql",
	})
}
