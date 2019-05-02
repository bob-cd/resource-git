### Reference Bob Resource: Git

This is a simple external resource enabling Bob to read git repositories.

#### Requirements
- Python 3.5+
- Git 1.7+
- [Poetry](https://poetry.eustace.io)

#### Running
- `poetry install` to install dependencies.
- `python3 server.py` will start the plugin on port 8000.

#### API
- `GET /bob_request`: Takes `repo` and `branch` as params, clones and
   responds back with a zip of the repo.
- `GET /ping`: Responds with an `Ack`.
