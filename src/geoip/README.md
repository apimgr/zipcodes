# GeoIP module

The server downloads local MMDB files from `sapics/ip-location-db` on first run. City records come from DB-IP under CC BY 4.0. Country and ASN records come from iptoasn.com under PDDL. The downloader verifies each file against its upstream SHA-256 checksum before replacing a live database.

The default directory is `{data_dir}/security/geoip`. Set `server.geoip.dir` to change it. Set `server.geoip.enabled: false` to disable downloads and lookups; the `asn`, `country`, and `city` flags independently select database categories.

The lookup is local and makes no request to a third-party API. It returns only fields available in the selected data files. Private, loopback, unspecified, and link-local addresses are not looked up. Missing databases do not block other application requests.

Visible attribution is included on `/server/about`:

- [IP Geolocation by DB-IP](https://db-ip.com/)

Country and ASN data are PDDL and require no attribution.

The application refreshes enabled databases on the configured schedule. The next run time and the most recent 1,000 task results persist in the SQLite database.
