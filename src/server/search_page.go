package server

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/apimgr/zipcodes/src/config"
	"github.com/apimgr/zipcodes/src/database"
)

type searchPageData struct {
	MaxRadius float64
	Query       string
	Zip         string
	Radius      string
	Latitude    string
	Longitude   string
	Error       string
	Results     []database.Zipcode
	Nearby      []database.NearbyZipcode
	ResultCount int
	Stats       map[string]interface{}
}

var searchPageTemplate = template.Must(template.New("search").Parse(`<!doctype html>
<html lang="en" data-theme="auto">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <meta name="color-scheme" content="dark light">
  <meta name="theme-color" content="#151923">
  <link rel="manifest" href="/manifest.json">
  <link rel="stylesheet" href="/static/css/main.css">
  <title>Zipcodes | United States ZIP code lookup</title>
  <style>
    body { max-width: 64rem; margin: 0 auto; padding: 1rem; line-height: 1.6; }
    main { display: grid; gap: 1.5rem; }
    form { display: grid; grid-template-columns: repeat(auto-fit, minmax(min(100%, 14rem), 1fr)); gap: .75rem; }
    label { display: grid; gap: .25rem; }
    input, select, button { min-height: 2.75rem; font: inherit; padding: .4rem .6rem; }
    .result { padding: .9rem; margin: .5rem 0; border: 1px solid currentColor; border-radius: .5rem; }
    .skip-link { position: absolute; left: -10000px; }
    .skip-link:focus { left: 1rem; top: 1rem; padding: .75rem; background: Canvas; color: CanvasText; }
    @media (prefers-color-scheme: dark) { html[data-theme="auto"] { color-scheme: dark; } }
    @media (prefers-color-scheme: light) { html[data-theme="auto"] { color-scheme: light; } }
  </style>
</head>
<body>
  <a class="skip-link" href="#main-content">Skip to main content</a>
  <header><h1>Zipcodes</h1><p>United States ZIP code lookup</p>
    <nav aria-label="Main navigation"><a href="/api/v1/server/swagger">API documentation</a> · <a href="/graphql">GraphQL</a> · <a href="/server/about">About</a></nav>
  </header>
  <main id="main-content">
    <section aria-labelledby="search-heading">
      <h2 id="search-heading">Search ZIP codes</h2>
      <form method="get" action="/search">
        <label>ZIP code, city, or state <input name="q" value="{{.Query}}" maxlength="80" autocomplete="postal-code" placeholder="90210 or Beverly Hills"></label>
        <label>Origin ZIP for radius search <input name="zip" value="{{.Zip}}" inputmode="numeric" pattern="[0-9]{5}" maxlength="5"></label>
        <label>Radius in miles <input name="radius" value="{{.Radius}}" type="number" min="0.1" max="{{.MaxRadius}}" step="0.1"></label>
        <label>Origin latitude <input name="lat" value="{{.Latitude}}" type="number" min="-90" max="90" step="any"></label>
        <label>Origin longitude <input name="lng" value="{{.Longitude}}" type="number" min="-180" max="180" step="any"></label>
        <button type="submit">Search</button>
      </form>
      <label>Color theme <select id="theme-choice" aria-label="Color theme"><option value="auto">Use device setting</option><option value="light">Light</option><option value="dark">Dark</option></select></label>
      {{if .Error}}<p role="alert">{{.Error}}</p>{{end}}
    </section>
    {{if .ResultCount}}
    <section aria-labelledby="results-heading">
      <h2 id="results-heading">{{.ResultCount}} results</h2>
      {{range .Results}}<article class="result"><h3>{{.Zip}} · {{.City}}, {{.State}}</h3><p>{{.StateFull}}{{if .County}} · {{.County}} County{{end}}</p><p>Coordinates: {{.Latitude}}, {{.Longitude}}</p><p>Timezone: {{.Timezone}}</p><p>Area codes: {{range $i, $code := .AreaCodes}}{{if $i}}, {{end}}{{$code}}{{else}}Not available{{end}}</p></article>{{end}}
      {{range .Nearby}}<article class="result"><h3>{{.Zipcode.Zip}} · {{.Zipcode.City}}, {{.Zipcode.State}}</h3><p>{{printf "%.2f" .DistanceMiles}} miles away · {{.Zipcode.County}} County · {{.Zipcode.Timezone}}</p></article>{{end}}
    </section>
    {{end}}
    <section aria-label="Dataset statistics"><h2>Dataset</h2><p>{{index .Stats "total_zipcodes"}} ZIP records · {{index .Stats "total_states"}} states · {{index .Stats "total_cities"}} cities</p></section>
  </main>
  <footer><p>ZIP data from Ready APIs curated-us-zips, sourced from U.S. Census Bureau, USPS, and SimpleMaps. Telephone area-code data is unavailable.</p></footer>
  <script src="/static/js/theme.js" defer></script>
</body>
</html>
`))

