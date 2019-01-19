package git

import (
	"errors"
	"github.com/hashicorp/go-version"
	"gitlab.osram.info/osram/deployer/cli_util"
	"gitlab.osram.info/osram/deployer/config"
	"gitlab.osram.info/osram/deployer/models"
	"gopkg.in/src-d/go-git.v4"
	gitConfig "gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"io"
	"sort"
	"strings"
)

var (
	InvalidCloneMethod = errors.New("invalid clone method")
	InvalidUsername    = errors.New("invalid username")
	InvalidPassword    = errors.New("invalid password")
	Authenticator      transport.AuthMethod
	Config             config.Config
)

type Controller struct {
	repository    *git.Repository
	memoryStorage *memory.Storage
}

func InitPackage() error {
	var err error
	Config = config.LoadConfig()
	Authenticator, err = GetAuth(Config)
	return err
}

func NewGitController() (*Controller, error) {
	var (
		err        error
		controller = &Controller{}
	)

	cloneOptions := &git.CloneOptions{
		Auth:          Authenticator,
		URL:           Config.RemoteUrl,
		ReferenceName: "refs/heads/develop",
	}

	err = controller.CloneIntoMemory(cloneOptions)
	if err != nil {
		return nil, err
	}

	return controller, nil
}

func (this *Controller) CloneIntoMemory(options *git.CloneOptions) error {
	memoryStorage := memory.NewStorage()
	repository, err := git.Clone(
		memoryStorage,
		nil,
		options,
	)
	if err != nil {
		return err
	}
	this.repository = repository
	return nil
}

func (this *Controller) GetCommitsUpToTag() ([]models.Commit, error) {
	var commits []models.Commit

	logOptions := &git.LogOptions{
		Order: git.LogOrderCommitterTime,
	}

	commitIter, err := this.repository.Log(logOptions)
	if err != nil {
		return nil, err
	}

	referenceIter, err := this.repository.Tags()
	if err != nil {
		return nil, err
	}

	var hashes []string
	referenceIter.ForEach(func(reference *plumbing.Reference) error {
		hashes = append(hashes, reference.Hash().String())
		return nil
	})

	for {
		commit, err := commitIter.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		commitHash := commit.Hash.String()
		found := false

		for _, hash := range hashes {
			if hash == commitHash {
				found = true
				break
			}
		}

		if found {
			break
		}

		commits = append(commits, models.Commit{
			Message: cli_util.SanitizeString(commit.Message),
			Author:  commit.Author.Email,
			Hash:    commitHash,
		})

	}

	return commits, nil
}

func (this *Controller) GetTagsByMetadata(metadata string) (version.Collection, error) {
	var versions version.Collection

	tags, err := this.repository.Tags()
	if err != nil {
		return nil, err
	}

	tags.ForEach(
		func(reference *plumbing.Reference) error {

			realTagName := strings.Replace(string(reference.Name()), "refs/tags/", "", -1)
			versionTag, err := version.NewVersion(realTagName)
			if err != nil {
				return err
			}
			if versionTag.Metadata() != metadata {
				return nil
			}
			versions = append(versions, versionTag)
			return nil
		},
	)

	if len(versions) > 0 {
		sort.Sort(sort.Reverse(versions))
	}

	return versions, nil
}

func (this *Controller) GetProductionTags() (version.Collection, error) {
	var versions version.Collection

	tags, err := this.repository.Tags()
	if err != nil {
		return nil, err
	}

	tags.ForEach(
		func(reference *plumbing.Reference) error {

			realTagName := strings.Replace(string(reference.Name()), "refs/tags/", "", -1)
			versionTag, err := version.NewVersion(realTagName)
			if err != nil {
				return err
			}
			if versionTag.Metadata() != "" || versionTag.Prerelease() != "" {
				return nil
			}
			versions = append(versions, versionTag)
			return nil
		},
	)

	if len(versions) > 0 {
		sort.Sort(sort.Reverse(versions))
	}

	return versions, nil
}

func (this *Controller) GetTagReference(tagName string) (*plumbing.Reference, error) {
	return this.repository.Tag(tagName)
}

func (this *Controller) AddTag(tagName string, referenceHash string) error {
	compiledHash := plumbing.NewHash(referenceHash)
	reference := plumbing.ReferenceName("refs/tags/" + tagName)
	hashReference := plumbing.NewHashReference(reference, compiledHash)
	return this.repository.Storer.SetReference(hashReference)
}

func (this *Controller) PushTags() error {
	return this.repository.Push(
		&git.PushOptions{
			Auth: Authenticator,
			RefSpecs: []gitConfig.RefSpec{
				gitConfig.RefSpec("refs/tags/*:refs/tags/*"),
			},
		},
	)
}
