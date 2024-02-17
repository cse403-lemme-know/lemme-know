# Contributing

## Introduction

Thank you for your interest in contributing to LemmeKnow! This guide contains the information you'll need to get started.

## Reporting bugs and proposing features

Please use our [issue tracker](https://github.com/cse403-lemmeknow/lemmeknow/issues) to report bugs and propose features.

Include details and context about the issue, using the corresponding template.

## Contributing code

We accept [pull requests](https://github.com/cse403-lemmeknow/lemmeknow/pulls)!

Keep your pull request tagged as a draft until you want it reviewed. Ensure automated tests pass!

A common cause for CI failures is improper formatting. To fix, run `make fmt` in the applicable directories and try again.

You can run `make lint` in [`frontend/`](./frontend/) or [`backend/`](./backend/) to proactively check for other errors that would be flagged by the CI.

## Contributing tests

### Backend
You may add tests for `filename.go` by adding `filename_test.go`, consisting of functions with the `Test` prefix and the appropriate arguments. Full instructions can be found [here](https://go.dev/doc/tutorial/add-a-test).

One unit test function should test one implementation function, possible with multiple argument variations or for multiple output or side effect criteria. Integration tests are primarily in `main_test.go`.

### Frontend

You may add unit tests for functions in `filename.js` by adding `filename.test.js`, in which you call the `test` function for each case. Full instructions can be found [here](https://vitest.dev/guide/).

You may add tests of the form `frontend/tests/featurename.spec.js`, that make various assertions about rendered pages. Full instructions can be found [here](https://playwright.dev/docs/writing-tests). 

## License

By contributing, you agree that your contributions will be licensed under our [AGPL-3.0 license](./LICENSE).