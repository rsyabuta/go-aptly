package aptly

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
)

type LocalRepo struct {
	Name                string `json:"Name"`
	Comment             string `json:"Comment"`
	DefaultDistribution string `json:"DefaultDistribution"`
	DefaultComponent    string `json:"DefaultComponent"`
}

func (aptly *Aptly) CreateRepo(repo *LocalRepo) (*LocalRepo, error) {
	if repo.Name == "" {
		return nil, errors.New("aptly: passed repo missing Name field")
	}
	reqBody, err := json.Marshal(repo)
	if err != nil {
		return nil, err
	}

	resp, err := aptly.Post("repos", "application/json", nil, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 400 {
		return nil, errors.New("aptly: repo already exists")
	}

	var newRepo LocalRepo
	err = json.NewDecoder(resp.Body).Decode(&newRepo)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return &newRepo, err
}

func (aptly *Aptly) DeleteRepo(repo *LocalRepo) error {
	return nil
}

func (aptly *Aptly) ListRepos() ([]byte, error) {
	resp, err := aptly.Get("repos")
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return body, nil
}
