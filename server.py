import os
import shutil
import tempfile

from git import Repo
from sanic import Sanic, response

app = Sanic()


@app.route("/bob_request")
async def handle(request):
    repo = request.args.get("repo")
    branch = request.args.get("branch")

    if any(map(lambda _: _ is None, [repo, branch])):
        return response.text("Invalid params", status=400)

    zip_file = "repo.zip"

    if os.path.exists(zip_file):
        os.remove(zip_file)

    clone_dir = tempfile.mkdtemp()

    Repo.clone_from(repo, clone_dir, branch=branch)

    shutil.make_archive("repo", "zip", clone_dir)

    shutil.rmtree(clone_dir)

    return await response.file(zip_file)


@app.route("/ping")
async def handle_ping(_):
    return response.text("Ack")


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8000)
