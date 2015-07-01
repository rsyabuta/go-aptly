package aptly

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"
)

type SnapshotService struct {
	Client *Client
}

type Snapshot struct {
	Name        string    `json:"Name"`
	Description string    `json:"Description"`
	CreatedAt   time.Time `json:"CreatedAt"`
}

type SnapshotCollection struct {
	Snapshots []Snapshot
}

func (service *SnapshotService) List() (*SnapshotCollection, error) {
	resp, err := service.Client.Get("snapshots")
	if err != nil {
		return nil, err
	}
	var collection SnapshotCollection
	err = json.NewDecoder(resp.Body).Decode(&collection.Snapshots)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return &collection, err
}

func (service *SnapshotService) CreateFromRepo(snapshot *Snapshot, repo *LocalRepo) (*Snapshot, error) {
	if snapshot.Name == "" {
		return nil, errors.New("aptly: passed snapshot missing Name field")
	}
	reqBody, err := json.Marshal(snapshot)
	if err != nil {
		return nil, err
	}

	resp, err := service.Client.Post(fmt.Sprintf("repos/%s/snapshots", repo.Name), "application/json", nil, bytes.NewBuffer(reqBody))
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 400 {
		return nil, errors.New(fmt.Sprintf("aptly: %s", body))
	} else if resp.StatusCode == 404 {
		return nil, errors.New(fmt.Sprintf("aptly: %s", body))
	}

	var newSnapshot Snapshot
	err = json.Unmarshal(body, &newSnapshot)

	if err != nil {
		return nil, err
	}
	return &newSnapshot, nil
}

func (service *SnapshotService) Delete(snapshot *Snapshot) error {
	if snapshot.Name == "" {
		return errors.New("aptly: passed repo missing Name field")
	}
	resp, err := service.Client.Delete(fmt.Sprintf("snapshots/%s", snapshot.Name), nil, nil)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode == 400 {
		return errors.New(fmt.Sprintf("aptly: %s", body))
	} else if resp.StatusCode == 409 {
		return errors.New(fmt.Sprintf("aptly: %s", body))
	}
	return nil
}
