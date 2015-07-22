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

type Package struct {
	Architecture  string `json:"Architecture"`
	Description   string `json:"Description"`
	Filename      string `json:"Filename"`
	FilesHash     string `json:"FilesHash"`
	Homepage      string `json:"Homepage"`
	InstalledSize string `json:"Installed-Size"`
	Key           string `json:"Key"`
	License       string `json:"License"`
	MD5Sum        string `json:"MD5Sum"`
	Maintainer    string `json:"Maintainer"`
	Package       string `json:"Package"`
	Priority      string `json:"Priority"`
	Recommends    string `json:"Recommends"`
	SHA1          string `json:"SHA1"`
	SHA256        string `json:"SHA256"`
	Section       string `json:"Section"`
	ShortKey      string `json:"ShortKey"`
	Size          string `json:"Size"`
	Vendor        string `json:"Vendor"`
	Version       string `json:"Version"`
}

func (pl *PackageCollection) Len() int {
	return len(pl.Packages)

}

func (pl *PackageCollection) Swap(a, b int) {
	pl.Packages[a], pl.Packages[b] = pl.Packages[b], pl.Packages[a]
}

func (pl *PackageCollection) Less(a, b int) bool {
	if pl.Packages[a].Package == pl.Packages[b].Package {
		return pl.Packages[a].Package < pl.Packages[b].Package
	}

	return pl.Packages[a].Version < pl.Packages[b].Version
}

type PackageCollection struct {
	Packages []Package
}

func (service *LocalRepoService) Get(name string) (*LocalRepo, error) {
	resp, err := service.Client.Get(fmt.Sprintf("repos/%s", name))
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	var repo LocalRepo
	err = json.NewDecoder(resp.Body).Decode(&repo)
	if err != nil {
		return nil, err
	}
	return &repo, err
}

func (service *LocalRepoService) Packages(repo *LocalRepo) (*PackageCollection, error) {
	params := map[string]string{
		"format": "details",
	}
	resp, err := service.Client.GetWithParams(fmt.Sprintf("repos/%s/packages", repo.Name), params)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	var pc PackageCollection
	err = json.NewDecoder(resp.Body).Decode(&pc.Packages)
	if err != nil {
		return nil, err
	}
	return &pc, err
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
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 400 {
		return nil, errors.New(fmt.Sprintf("aptly: repo %s already exists", repo.Name))
	}

	var newRepo LocalRepo
	err = json.NewDecoder(resp.Body).Decode(&newRepo)
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
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode == 404 {
		return errors.New(fmt.Sprintf("aptly: %s", body))
	} else if resp.StatusCode == 409 {
		return errors.New(fmt.Sprintf("aptly: %s", body))
	}
	return nil
}

func (service *LocalRepoService) AddFile(file string, repo *LocalRepo) (*FileResponse, error) {
	if repo.Name == "" {
		return nil, errors.New("aptly: passed repo missing Name field")
	}

	resp, err := service.Client.Post(fmt.Sprintf("repos/%s/file/%s", repo.Name, file), "application/json", nil, nil)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 404 {
		return nil, errors.New(fmt.Sprintf("aptly: repo %s doesn't exist", repo.Name))
	}
	var result FileResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, err
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
