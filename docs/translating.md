# Translation guide

The application uses JSON locale catalogs embedded in its Go binaries. English (`en.json`) is the source catalog; translators work in a copy for their language and never translate or rename keys. The current catalog set has six languages. The supported-language registry, locale files, and this list must stay in sync.

| Code | Language | Direction |
|---|---|---|
| `en` | English | `ltr` |
| `es` | Spanish | `ltr` |
| `fr` | French | `ltr` |
| `ar` | Arabic | `rtl` |
| `de` | German | `ltr` |
| `it` | Italian | `ltr` |

## Adding or updating a translation

For a new language, copy `src/common/i18n/locales/en.json` to `src/common/i18n/locales/{language}.json`, using its lowercase BCP 47 language code (for example `pt` or `pt-br`). Translate every value and preserve the JSON structure and key order. Fill the `meta` fields (`language`, `name`, `native_name`, `direction`, and `version`), then add the language to the single language registry with its text direction and CLDR plural categories. Register it in the configuration defaults as required by the implementation contract. The locale file, registry entry, and supported-language documentation must be updated together. Existing translations should change only the locale file; update its metadata when needed.

Locale files use UTF-8 JSON with two-space indentation, no comments, and no trailing commas. Keep keys unchanged. Preserve every interpolation placeholder exactly, including its spelling and braces, such as `{count}`, `{current}`, `{total}`, and `{datetime}`. Provide the plural forms required by the language. Each value should be a complete phrase rather than a fragment intended to be joined to another string.

Every locale must have the same message-key set as English. For each plural message, provide the CLDR categories used by that language and always include `other`; an optional `zero` form is allowed. `{count}` may be omitted in a plural category when its wording does not need the number, but all other placeholders must appear with the same spelling in every category. A key, required category, or placeholder mismatch fails the locale validation tests.

The required plural categories for the current catalog set are:

| Language | Categories |
|---|---|
| English (`en`), Spanish (`es`), German (`de`) | `one`, `other` |
| French (`fr`) | `one`, `other` (zero uses `one`) |
| Arabic (`ar`) | `zero`, `one`, `two`, `few`, `many`, `other` |
| Italian (`it`) | `one`, `many`, `other` (`many` applies to round multiples of 1,000,000) |

Use a consistent, neutral register. Follow the language's punctuation and typography conventions, including French spacing before punctuation where required. Keep product and brand names, protocol and technology names, code, file names, configuration keys, URLs, and literal values users enter (such as `dark` and `auto`) untranslated.

Translate complete sentences and messages, not sentence fragments intended to be joined together. Keep the license-required attribution text and link text intact. The service worker's English offline fallback is used only before a page has saved translated offline strings; keep the offline title, heading, message, and reload label synchronized in every locale.

## Glossary

Use one term consistently throughout a locale. Preserve these concepts and the English technical terms where a localized equivalent would make an identifier or protocol unclear:

| English term | Translation guidance |
|---|---|
| ZIP code | Use the established postal-code term in the target language; keep `ZIP` in API fields and code. |
| API | Keep `API`. |
| token | Use one term for an API credential; do not translate token strings or field names. |
| primary token | Keep distinct from a linked token; preserve the primary credential's administrative meaning. |
| linked token | Use the same term wherever a token is linked to an owner. |
| rotate / revoke | Keep these actions distinct: rotation replaces a credential; revocation disables it. |
| theme | Use one term for the visual appearance preference. |
| language | Use one term for the selected interface language. |
| owner | Distinguish the token owner from an anonymous visitor; do not imply account registration. |

When the source phrase is unclear or incomplete, report it for correction instead of translating around the ambiguity.

## Review and credit

The server and client use the same supported language set. A language must appear in the server's preferences selector, the client `--lang` option, and the supported-language documentation only when its embedded locale and registry entry are present. Arabic requires a mirrored right-to-left layout. Use logical CSS properties such as `margin-inline-start` and `text-align: start`; avoid left/right-only layout rules.

Review the interface at 320, 768, and 1280 pixels; confirm translated text is not clipped or overflowing. Check right-to-left languages in a mirrored layout. Languages whose scripts may not be covered by the system font need an appropriate font-stack update and a long unspaced text wrapping check. Dates and times use each locale's `format.date`, `format.time`, and `format.datetime` patterns; numbers, percentages, and currencies use locale-aware formatting rather than translated punctuation in strings.

AI-assisted translation is welcome when it can produce natural, reliable wording. Human corrections take precedence. A maintainer reviews key and placeholder parity, plural categories, glossary use, and layout before merging. Native-speaker review is welcome but is not required. Credit translators by their chosen name or handle in this guide and in the locale file's optional `meta.translators` field. Do not invent contributor credits; add them when contributors provide their preferred attribution.

Translate a new string into every language where the wording can be translated reliably in the same change as the English source. If a language is outside the translator's reliable fluency, leave that translation for a human contributor instead of shipping a guess. Human corrections remain authoritative in later translation updates.

Locale files also define `format.date`, `format.time`, and `format.datetime` using Go reference-time layouts. These formats, optional `format.am` / `format.pm` labels, and optional `format.digits` are data for formatting routines, not sentence strings. Numbers, percentages, and currency use locale-aware formatting libraries; do not hard-code translated punctuation into UI text.

For a contribution, submit one locale file for an existing language. For a new language, include its locale file and registry/configuration changes in the same contribution. Keep every locale synchronized whenever English strings change. The project test suite validates locale key parity, placeholder parity, and plural categories as part of `make test`; it does not use a separate translation-validation Make target.
