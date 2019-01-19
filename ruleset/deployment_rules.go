package ruleset

import (
	"gitlab.osram.info/osram/deployer/cli_util"
	"fmt"
	"github.com/hashicorp/go-version"
	"gitlab.osram.info/osram/deployer/models"
	"gitlab.osram.info/osram/deployer/semantic_version"
	"strconv"
)

func IsValidTarget(targetName string) bool {
	for _, target := range Targets {
		if target.GetId() == targetName {
			return true
		}
	}
	return false
}

func IsValidDeploymentType(deploymentType string) bool {
	for target := range DeploymentTypes {
		if target == deploymentType {
			return true
		}
	}
	return false
}

func GetSelectionStepPosition(stepName string) int {
	return SelectionStepPositions[stepName]
}

func GetSelectionResultBoxes(stepName string) *cli_util.TextBoxElement {
	return SelectionResultBoxes[stepName]
}

func GetTarget(id string) *models.Target {
	for _, target := range Targets {
		if target.GetId() == id {
			return &target
		}
	}
	return nil
}

func NewVersion(oldVersion version.Version, incrementType string, newMetadata string) (*version.Version, error) {
	var (
		patchVersion int
		minorVersion int
		majorVersion int
		iteration    int
	)

	oldPre := oldVersion.Prerelease()
	if oldPre != "" {
		intVersion, err := strconv.Atoi(oldPre)
		if err != nil {
			return nil, err
		}
		iteration = intVersion
	}

	switch incrementType {
	case IncrementVersion:
		patchVersion = oldVersion.Segments()[2]
		minorVersion = oldVersion.Segments()[1]
		majorVersion = oldVersion.Segments()[0]
		iteration = iteration + 1
	case PatchVersion:
		patchVersion = oldVersion.Segments()[2] + 1
		minorVersion = oldVersion.Segments()[1]
		majorVersion = oldVersion.Segments()[0]
		iteration = 1
	case MinorVersion:
		patchVersion = 0
		minorVersion = oldVersion.Segments()[1] + 1
		majorVersion = oldVersion.Segments()[0]
		iteration = 1
	case MajorVersion:
		patchVersion = 0
		minorVersion = 0
		majorVersion = oldVersion.Segments()[0] + 1
		iteration = 1
	}

	if newMetadata == "" {
		newMetadata = oldVersion.Metadata()
	}

	return version.NewVersion(
		fmt.Sprintf("%v.%v.%v-%v+%v", majorVersion, minorVersion, patchVersion, iteration, newMetadata),
	)
}

func (this *Context) GetCurrentVersion() string {
	var (
		currentVersion                 *version.Version
		simpleCurrentVersion           string
		currentProductionVersion       = this.GetLatestProductionVersion()
		simpleCurrentProductionVersion string
		headerString                   string
	)

	simpleCurrentProductionVersion = semantic_version.GetSimpleVersion(currentProductionVersion)

	switch this.CommandReference {
	case DeployDevCommand:
		currentVersion = this.GetLatestDevelopmentVersion(DevTagMetadata)
	case DeployQaCommand:
		currentVersion = this.GetLatestDevelopmentVersion(QaTagMetadata)
	case DeployProdCommand:
		return fmt.Sprintf(
			"Current Version on Target Server: %v",
			simpleCurrentProductionVersion,
		)
	default:
		return headerString
	}

	simpleCurrentVersion = semantic_version.GetSimpleVersion(currentVersion)

	if simpleCurrentVersion == "" {
		headerString = ""
	} else if simpleCurrentProductionVersion == simpleCurrentVersion {
		headerString = fmt.Sprintf(
			"Current Version: %v (already deployed on production)",
			simpleCurrentProductionVersion,
		)
	} else {
		headerString = fmt.Sprintf(
			"Current Version: %v",
			currentVersion.String(),
		)
	}

	return headerString
}

func (this *Context) GetDeploymentTypes() []*models.DeployType {
	var (
		currentProductionVersion = this.GetLatestProductionVersion()
		currentVersion           = this.GetLatestDevelopmentVersion(DevTagMetadata)
		availableTypes           []*models.DeployType
		deploymentTypes          = DeploymentTypes
	)

	for _, deployType := range deploymentTypes {
		if deployType.GetId() == IncrementVersion {
			continue
		}
		var newVersion *version.Version
		if currentVersion != nil {
			newVersion, _ = NewVersion(*currentVersion, deployType.GetId(), "")
		} else {
			newVersion, _ = NewVersion(*PlaceholderVersion, deployType.GetId(), DevTagMetadata)
		}
		availableTypes = append(
			availableTypes,
			models.NewDeployType(
				deployType.GetId(),
				deployType.GetFormattedText()+" -> "+newVersion.String(),
			),
		)
	}

	if currentVersion != nil &&
		semantic_version.GetSimpleVersion(currentProductionVersion) != semantic_version.GetSimpleVersion(currentVersion) {
		incrementDeployType := deploymentTypes[IncrementVersion]
		newVersion, _ := NewVersion(*currentVersion, IncrementVersion, "")

		availableTypes = append(
			availableTypes,
			models.NewDeployType(
				IncrementVersion,
				incrementDeployType.GetFormattedText()+"   -> "+newVersion.String(),
			),
		)
	}
	return availableTypes
}
