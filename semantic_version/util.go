package semantic_version

import (
	"errors"
	"fmt"
	"sort"

	"github.com/hashicorp/go-version"

	"github.com/fgehrlicher/deployer/cli_util"
)

type Version struct {
	UnderlyingVersion *version.Version
	TextPrefix        string
}

func (this Version) GetId() string {
	return this.UnderlyingVersion.String()
}

func (this Version) GetFormattedText() string {
	return this.UnderlyingVersion.String() + " " + this.TextPrefix
}

func (this Version) GetDisplayText() string {
	return this.UnderlyingVersion.String() + " " + this.TextPrefix
}

func ChangeMetadata(oldVersion *version.Version, newMetadata string) (*version.Version, error) {
	oldVersionSegments := oldVersion.Segments()
	return version.NewVersion(fmt.Sprintf(
		"%v.%v.%v-%v+%v",
		oldVersionSegments[0],
		oldVersionSegments[1],
		oldVersionSegments[2],
		oldVersion.Prerelease(),
		newMetadata,
	))
}

func GetDiff(master version.Collection, comparison version.Collection, ignorePrerelease bool) (version.Collection, error) {
	if len(master) == 0 {
		return nil, errors.New("the master collection cant be empty")
	}
	sort.Sort(sort.Reverse(master))

	if len(comparison) == 0 {
		return master, nil
	}

	var latestComparison string
	if ignorePrerelease {
		latestComparison = GetSimpleVersion(comparison[0])
	} else {
		latestComparison = GetVersionStringWithoutMetaData(*comparison[0])
	}

	var (
		found      bool
		foundIndex int
	)

	for index, value := range master {
		var currentVersion string
		if ignorePrerelease {
			currentVersion = GetSimpleVersion(value)
		} else {
			currentVersion = GetVersionStringWithoutMetaData(*value)
		}
		if currentVersion == latestComparison {
			found = true
			foundIndex = index
			break
		}
	}

	if !found {
		return nil, errors.New("latest element in comparison is not in master")
	}

	if foundIndex == 0 {
		return version.Collection{}, nil
	}

	return master[0:foundIndex], nil
}

func ConvertCollectionToSliceOfSelectables(collection version.Collection) []cli_util.Selectable {
	result := make([]cli_util.Selectable, len(collection))
	for index, value := range collection {
		result[index] = Version{UnderlyingVersion: value}
	}
	return result
}

func GetVersionStringWithoutMetaData(version version.Version) string {
	versionSegments := version.Segments()
	return fmt.Sprintf(
		"%v.%v.%v-%v",
		versionSegments[0],
		versionSegments[1],
		versionSegments[2],
		version.Prerelease(),
	)
}

func GetSimpleVersion(version *version.Version) string {
	if version == nil {
		return ""
	}

	versionSegments := version.Segments()
	return fmt.Sprintf(
		"%v.%v.%v",
		versionSegments[0],
		versionSegments[1],
		versionSegments[2],
	)
}

func StripIterationsOfSameVersion(collection version.Collection) version.Collection {
	var returnCollection version.Collection
	if len(collection) == 0 {
		return returnCollection
	}

	sort.Sort(sort.Reverse(collection))
	var foundVersions []string

	for _, value := range collection {
		rawVersion := GetSimpleVersion(value)
		if !InSlice(rawVersion, foundVersions) {
			returnCollection = append(returnCollection, value)
			foundVersions = append(foundVersions, rawVersion)
		}
	}

	return returnCollection
}

func InSlice(needle string, haystack []string) bool {
	for _, value := range haystack {
		if value == needle {
			return true
		}
	}
	return false
}
