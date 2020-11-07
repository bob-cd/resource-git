/* This file is part of resource-git.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"archive/tar"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/client"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

func Tar(src string, archiveName string) error {
	if _, err := os.Stat(src); err != nil {
		return err
	}

	tarfile, err := os.Create(archiveName)
	if err != nil {
		return err
	}
	defer tarfile.Close()

	var fileWriter io.WriteCloser = tarfile

	tw := tar.NewWriter(fileWriter)
	defer tw.Close()

	err = filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !fi.Mode().IsRegular() {
			return nil
		}
		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}
		header.Name = strings.TrimPrefix(strings.Replace(file, src, "", -1), string(filepath.Separator))

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		f, err := os.Open(file)
		if err != nil {
			return err
		}

		if _, err := io.Copy(tw, f); err != nil {
			return err
		}

		f.Close()

		return nil
	})

	return nil
}

func Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Ack")
}

func Clone(w http.ResponseWriter, r *http.Request) {
	repo := r.URL.Query().Get("repo")
	branch := r.URL.Query().Get("branch")

	if repo == "" || branch == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid params: both repo and branch required.")

		return
	}

	archive := fmt.Sprintf("%d", time.Now().UnixNano())
	dir := "repo-" + archive
	os.MkdirAll(dir, os.ModePerm)
	defer os.RemoveAll(dir)
	defer os.Remove(archive)

	_, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:           repo,
		ReferenceName: plumbing.ReferenceName("refs/heads/" + branch),
		SingleBranch:  true,
		Depth:         1,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())

		return
	}

	if err := Tar(dir, archive); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())

		return
	}

	http.ServeFile(w, r, archive)
}

func main() {
	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "8000"
	}

	customClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	client.InstallProtocol("https", githttp.NewClient(customClient))

	http.HandleFunc("/ping", Ping)
	http.HandleFunc("/bob_resource", Clone)

	http.ListenAndServe(":"+port, nil)
}
