# Translation guide

The web interface currently ships English copy. A translation contribution should preserve the meaning of ZIP code fields, validation messages, privacy statements, and dataset attribution.

Before translating, search the server templates and handlers for the source text and identify every occurrence. Keep terminology consistent across the search form, result pages, API documentation UI, and server information pages. Do not translate route paths, JSON field names, ZIP codes, state abbreviations, IANA timezone identifiers, or API examples.

Review translated pages at mobile and desktop widths, with JavaScript disabled, and check that labels, errors, and navigation remain understandable to keyboard and screen-reader users. Include the language name and reviewer in the contribution description.

There is no localization framework or translation catalog in the current implementation; changes should not imply that additional languages are already available.
