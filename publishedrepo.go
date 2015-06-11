package aptly

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

type PublishedRepoService struct {
	Client *Client
}

type PublishedRepo struct {
	Storage       string   `json:"Storage"`
	Prefix        string   `json:"Prefix"`
	Distribution  string   `json:"Distribution"`
	SourceKind    string   `json:"SourceKind"`
	Sources       []Source `json:"Sources"`
	Architectures []string `json:"Architectures"`
	Label         string   `json:"Label"`
	Origin        string   `json:"Origin"`
	Signing       Signing  `json:"Signing"`
}

type Source struct {
	Name      string `json:"Name"`
	Component string `json:"Component"`
}

type Signing struct {
	Skip           bool   `json:"Skip"`
	Batch          bool   `json:"Batch"`
	GpgKey         string `json:"GpgKey"`
	Keyring        string `json:"Keyring"`
	SecretKeyring  string `json:"SecretKeyring"`
	Passphrase     string `json:"Passphrase"`
	PassphraseFile string `json:"PassphraseFile"`
}

type PublishedRepoCollection struct {
	PublishedRepos []PublishedRepo
}

func (service *PublishedRepoService) List() (*PublishedRepoCollection, error) {
	resp, err := service.Client.Get("publish")
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	var collection PublishedRepoCollection
	err = json.NewDecoder(resp.Body).Decode(&collection.PublishedRepos)
	if err != nil {
		return nil, err
	}
	return &collection, err
}

func (service *PublishedRepoService) Publish(publishedrepo *PublishedRepo) (*PublishedRepo, error) {
	reqBody, err := json.Marshal(publishedrepo)
	if err != nil {
		return nil, err
	}

	resp, err := service.Client.Post("publish", "application/json", nil, bytes.NewBuffer(reqBody))
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
	} else if resp.StatusCode == 500 {
		return nil, errors.New(fmt.Sprintf("aptly: %s", body))
	}

	var pr PublishedRepo
	err = json.Unmarshal(body, &pr)

	if err != nil {
		return nil, err
	}
	return &pr, nil
}

func (service *PublishedRepoService) PublishFromLocalRepo(localrepo *LocalRepo) (*PublishedRepo, error) {
	if localrepo.Name == "" {
		return nil, errors.New("aptly: passed repo missing Name field")
	}
	source := service.BuildSourceFromLocalRepo(localrepo, "main")
	publishedrepo := &PublishedRepo{
		SourceKind:   "local",
		Distribution: "trusty",
		Sources:      source,
		Signing: Signing{
			Skip: true,
		},
	}
	pr, err := service.Publish(publishedrepo)
	if err != nil {
		return nil, err
	}
	return pr, nil
}

func (service *PublishedRepoService) PublishFromSnapshot(snapshot *Snapshot) (*PublishedRepo, error) {
	if snapshot.Name == "" {
		return nil, errors.New("aptly: passed repo missing Name field")
	}
	source := service.BuildSourceFromSnapshot(snapshot, "main")
	publishedrepo := &PublishedRepo{
		SourceKind:   "snapshot",
		Distribution: "trusty",
		Sources:      source,
		Signing: Signing{
			Skip: true,
		},
	}
	pr, err := service.Publish(publishedrepo)
	if err != nil {
		return nil, err
	}
	return pr, nil
}

func (service *PublishedRepoService) BuildSourceFromSnapshot(snapshot *Snapshot, component string) []Source {
	s := Source{
		Name:      snapshot.Name,
		Component: component}
	sl := []Source{s}
	return sl
}

func (service *PublishedRepoService) BuildSourceFromLocalRepo(localrepo *LocalRepo, component string) []Source {
	s := Source{
		Name:      localrepo.Name,
		Component: component}
	sl := []Source{s}
	return sl
}

func (service *PublishedRepoService) Drop(publishedrepo *PublishedRepo) error {
	resp, err := service.Client.Delete(fmt.Sprintf("publish/%s/%s", publishedrepo.Prefix, publishedrepo.Distribution), nil)
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
