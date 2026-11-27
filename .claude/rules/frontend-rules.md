# Frontend Rules (PART 16)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Never make JavaScript necessary for a core page or form flow.
- Never render private or visitor-specific content into a shared offline cache.
- Never use inline styles/scripts, inaccessible controls, or stale asset URLs.

## CRITICAL - ALWAYS DO
- Render pages server-side through the shared accessible shell and registry.
- Keep mobile, keyboard, reduced-motion, dark/light/auto theme, and offline behavior usable.
- Cache only explicitly marked public pages; synthesize a real offline Response on failures.
- Keep browser pages, JSON/text forms, SEO, announcements, and preferences consistent.

For complete details, see AI.md PART 16.
