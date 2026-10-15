# GeoIP module

The server downloads local MMDB files from `sapics/ip-location-db` on first run. City records come from DB-IP; country and ASN records come from the Number Resource Organization (NRO). All three sources are CC BY 4.0. IDEA.md describes this source, following AI.md PART 19 license and source rules.

The default directory is `{data_dir}/security/geoip`. Set `server.geoip.dir` to change it. Set `server.geoip.enabled: false` to disable downloads and lookups; the `asn`, `country`, and `city` flags independently select database categories.

The lookup is local and makes no request to a third-party API. It returns only fields available in the selected data files. Private, loopback, unspecified, and link-local addresses are not looked up. Missing databases do not block other application requests.

Visible attribution is included on `/server/about`:

- [IP Geolocation by DB-IP](https://db-ip.com/), CC BY 4.0
- Country and ASN data licensed CC BY 4.0 by the Number Resource Organization (NRO)

The application refreshes enabled databases on the configured schedule. The next run time and the most recent 1,000 task results persist in the SQLite database.
