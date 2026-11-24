# Translation guide

The application uses JSON locale catalogs embedded in its Go binaries. English (`en.json`) is the source catalog; translators work in a copy for their language and never translate or rename keys. The current catalog set has six languages: English (`en`), Spanish (`es`), French (`fr`), Arabic (`ar`), German (`de`), and Italian (`it`). The supported-language registry and the locale files must stay in sync.

## Adding or updating a translation

For a new language, copy `src/common/i18n/locales/en.json` to `src/common/i18n/locales/{language}.json`, using the lowercase language code. Translate every value and preserve the JSON structure and key order. Fill the `meta` fields (`language`, `name`, `native_name`, `direction`, and `version`), then add the language to the registry with its text direction and CLDR plural categories. Existing translations should change only the locale file; update its metadata when needed.

Locale files use UTF-8 JSON with two-space indentation, no comments, and no trailing commas. Keep keys unchanged. Preserve every interpolation placeholder exactly, including its spelling and braces, such as `{count}`, `{current}`, `{total}`, and `{datetime}`. Provide the plural forms required by the language. Each value should be a complete phrase rather than a fragment intended to be joined to another string.

Use a consistent, neutral register. Follow the language's punctuation and typography conventions, including French spacing before punctuation where required. Keep product and brand names, protocol and technology names, code, file names, configuration keys, URLs, and literal values users enter (such as `dark` and `auto`) untranslated.

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

Check that the locale has the same keys, placeholders, and required plural forms as English. Review the interface at 320, 768, and 1280 pixels; confirm translated text is not clipped or overflowing. Check right-to-left languages in a mirrored layout. Languages whose scripts may not be covered by the system font need an appropriate font-stack update and a long unspaced text wrapping check.

AI-assisted translation is welcome when it can produce natural, reliable wording. Human corrections take precedence. A maintainer reviews key and placeholder parity, glossary use, and layout before merging. Native-speaker review is welcome but is not required. Credit translators by their chosen name or handle in this guide and in the locale file's optional `meta.translators` field. Do not invent contributor credits; add them when contributors provide their preferred attribution.

For a contribution, submit one locale file for an existing language. For a new language, include its locale file and registry changes in the same contribution. Keep all locale files synchronized whenever English strings change.
