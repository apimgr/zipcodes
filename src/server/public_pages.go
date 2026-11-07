package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type publicPage struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"paragraphs"`
	GeoIPAttribution bool `json:"geoip_attribution,omitempty"`
}

var publicPageTemplate = template.Must(template.New("public-page").Parse(`<!doctype html>
<html lang="en" data-theme="dark">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <meta name="color-scheme" content="dark light">
  <title>{{.Title}} · Zipcodes</title>
  <style>
    :root { color-scheme: dark light; font-family: system-ui, sans-serif; }
    body { margin: 0 auto; max-width: 48rem; padding: 1.5rem; line-height: 1.6; }
    main { padding: 1rem 0; }
    .skip-link { position: absolute; left: -10000px; }
    .skip-link:focus { left: 1rem; top: 1rem; padding: .75rem; background: Canvas; color: CanvasText; }
    @media (max-width: 40rem) { body { padding: 1rem; } }
    @media (prefers-reduced-motion: reduce) { *, *::before, *::after { scroll-behavior: auto !important; transition: none !important; } }
  </style>
</head>
<body>
  <a class="skip-link" href="#main-content">Skip to main content</a>
  <header><p>Zipcodes API</p></header>
  <main id="main-content">
    <h1>{{.Title}}</h1>
    {{range .Paragraphs}}<p>{{.}}</p>{{end}}
    {{if .GeoIPAttribution}}<p><a href="https://db-ip.com/">IP Geolocation by DB-IP</a> (CC BY 4.0).</p><p>Country and ASN data licensed CC BY 4.0 by the Number Resource Organization (NRO).</p>{{end}}
  </main>
  <footer><p>ZIP data from <a href="https://github.com/ReadyAPIs-com/curated-us-zips">Ready APIs curated-us-zips</a>, sourced from U.S. Census Bureau, USPS, and SimpleMaps; licensed under <a href="https://creativecommons.org/licenses/by/4.0/">CC BY 4.0</a>.</p></footer>
</body>
</html>
`))

func (s *Server) handleAbout(w http.ResponseWriter, r *http.Request) {
	s.writePublicPage(w, r, publicPage{
		Title: "About Zipcodes",
		Paragraphs: []string{
			"Zipcodes provides public, read-only lookups for United States ZIP codes.",
			"The embedded dataset contains 33,100 geographic ZIP records with city, state, county, coordinates, and IANA timezone. It is generated from Ready APIs curated-us-zips at a pinned source revision under CC BY 4.0.",
			"Attribution: ZIP data from Ready APIs curated-us-zips, sourced from U.S. Census Bureau, USPS, and SimpleMaps.",
			"The source excludes military APO/FPO/DPO ZIPs and does not provide telephone area codes. The area_codes field is empty until a redistributable ZIP-to-area-code mapping is available.",
			"No account or API key is required to use the lookup API.",
		},
		GeoIPAttribution: true,
	})
}

func (s *Server) handleHelp(w http.ResponseWriter, r *http.Request) {
	s.writePublicPage(w, r, publicPage{
		Title: "Help",
		Paragraphs: []string{
			"Look up a ZIP code with GET /api/v1/zipcodes/{code}; search by city or state using the web form.",
			"Use GET /api/v1/zipcodes/nearby?zip=90210&radius=10 for a bounded radius search, or provide lat and lng coordinates.",
			"The CLI requires --server URL unless its server is configured in ZIPCODES_SERVER_PRIMARY or cli.yml.",
			"Search by city or state with GET /api/v1/zipcodes/city/{city} and GET /api/v1/zipcodes/state/{state}.",
			"The API is public and read-only. Use the OpenAPI page for endpoint documentation.",
		},
	})
}

func (s *Server) handlePrivacy(w http.ResponseWriter, r *http.Request) {
	s.writePublicPage(w, r, publicPage{
		Title: "Privacy",
		Paragraphs: []string{
			"ZIP code records contain public geographic reference data and do not contain personal information.",
			"GeoIP results are computed for a request and are not saved in the ZIP code database. If GeoIP data is unavailable, location fields are omitted.",

			"The server operator controls deployment and access logs. Contact that operator for questions about log retention.",
		},
	})
}

func (s *Server) handleTerms(w http.ResponseWriter, r *http.Request) {
	s.writePublicPage(w, r, publicPage{
		Title: "Terms of Use",
		Paragraphs: []string{
			"The API provides geographic reference information for convenience. Validate important address and delivery decisions with the appropriate postal authority.",
			"The service is provided without warranty. The server operator may change or discontinue access to a self-hosted instance.",
		},
	})
}

func (s *Server) handleHealthPage(w http.ResponseWriter, r *http.Request) {
	if _, err := s.db.GetStats(); err != nil {
		s.writePublicPageStatus(w, r, http.StatusServiceUnavailable, publicPage{
			Title:      "Service status",
			Paragraphs: []string{"ERROR: The service is unavailable."},
		})
		return
	}
	s.writePublicPage(w, r, publicPage{
		Title:      "Service status",
		Paragraphs: []string{"OK: Zip code lookup is available."},
	})
}

func (s *Server) writePublicPage(w http.ResponseWriter, r *http.Request, page publicPage) {
	s.writePublicPageStatus(w, r, http.StatusOK, page)
}

func (s *Server) writePublicPageStatus(w http.ResponseWriter, r *http.Request, status int, page publicPage) {
	accept := strings.ToLower(r.Header.Get("Accept"))
	if strings.Contains(accept, "text/plain") && !strings.Contains(accept, "application/json") {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(status)
		if status >= http.StatusBadRequest { _, _ = fmt.Fprintln(w, "ERROR: The service is unavailable."); return }
		_, _ = fmt.Fprintln(w, page.Title)
		for _, paragraph := range page.Paragraphs { _, _ = fmt.Fprintln(w, paragraph) }
		return
	}
	if strings.HasPrefix(r.URL.Path, "/api/") || strings.Contains(accept, "application/json") {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(status)
		if status >= http.StatusBadRequest {
			_ = json.NewEncoder(w).Encode(map[string]any{"ok": false, "error": "SERVICE_UNAVAILABLE", "message": "The service is unavailable."})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true, "data": page})
		return
	}
	if strings.Contains(r.Header.Get("Accept"), "text/plain") {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(status)
		_, _ = fmt.Fprintf(w, "%s\n%s\n", page.Title, strings.Join(page.Paragraphs, "\n"))
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := publicPageTemplate.Execute(w, page); err != nil {
		return
	}
}
