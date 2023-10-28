package main

import (
	"bytes"
	"fmt"
	"github.com/go-git/go-billy/v5/memfs"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"os"
)

func cloneGitRepo() {
	// create a new in-memory filesystem
	fs := memfs.New()

	// create a new Git memory storage
	memStorage := memory.NewStorage()

	// create the authentication method
	auth := &http.BasicAuth{
		Username: "dawnbreather",
		Password: "glpat-fPRu9o4cxWYxQDqmsZSF",
	}

	// clone the repository into the in-memory filesystem using the authentication method
	repo, err := git.Clone(memStorage, fs, &git.CloneOptions{
		URL:      "https://gitlab.com/dawnbreather/test-tmp.git",
		Auth:     auth,
		Progress: os.Stdout,
	})
	if err != nil && err != transport.ErrEmptyRemoteRepository {
		fmt.Println("Error cloning repository:", err)
		os.Exit(1)
	}

	// checkout the specific branch
	wt, err := repo.Worktree()
	if err != nil {
		fmt.Println("Error getting worktree:", err)
		os.Exit(1)
	}

	err = wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName("main"),
	})
	if err != nil {
		fmt.Println("Error checking out branch:", err)
		os.Exit(1)
	}

	// modify a file
	file, err := fs.OpenFile("README.md", os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(file); err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	if _, err := file.Seek(0, 0); err != nil {
		fmt.Println("Error seeking file:", err)
		os.Exit(1)
	}

	if _, err := file.Write([]byte("modified content")); err != nil {
		fmt.Println("Error writing file:", err)
		os.Exit(1)
	}

	// commit the changes
	_, err = wt.Commit("commit message", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Automated Deployer",
			Email: "devops@abcloudz.com",
		},
	})
	if err != nil {
		fmt.Println("Error committing changes:", err)
		os.Exit(1)
	}

	// push the changes to the remote repository
	err = repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       auth,
		Progress:   os.Stdout,
	})
	if err != nil {
		fmt.Println("Error pushing changes:", err)
		os.Exit(1)
	}

	fmt.Println("Changes pushed successfully")

	err = wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName("dev"),
	})
	if err != nil {
		fmt.Println("Error checking out branch:", err)
		os.Exit(1)
	}

	// modify a file
	file, err = fs.OpenFile("README.md", os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	buf = new(bytes.Buffer)
	if _, err := buf.ReadFrom(file); err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	if _, err := file.Seek(0, 0); err != nil {
		fmt.Println("Error seeking file:", err)
		os.Exit(1)
	}

	if _, err := file.Write([]byte("modified content")); err != nil {
		fmt.Println("Error writing file:", err)
		os.Exit(1)
	}

	// commit the changes
	_, err = wt.Commit("commit message", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Automated Deployer",
			Email: "devops@abcloudz.com",
		},
	})
	if err != nil {
		fmt.Println("Error committing changes:", err)
		os.Exit(1)
	}

	// push the changes to the remote repository
	err = repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       auth,
		Progress:   os.Stdout,
	})
	if err != nil {
		fmt.Println("Error pushing changes:", err)
		os.Exit(1)
	}

	fmt.Println("Changes pushed successfully")

}