func (s *Server) handleSearchPage(w http.ResponseWriter, r *http.Request) {
	maxRadius := config.Get().Server.MaxRadiusMiles
	if maxRadius <= 0 {
		maxRadius = config.DefaultConfig().Server.MaxRadiusMiles
	}
	data := searchPageData{
		Query: r.URL.Query().Get("q"), Zip: r.URL.Query().Get("zip"),
		Radius: r.URL.Query().Get("radius"), MaxRadius: maxRadius, Latitude: r.URL.Query().Get("lat"),
		Longitude: r.URL.Query().Get("lng"),
	}
	if data.Radius != "" {
		data.loadNearby(s.db, r)
	} else if strings.TrimSpace(data.Query) != "" {
		data.loadTextSearch(s.db)
	}
	data.Stats, _ = s.db.GetStats()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := searchPageTemplate.Execute(w, data); err != nil {
		return
	}
}

func (d *searchPageData) loadNearby(db *database.DB, r *http.Request) {
	radius, err := strconv.ParseFloat(d.Radius, 64)
	if err != nil { d.Error = "Enter a valid radius."; return }
	var latitude, longitude float64
	if d.Zip != "" {
		if len(d.Zip) != 5 || !onlyDigits(d.Zip) { d.Error = "Enter a five-digit origin ZIP."; return }
		origin, err := strconv.Atoi(d.Zip)
		if err != nil { d.Error = "Enter a five-digit origin ZIP."; return }
		zipcode, err := db.SearchByZipCode(origin)
		if err != nil || zipcode == nil { d.Error = "Origin ZIP was not found."; return }
		latitude, err = strconv.ParseFloat(zipcode.GetLatitudeString(), 64)
		if err != nil { d.Error = "Origin ZIP has no coordinates."; return }
		longitude, err = strconv.ParseFloat(zipcode.GetLongitudeString(), 64)
		if err != nil { d.Error = "Origin ZIP has no coordinates."; return }
	} else {
		latitude, err = strconv.ParseFloat(d.Latitude, 64)
		if err != nil { d.Error = "Enter an origin ZIP or valid coordinates."; return }
		longitude, err = strconv.ParseFloat(d.Longitude, 64)
		if err != nil { d.Error = "Enter an origin ZIP or valid coordinates."; return }
	}
	maxRadius := config.Get().Server.MaxRadiusMiles
	if maxRadius <= 0 {
		maxRadius = config.DefaultConfig().Server.MaxRadiusMiles
	}
	if radius <= 0 || radius > maxRadius {
		d.Error = "Radius exceeds the configured maximum."
		return
	}
	results, err := db.SearchWithinRadius(latitude, longitude, radius)
	if err != nil { d.Error = "Radius search could not be completed."; return }
	d.Nearby = results
	d.ResultCount = len(results)
}

func (d *searchPageData) loadTextSearch(db *database.DB) {
	query := strings.TrimSpace(d.Query)
	var results []database.Zipcode
	var err error
	if len(query) == 5 && onlyDigits(query) {
		code, parseErr := strconv.Atoi(query)
		if parseErr != nil { d.Error = "Enter a valid five-digit ZIP code."; return }
		var result *database.Zipcode
		result, err = db.SearchByZipCode(code)
		if result != nil { results = []database.Zipcode{*result} }
	} else if len(query) == 2 && onlyLetters(query) {
		results, err = db.SearchByState(query)
	} else if strings.Contains(query, ",") {
		parts := strings.SplitN(query, ",", 2)
		results, err = db.SearchByStateAndCity(strings.TrimSpace(parts[1]), strings.TrimSpace(parts[0]))
	} else {
		results, err = db.SearchByCity(query)
	}
	if err != nil { d.Error = "Search could not be completed."; return }
	d.Results = results
	d.ResultCount = len(results)
}

func onlyDigits(value string) bool {
	if value == "" { return false }
	for _, character := range value { if character < '0' || character > '9' { return false } }
	return true
}

func onlyLetters(value string) bool {
	if value == "" { return false }
	for _, character := range value { if (character < 'A' || character > 'Z') && (character < 'a' || character > 'z') { return false } }
	return true
}
