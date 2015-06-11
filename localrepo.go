package aptly

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

type LocalRepoService struct {
	Client *Client
}

type LocalRepo struct {
	Name                string `json:"Name"`
	Comment             string `json:"Comment"`
	DefaultDistribution string `json:"DefaultDistribution"`
	DefaultComponent    string `json:"DefaultComponent"`
}

type LocalRepoCollection struct {
	LocalRepos []LocalRepo
}

type FileResponse struct {
	FailedFiles []string   `json:"FailedFiles"`
	Report      FileReport `json:"Report"`
}

type FileReport struct {
	Warnings []string `json:"Warnings"`
	Added    []string `json:"Added"`
	Deleted  []string `json:"Deleted"`
}

func (service *LocalRepoService) Create(repo *LocalRepo) (*LocalRepo, error) {
	if repo.Name == "" {
		return nil, errors.New("aptly: passed repo missing Name field")
	}
	reqBody, err := json.Marshal(repo)
	if err != nil {
		return nil, err
	}

	resp, err := service.Client.Post("repos", "application/json", nil, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 400 {
		return nil, errors.New(fmt.Sprintf("aptly: repo %s already exists", repo.Name))
	}

	var newRepo LocalRepo
	err = json.NewDecoder(resp.Body).Decode(&newRepo)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return &newRepo, err
}

func (service *LocalRepoService) Delete(repo *LocalRepo) error {
	if repo.Name == "" {
		return errors.New("aptly: passed repo missing Name field")
	}
	resp, err := service.Client.Delete(fmt.Sprintf("repos/%s", repo.Name), nil)
	if err != nil {
		return err
	}

	if resp.StatusCode == 404 {
		return errors.New(fmt.Sprintf("aptly: repo %s doesn't exist", repo.Name))
	} else if resp.StatusCode == 409 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		return errors.New(fmt.Sprintf("%q", body))
	}
	return nil
}

func (service *LocalRepoService) AddFile(file string) (*FileReport, error) {
	return nil, nil
}

func (service *LocalRepoService) List() (*LocalRepoCollection, error) {
	resp, err := service.Client.Get("repos")
	if err != nil {
		return nil, err
	}
	var collection LocalRepoCollection
	err = json.NewDecoder(resp.Body).Decode(&collection.LocalRepos)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return &collection, err
}
