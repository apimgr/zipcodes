# GeoIP module

The server downloads local MMDB files from `sapics/ip-location-db` on first run. City records come from DB-IP under CC BY 4.0. Country and ASN records come from iptoasn.com under PDDL. The downloader verifies each file against its upstream SHA-256 checksum before replacing a live database.

The fixed directory is `{data_dir}/security/geoip`. Set `server.geoip.enabled: false` to disable downloads, updates, lookups, and country restrictions. The application always manages the ASN, Country, City IPv4, and City IPv6 datasets together, with a weekly update schedule.

The lookup is local and makes no request to a third-party API. It returns only fields available in the selected data files. Private, loopback, unspecified, and link-local addresses are not looked up. Missing databases do not block other application requests.

Visible attribution is included on `/server/about`:

- [IP Geolocation by DB-IP](https://db-ip.com/)

Country and ASN data are PDDL and require no attribution.

The application refreshes enabled databases on the configured schedule and records each successful update in `.last_updated`. The previous databases remain available after update failures; their health status becomes degraded only after the recorded successful update is more than 21 days old. The next run time and the most recent 1,000 task results persist in the SQLite database.
