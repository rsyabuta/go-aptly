package aptly

import (
	"fmt"
	"testing"

	"github.com/kr/pretty"
)

func TestListSnapshots(t *testing.T) {
	client := testClient()
	snapshots, err := client.Snapshot.List()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%# v\n", pretty.Formatter(snapshots))
}

func TestCreateSnapshot(t *testing.T) {
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

	_, err = client.LocalRepo.AddFile("libboost-program-options-dev_1.49.0.1_i386.deb")
	if err != nil {
		t.Error(err)
	}
	_, err = client.Snapshot.CreateFromRepo(snapshot, repo)
	if err != nil {
		t.Error(err)
	}

}

/*
func TestDeleteRepo(t *testing.T) {
	client := testClient()
	repo := &LocalRepo{
		Name: "test_repo"}
	err := client.DeleteRepo(repo)
	if err != nil {
		t.Error(err)
	}
}
*/
