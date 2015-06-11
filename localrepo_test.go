package aptly

import (
	"fmt"
	"testing"

	"github.com/kr/pretty"
)

func TestCreateRepo(t *testing.T) {
	client := testClient()
	repo := &LocalRepo{
		Name: "test_repo"}
	createdRepo, err := client.LocalRepo.Create(repo)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%# v\n", pretty.Formatter(createdRepo))
}

func TestListRepos(t *testing.T) {
	client := testClient()
	repos, err := client.LocalRepo.List()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%# v\n", pretty.Formatter(repos))
}

func TestDeleteRepo(t *testing.T) {
	client := testClient()
	repo := &LocalRepo{
		Name: "test_repo"}
	err := client.LocalRepo.Delete(repo)
	if err != nil {
		t.Error(err)
	}
}
