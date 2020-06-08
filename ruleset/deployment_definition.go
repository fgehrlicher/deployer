package ruleset

import (
	"github.com/hashicorp/go-version"

	"github.com/fgehrlicher/deployer/cli_util"
	"github.com/fgehrlicher/deployer/models"
)

const DevTarget = "dev"
const QaTarget = "qa"
const ProdTarget = "prod"

const IncrementVersion = "increment"
const PatchVersion = "patch"
const MinorVersion = "minor"
const MajorVersion = "major"

const DevTagMetadata = "dev"
const QaTagMetadata = "qa"

const TargetSelectionStep = "targetselection"
const CommitSelectionStep = "commitselection"
const TagSelectionStep = "tagselection"
const DeploymentTypeSelection = "deploymenttypeselection"

const DeployCommand = 1
const DeployDevCommand = 2
const DeployQaCommand = 3
const DeployProdCommand = 4

var (
	EmptyPlaceholder        = models.NewPlaceholder("     ")
	NotAvailablePlaceholder = models.NewPlaceholder("-")
	PlaceholderVersion, _   = version.NewVersion("0.0.0")
)

var TagMetadata = []string{
	DevTagMetadata,
	QaTagMetadata,
}

var Targets = []models.Target{
	*models.NewTarget(DevTarget, "Development Server"),
	*models.NewTarget(QaTarget, "Quality Assurance Server"),
	*models.NewTarget(ProdTarget, "Production Server"),
}

var DeploymentTypes = map[string]*models.DeployType{
	MajorVersion:     models.NewDeployType(MajorVersion, "Create new Major Version"),
	MinorVersion:     models.NewDeployType(MinorVersion, "Create new Minor Version"),
	PatchVersion:     models.NewDeployType(PatchVersion, "Create new Patch Version"),
	IncrementVersion: models.NewDeployType(IncrementVersion, "Deploy Current Version"),
}

var SelectionStepPositions = map[string]int{
	TargetSelectionStep:     0,
	CommitSelectionStep:     1,
	TagSelectionStep:        1,
	DeploymentTypeSelection: 2,
}

var SelectionResultBoxes = map[string]*cli_util.TextBoxElement{
	TargetSelectionStep:     cli_util.NewTextBoxElement("Selected Target", EmptyPlaceholder),
	CommitSelectionStep:     cli_util.NewTextBoxElement("Selected Commit", EmptyPlaceholder),
	TagSelectionStep:        cli_util.NewTextBoxElement("Selected Version", EmptyPlaceholder),
	DeploymentTypeSelection: cli_util.NewTextBoxElement("Selected Deployment Type", EmptyPlaceholder),
}
