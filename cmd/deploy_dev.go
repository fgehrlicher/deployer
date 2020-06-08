package cmd

import (
	"fmt"
	"os"

	"github.com/urfave/cli"

	"github.com/fgehrlicher/deployer/cli_util"
	"github.com/fgehrlicher/deployer/git"
	"github.com/fgehrlicher/deployer/ruleset"
)

func DeployDevCommand(
	headerBox *cli_util.TextBox,
	cliContext *cli.Context,
	gitController *git.Controller,
	deploymentContext *ruleset.Context,
) error {

	deploymentContext.CommandReference = ruleset.DeployDevCommand

	headerBox.SetSuffix(
		deploymentContext.GetCurrentVersion(),
	)

	// Commit Selection
	commits, err := gitController.GetCommitsUpToTag()
	if err != nil {
		return errorWrap(err)
	}

	if len(commits) == 0 {
		headerBox.Elements[ruleset.GetSelectionStepPosition(ruleset.CommitSelectionStep)].Item = ruleset.NotAvailablePlaceholder
		headerBox.Elements[ruleset.GetSelectionStepPosition(ruleset.DeploymentTypeSelection)].Item = ruleset.NotAvailablePlaceholder
		headerBox.SetSuffix("No available Commits to deploy.\n" +
			" It´s not possible to deploy to the " + cli_util.Bold("Development Server") +
			"\n without Commits that aren´t already included in previous \n " +
			cli_util.Bold("Development Server") + " Deployments.\n",
		)
		os.Exit(0)
	}

	var selectedCommit cli_util.Selectable
	if cliContext.Bool("latest") {
		selectedCommit = commits[0]
	} else {
		var selectableCommits []cli_util.Selectable
		for i := 0; i < len(commits); i++ {
			selectableCommits = append(selectableCommits, cli_util.Selectable(commits[i]))
		}
		selectedCommit, err = cli_util.NewSingleSelect(selectableCommits, "Select Commit:").RenderSelect()
		if err != nil {
			return errorWrap(err)
		}
	}

	headerBox.Elements[ruleset.GetSelectionStepPosition(ruleset.CommitSelectionStep)].Item = selectedCommit.(cli_util.Showable)
	headerBox.ReRenderBox()

	// Deployment Type Selection
	var (
		selectedDeployType    cli_util.Selectable
		preselectedDeployType = cliContext.String("type")
		availableDeployTypes  = deploymentContext.GetDeploymentTypes()
		validPreselect        bool
	)

	if preselectedDeployType != "" {
		for _, item := range availableDeployTypes {
			if item.GetId() == preselectedDeployType {
				validPreselect = true
				selectedDeployType = item
			}
		}
	}

	if !validPreselect {
		var selectableDeployTypes []cli_util.Selectable
		for i := 0; i < len(availableDeployTypes); i++ {
			selectableDeployTypes = append(selectableDeployTypes, cli_util.Selectable(availableDeployTypes[i]))
		}

		selectedDeployType, err = cli_util.NewSingleSelect(selectableDeployTypes, "Select Deployment Type:").RenderSelect()
		if err != nil {
			return errorWrap(err)
		}
	}

	headerBox.Elements[ruleset.GetSelectionStepPosition(ruleset.DeploymentTypeSelection)].Item = selectedDeployType.(cli_util.Showable)
	headerBox.ReRenderBox()

	currentVersion := deploymentContext.GetLatestDevelopmentVersion(ruleset.DevTagMetadata)
	if currentVersion == nil {
		currentVersion = ruleset.PlaceholderVersion
	}

	newVersion, err := ruleset.NewVersion(*currentVersion, selectedDeployType.GetId(), ruleset.DevTagMetadata)
	if err != nil {
		return errorWrap(err)
	}

	headerBox.SetSuffix(
		fmt.Sprintf("Will result in version: %v", newVersion.String()),
	)

	confirmation := cli_util.NewConfirmation(" Confirm deployment? (y/n):")
	err = confirmation.Render()
	if err != nil {
		return errorWrap(err)
	}

	err = gitController.AddTag(newVersion.String(), selectedCommit.GetId())
	if err != nil {
		return errorWrap(err)
	}

	err = gitController.PushTags()
	if err != nil {
		return errorWrap(err)
	}

	confirmation.Cleanup()
	headerBox.SetSuffix("Marked commit for deployment.\n")

	return nil
}
