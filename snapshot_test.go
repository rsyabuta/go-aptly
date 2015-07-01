package aptly

import (
	"fmt"
	"testing"

	"github.com/kr/pretty"
)

func TestDeleteAllSnapshots(t *testing.T) {
	client := testClient()
	snaps, err := client.Snapshot.List()
	if err != nil {
		t.Error(err)
	}
	for _, s := range snaps.Snapshots {
		err = client.Snapshot.Delete(&s)
		if err != nil {
			t.Error(err)
		}
	}
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

	fr, err := client.LocalRepo.AddFile("libboost-program-options-dev_1.49.0.1", repo)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%# v\n", pretty.Formatter(fr))

	ns, err := client.Snapshot.CreateFromRepo(snapshot, repo)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%# v\n", pretty.Formatter(ns))

}

func TestListSnapshots(t *testing.T) {
	client := testClient()
	snapshots, err := client.Snapshot.List()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%# v\n", pretty.Formatter(snapshots))
}

func TestDeleteSnapshot(t *testing.T) {
	client := testClient()
	repo := &LocalRepo{
		Name: "test_repo",
	}
	snapshot := &Snapshot{
		Name: "test_snap",
	}
	err := client.Snapshot.Delete(snapshot)
	if err != nil {
		t.Error(err)
	}
	err = client.LocalRepo.Delete(repo)
	if err != nil {
		t.Error(err)
	}

}
