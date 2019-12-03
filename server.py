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
from urllib.parse import urljoin

import requests
from git import Repo
from sanic import Sanic, response

app = Sanic()
PORT = 8000


@app.route("/bob_resource")
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


@app.route("/register")
async def handle_register(request):
    host = request.args.get("host") or "http://localhost:7777"
    path = "/api/external-resource/git"
    url = urljoin(host, path)
    data = {"url": f"http://localhost:{PORT}"}

    requests.post(url, json=data)

    return response.text("Ok")


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=PORT)
