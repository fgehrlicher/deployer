package cmd

import (
	"github.com/urfave/cli"
	"gitlab.osram.info/osram/deployer/git"
	"gitlab.osram.info/osram/deployer/cli_util"
	"gitlab.osram.info/osram/deployer/ruleset"
	"errors"
)

func DeployCommand(cliContext *cli.Context) error {
	Init(cliContext)

	err := git.InitPackage()
	if err != nil {
		return errorWrap(err)
	}

	gitController, err := git.NewGitController()
	if err != nil {
		return errorWrap(err)
	}

	deploymentContext, err := ruleset.CreateContext(gitController)
	if err != nil {
		return errorWrap(err)
	}

	deploymentContext.CommandReference = ruleset.DeployCommand

	selectionInformation := cli_util.NewTextBox()
	selectionInformation.Elements[ruleset.GetSelectionStepPosition(ruleset.TargetSelectionStep)] =
		ruleset.GetSelectionResultBoxes(ruleset.TargetSelectionStep)

	selectionInformation.Elements[ruleset.GetSelectionStepPosition(ruleset.CommitSelectionStep)] =
		ruleset.GetSelectionResultBoxes(ruleset.CommitSelectionStep)

	selectionInformation.Elements[ruleset.GetSelectionStepPosition(ruleset.DeploymentTypeSelection)] =
		ruleset.GetSelectionResultBoxes(ruleset.DeploymentTypeSelection)
	selectionInformation.RenderBox()


	var (
		predefinedTarget = cliContext.String("stage")
		selectedTarget   cli_util.Selectable
	)

	if ruleset.IsValidTarget(predefinedTarget) {
		predefinedTargetSelectable := cli_util.Selectable(ruleset.GetTarget(predefinedTarget))
		selectedTarget = predefinedTargetSelectable
	} else {
		var selectableTargets []cli_util.Selectable
		for i := 0; i < len(ruleset.Targets); i++ {
			selectableTargets = append(selectableTargets, cli_util.Selectable(ruleset.Targets[i]))
		}
		selectedTarget, err = cli_util.NewSingleSelect(selectableTargets, "Select Target:").RenderSelect()
		if err != nil {
			return errorWrap(err)
		}
	}

	showableSelectedTarget := selectedTarget.(cli_util.Showable)
	selectionInformation.Elements[ruleset.GetSelectionStepPosition(ruleset.TargetSelectionStep)].Item = showableSelectedTarget
	selectionInformation.ReRenderBox()

	return ForwardToSpecificDeployment(
		selectionInformation,
		cliContext,
		gitController,
		deploymentContext,
	)

}

func ForwardToSpecificDeployment(
	selectionInformation *cli_util.TextBox,
	cliContext *cli.Context,
	controller *git.Controller,
	deploymentContext *ruleset.Context,
) error {
	switch selectionInformation.Elements[ruleset.GetSelectionStepPosition(ruleset.TargetSelectionStep)].Item.GetId() {
	case ruleset.DevTarget:
		return DeployDevCommand(
			selectionInformation,
			cliContext,
			controller,
			deploymentContext,
		)
	case ruleset.QaTarget:
		return DeployQaCommand(
			selectionInformation,
			cliContext,
			controller,
			deploymentContext,
		)
	case ruleset.ProdTarget:
		return DeployProdCommand(
			selectionInformation,
			cliContext,
			controller,
			deploymentContext,
		)
	}
	return errorWrap(errors.New("invalid Target"))
}
