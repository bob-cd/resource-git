/*
* Copyright 2018- Rahul De
*
* Use of this source code is governed by an MIT-style
* license that can be found in the LICENSE file or at
* https://opensource.org/licenses/MIT.
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

func makeTar(src string, archiveName string) error {
	if _, err := os.Stat(src); err != nil {
		return err
	}

	tarfile, err := os.Create(archiveName)
	if err != nil {
		return err
	}
	defer tarfile.Close()

	tw := tar.NewWriter(tarfile)
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

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Ack")
}

func clone(w http.ResponseWriter, r *http.Request) {
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

	if err := makeTar(dir, archive); err != nil {
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

	http.HandleFunc("/ping", ping)
	http.HandleFunc("/bob_resource", clone)

	http.ListenAndServe(":"+port, nil)
}
