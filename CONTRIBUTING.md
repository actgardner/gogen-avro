# Contributing to gogen-avro

Want to fix a bug or add a feature to gogen-avro? That's awesome!
Currently the project has a single maintainer, and I work on it part-time along with my day job, so it can be hard to review large or complex PRs.
I do try to review and give feedback on every PR in a timely manner.

Some tips to get PRs merged quickly:

- **Keep PRs small and focused**. If your PR adds a new feature, try to only add that feature and avoid making other changes at the same time (like refactoring surrounding code).
Small, focused PRs that do a single thing are easier to discuss and accept.

- **Refactor first if necessary, in a separate PR**. If there's some deep architectural changes required to add a new feature or fix a bug, try and split them into a separate PR that doesn't change any existing functionality first.
PRs that only clean up or reorganize the codebase are fine. Try to make them as small and focused as possible, either by refactoring a single piece of the codebase, or by making one type of change everywhere.

- **Add integration tests as necessary**. If your PR fixes a bug or adds new functionality, add corresponding tests to the `tests` directory to exercise the functionality.
