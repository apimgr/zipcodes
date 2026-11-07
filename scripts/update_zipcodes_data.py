#!/usr/bin/env python3
"""Fetch and normalize the pinned CC BY 4.0 ZIP dataset."""

import csv
import io
import json
import urllib.request
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
COMMIT = "64ac28b620a39d9efabe39f91dc03181a13a8a76"
SOURCE = (
    "https://raw.githubusercontent.com/ReadyAPIs-com/curated-us-zips/"
    f"{COMMIT}/data/us-zips.csv"
)
STATE_NAMES = {
    "AL": "Alabama", "AK": "Alaska", "AZ": "Arizona", "AR": "Arkansas",
    "CA": "California", "CO": "Colorado", "CT": "Connecticut", "DE": "Delaware",
    "DC": "District of Columbia", "FL": "Florida", "GA": "Georgia", "HI": "Hawaii",
    "ID": "Idaho", "IL": "Illinois", "IN": "Indiana", "IA": "Iowa",
    "KS": "Kansas", "KY": "Kentucky", "LA": "Louisiana", "ME": "Maine",
    "MD": "Maryland", "MA": "Massachusetts", "MI": "Michigan", "MN": "Minnesota",
    "MS": "Mississippi", "MO": "Missouri", "MT": "Montana", "NE": "Nebraska",
    "NV": "Nevada", "NH": "New Hampshire", "NJ": "New Jersey", "NM": "New Mexico",
    "NY": "New York", "NC": "North Carolina", "ND": "North Dakota", "OH": "Ohio",
    "OK": "Oklahoma", "OR": "Oregon", "PA": "Pennsylvania", "RI": "Rhode Island",
    "SC": "South Carolina", "SD": "South Dakota", "TN": "Tennessee", "TX": "Texas",
    "UT": "Utah", "VT": "Vermont", "VA": "Virginia", "WA": "Washington",
    "WV": "West Virginia", "WI": "Wisconsin", "WY": "Wyoming", "PR": "Puerto Rico",
    "VI": "U.S. Virgin Islands", "GU": "Guam", "AS": "American Samoa",
    "MP": "Northern Mariana Islands",
}


def main():
    with urllib.request.urlopen(SOURCE, timeout=30) as response:
        content = response.read().decode("utf-8")
    records = []
    for row in csv.DictReader(io.StringIO(content)):
        code = row["zip_code"]
        latitude = row["latitude"]
        longitude = row["longitude"]
        state = row["state"]
        records.append({
            "zip": code,
            "zip_code": int(code),
            "city": row["city"],
            "state": state,
            "state_full": STATE_NAMES.get(state, ""),
            "county": row["county"],
            "latitude": float(latitude) if latitude else None,
            "longitude": float(longitude) if longitude else None,
            "timezone": row["timezone"],
            "area_codes": [],
        })
    if len(records) != 33100:
        raise ValueError(f"expected 33100 source records, received {len(records)}")
    if any(not record["state_full"] for record in records):
        raise ValueError("source contains an unmapped state or territory code")
    if any(not record["city"] or not record["county"] or not record["timezone"] for record in records):
        raise ValueError("source contains an incomplete geographic record")
    codes = [record["zip"] for record in records]
    if len(set(codes)) != len(codes):
        raise ValueError("source contains duplicate ZIP codes")
    destination = ROOT / "src" / "data" / "zipcodes.json"
    destination.write_text(json.dumps(records, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    print(f"Wrote {len(records)} ZIP records to {destination}")


if __name__ == "__main__":
    main()
