# tinyclaw: Repo-Specific Review Rules

Inherits all rules from `_org-enforced-rules.md`.

## General

- Keep the utility focused and minimal. Flag scope creep in PRs.
- Error messages must be user-friendly. No raw panics, stack traces, or cryptic exit codes.
- CLI argument parsing must validate inputs and provide help text.
- New functionality must include tests.

## Code Quality

- Prefer standard library over external dependencies for simple operations.
- Functions should do one thing. Flag functions over ~50 lines for potential splitting.
- Handle edge cases: empty input, missing files, permission errors.
