# This file is part of resource-git.
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as
# published by the Free Software Foundation, either version 3 of the
# License, or (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.

import os
import shutil
import tempfile
import tarfile
from urllib.parse import urljoin

import requests
from git import Repo
from sanic import Sanic, response

app = Sanic(name="bob_resource_git")
PORT = 8000


def tar_dir(tarfile_name, source_dir):
    with tarfile.open(tarfile_name, "w") as tar:
        oldpwd = os.getcwd()
        os.chdir(source_dir)

        for f in os.listdir("."):
            tar.add(f)

        os.chdir(oldpwd)


@app.route("/bob_resource")
async def handle(request):
    repo = request.args.get("repo")
    branch = request.args.get("branch")

    if any(map(lambda _: _ is None, [repo, branch])):
        return response.text("Invalid params", status=400)

    archive = "repo.tar"

    if os.path.exists(archive):
        os.remove(archive)

    clone_dir = tempfile.mkdtemp()

    Repo.clone_from(repo, clone_dir, branch=branch)

    tar_dir(archive, clone_dir)

    shutil.rmtree(clone_dir)

    return await response.file_stream(archive)


@app.route("/ping")
async def handle_ping(_):
    return response.text("Ack")


@app.route("/register")
async def handle_register(request):
    host = request.args.get("host") or "http://localhost:7777"
    path = "/api/external-resource/git"
    url = urljoin(host, path)
    data = {"url": f"http://localhost:{PORT}"}

    requests.post(url, json=data)

    return response.text("Ok")


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=PORT, workers=os.cpu_count() + 1)
