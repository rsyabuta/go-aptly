package aptly

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testAptly() *Aptly {
	a := &Aptly{
		Url: "http://localhost:8080",
	}
	return a
}

func TestListRepos(t *testing.T) {
	aptly := testAptly()
	repos, err := aptly.ListRepos()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%q\n", repos)
}

func TestMalformedUrl(t *testing.T) {
	aptly := new(Aptly)
	aptly.Url = "localhost:8080"
	_, err := aptly.ListRepos()
	assert.Error(t, err, "Malformed URLs should error")
}

func TestCreateRepo(t *testing.T) {
	aptly := testAptly()
	repo := &LocalRepo{
		Name: "test_repo"}
	createdRepo, err := aptly.CreateRepo(repo)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", createdRepo)
}

func TestDeleteRepo(t *testing.T) {
	aptly := testAptly()
	repo := &LocalRepo{
		Name: "test_repo"}
	err := aptly.DeleteRepo(repo)
	if err != nil {
		t.Error(err)
	}
}
