package ruleset

import (
	"github.com/hashicorp/go-version"
	"gitlab.osram.info/osram/deployer/git"
	"sort"
)

type Context struct {
	GitTags struct {
		ProductionTags  version.Collection
		DevelopmentTags map[string]version.Collection
	}
	CommandReference int
}

func CreateContext(controller *git.Controller) (*Context, error) {
	context := Context{}

	productionTags, err := controller.GetProductionTags()
	if err != nil {
		return nil, err
	}
	sort.Sort(sort.Reverse(productionTags))
	context.GitTags.ProductionTags = productionTags
	context.GitTags.DevelopmentTags = make(map[string]version.Collection)


	for _, metadata := range TagMetadata {
		metadataTags, err := controller.GetTagsByMetadata(metadata)
		if err != nil {
			return nil, err
		}
		sort.Sort(sort.Reverse(metadataTags))
		context.GitTags.DevelopmentTags[metadata] = metadataTags
	}

	return &context, nil
}

func (this *Context) GetLatestDevelopmentVersion(metadata string) *version.Version {
	elements, ok := this.GitTags.DevelopmentTags[metadata]
	if !ok || len(elements) == 0 {
		return nil
	}
	return elements[0]
}

func (this *Context) GetLatestProductionVersion() *version.Version {
	if len(this.GitTags.ProductionTags) > 0 {
		return this.GitTags.ProductionTags[0]
	}
	return nil
}
