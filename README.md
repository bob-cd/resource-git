### Reference Bob Resource: Git

This is a simple external resource enabling Bob to read git repositories.

#### Requirements
- [Go](https://golang.org/dl/) 1.14+

#### Running
- `go build main.go` to compile the code and obtain a binary `main`.
- `./main` will start on port `8000` by default, set the env var `PORT` to change.

#### API
- `GET /bob_resource`: Takes `repo` and `branch` as params, clones and
   responds back with a tar of the repo.
- `GET /ping`: Responds with an `Ack`.
