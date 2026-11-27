# Frontend Rules (PART 16)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Never replace server-rendered pages or core forms with client-side rendering.
- Never use inline CSS, JavaScript alerts, or gratuitous JavaScript.
- Never hide required functionality behind JavaScript.
- Never use non-semantic controls or remove visible keyboard focus.

## CRITICAL - ALWAYS DO
- Render accessible, mobile-first pages with semantic HTML and no-JS core flows.
- Default to dark theme; support light/dark/auto with server-rendered theme state.
- Meet WCAG AA contrast, keyboard, screen-reader, focus, reduced-motion, and RTL requirements.
- Use external assets, native HTML disclosure/dialog/form features, and localized text.
- Keep public pages and forms consistent with canonical API and product routes.

## KEY DECISIONS
| Question | Answer | Spec Reference |
|---|---|---|
| How does a no-JS visitor use core pages? | Standard links and POST forms rendered by the server. | PART 16 |
| What is the default theme? | Dark. | PART 16 |

For complete details, see AI.md PART 16.
