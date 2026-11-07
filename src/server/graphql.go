package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/apimgr/zipcodes/src/database"
	"github.com/graphql-go/graphql"
)

type graphQLRequest struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

func (s *Server) serveGraphQL(w http.ResponseWriter, r *http.Request) {
	request, err := decodeGraphQLRequest(w, r)
	if err != nil {
		writeGraphQLJSON(w, http.StatusBadRequest, map[string]interface{}{"errors": []map[string]string{{"message": err.Error()}}})
		return
	}
	schema, err := s.graphQLSchema()
	if err != nil {
		writeGraphQLJSON(w, http.StatusInternalServerError, map[string]interface{}{"errors": []map[string]string{{"message": "GraphQL schema is unavailable"}}})
		return
	}
	result := graphql.Do(graphql.Params{
		Schema:         schema,
		RequestString:  request.Query,
		OperationName:  request.OperationName,
		VariableValues: request.Variables,
		Context:        r.Context(),
	})
	writeGraphQLJSON(w, http.StatusOK, result)
}

func decodeGraphQLRequest(w http.ResponseWriter, r *http.Request) (graphQLRequest, error) {
	if r.Method == http.MethodGet {
		return graphQLRequest{Query: r.URL.Query().Get("query"), OperationName: r.URL.Query().Get("operationName")}, nil
	}
	if r.Method != http.MethodPost {
		return graphQLRequest{}, fmt.Errorf("only GET and POST are supported")
	}
	mediaType := strings.ToLower(strings.TrimSpace(strings.Split(r.Header.Get("Content-Type"), ";")[0]))
	if mediaType == "application/graphql" {
		body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
		if err != nil { return graphQLRequest{}, fmt.Errorf("GraphQL request is too large or unreadable") }
		if strings.TrimSpace(string(body)) == "" { return graphQLRequest{}, fmt.Errorf("query is required") }
		return graphQLRequest{Query: string(body)}, nil
	}
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	var request graphQLRequest
	if err := decoder.Decode(&request); err != nil {
		return graphQLRequest{}, fmt.Errorf("request must contain a valid GraphQL JSON document")
	}
	if strings.TrimSpace(request.Query) == "" {
		return graphQLRequest{}, fmt.Errorf("query is required")
	}
	return request, nil
}

func writeGraphQLJSON(w http.ResponseWriter, status int, value interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(value)
}

func (s *Server) graphQLSchema() (graphql.Schema, error) {
	coordinates := graphql.NewObject(graphql.ObjectConfig{
		Name: "Coordinates",
		Fields: graphql.Fields{
			"latitude":  &graphql.Field{Type: graphql.Float, Resolve: func(p graphql.ResolveParams) (interface{}, error) { return parseCoordinate(p.Source.(database.Zipcode).Latitude) }},
			"longitude": &graphql.Field{Type: graphql.Float, Resolve: func(p graphql.ResolveParams) (interface{}, error) { return parseCoordinate(p.Source.(database.Zipcode).Longitude) }},
		},
	})
	zipType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Zipcode",
		Fields: graphql.Fields{
			"zip":       &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) { return p.Source.(database.Zipcode).Zip, nil }},
			"zipcode":   &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) { return p.Source.(database.Zipcode).Zip, nil }},
			"zipCode":   &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) { return p.Source.(database.Zipcode).ZipCode, nil }},
			"city":      &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) { return p.Source.(database.Zipcode).City, nil }},
			"state":     &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) { return p.Source.(database.Zipcode).State, nil }},
			"stateFull": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) { return p.Source.(database.Zipcode).StateFull, nil }},
			"county":    &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) { return p.Source.(database.Zipcode).County, nil }},
			"timezone":  &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) { return p.Source.(database.Zipcode).Timezone, nil }},
			"areaCodes": &graphql.Field{Type: graphql.NewList(graphql.String), Resolve: func(p graphql.ResolveParams) (interface{}, error) { return p.Source.(database.Zipcode).AreaCodes, nil }},
			"latitude":  &graphql.Field{Type: graphql.Float, Resolve: func(p graphql.ResolveParams) (interface{}, error) { return parseCoordinate(p.Source.(database.Zipcode).Latitude) }},
			"longitude": &graphql.Field{Type: graphql.Float, Resolve: func(p graphql.ResolveParams) (interface{}, error) { return parseCoordinate(p.Source.(database.Zipcode).Longitude) }},
			"coordinates": &graphql.Field{Type: coordinates, Resolve: func(p graphql.ResolveParams) (interface{}, error) { return p.Source, nil }},
		},
	})
	statsType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Stats",
		Fields: graphql.Fields{
			"totalZipcodes": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) { return p.Source.(map[string]interface{})["total_zipcodes"], nil }},
			"totalCities":   &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) { return p.Source.(map[string]interface{})["total_cities"], nil }},
			"totalStates":   &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) { return p.Source.(map[string]interface{})["total_states"], nil }},
		},
	})
	root := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"zipcode": &graphql.Field{
				Type: zipType,
				Args: graphql.FieldConfigArgument{"code": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)}},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					code, ok := p.Args["code"].(string)
					if !ok || len(code) != 5 || !isDigits(code) { return nil, fmt.Errorf("code must be a five-digit ZIP code") }
					value, err := strconv.Atoi(code)
					if err != nil { return nil, fmt.Errorf("code must be a five-digit ZIP code") }
					result, err := s.db.SearchByZipCode(value)
					if err != nil { return nil, fmt.Errorf("ZIP lookup failed") }
					return result, nil
				},
			},
			"search": &graphql.Field{
				Type: graphql.NewList(zipType),
				Args: graphql.FieldConfigArgument{
					"city": &graphql.ArgumentConfig{Type: graphql.String},
					"state": &graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					city, _ := p.Args["city"].(string)
					state, _ := p.Args["state"].(string)
					var results []database.Zipcode
					var err error
					switch {
					case city != "" && state != "": results, err = s.db.SearchByStateAndCity(state, city)
					case city != "": results, err = s.db.SearchByCity(city)
					case state != "": results, err = s.db.SearchByState(state)
					default: return nil, fmt.Errorf("city or state is required")
					}
					if err != nil { return nil, fmt.Errorf("ZIP search failed") }
					return results, nil
				},
			},
			"stats": &graphql.Field{Type: statsType, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				stats, err := s.db.GetStats()
				if err != nil { return nil, fmt.Errorf("statistics unavailable") }
				return stats, nil
			}},
		},
	})
	return graphql.NewSchema(graphql.SchemaConfig{Query: root})
}

func parseCoordinate(value interface{}) (interface{}, error) {
	if value == nil { return nil, nil }
	if number, ok := value.(float64); ok { return number, nil }
	return strconv.ParseFloat(fmt.Sprint(value), 64)
}

func isDigits(value string) bool {
	for _, character := range value {
		if character < '0' || character > '9' { return false }
	}
	return true
}
