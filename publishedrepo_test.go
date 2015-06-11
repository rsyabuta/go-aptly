package aptly

import (
	"fmt"
	"testing"

	"github.com/kr/pretty"
)

func TestPublishFromSnapshot(t *testing.T) {
	client := testClient()
	repo := &LocalRepo{
		Name: "test_repo",
	}
	snapshot := &Snapshot{
		Name: "test_snap",
	}
	_, err := client.LocalRepo.Create(repo)
	if err != nil {
		t.Error(err)
	}

	_, err = client.File.UploadFile("libboost-program-options-dev_1.49.0.1_i386.deb")
	if err != nil {
		t.Error(err)
	}

	_, err = client.LocalRepo.AddFile("libboost-program-options-dev_1.49.0.1", repo)
	if err != nil {
		t.Error(err)
	}

	s, err := client.Snapshot.CreateFromRepo(snapshot, repo)
	if err != nil {
		t.Error(err)
	}

	pr, err := client.PublishedRepo.PublishFromSnapshot(s)
	if err != nil {
		t.Error(err)
	}

	publishedrepos, err := client.PublishedRepo.List()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%# v\n", pretty.Formatter(publishedrepos))

	err = client.PublishedRepo.Drop(pr)
	if err != nil {
		t.Error(err)
	}

	err = client.Snapshot.Delete(snapshot)
	if err != nil {
		t.Error(err)
	}
	err = client.LocalRepo.Delete(repo)
	if err != nil {
		t.Error(err)
	}

}
