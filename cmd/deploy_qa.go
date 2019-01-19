package cmd

import (
	"gitlab.osram.info/osram/deployer/cli_util"
	"github.com/urfave/cli"
	"gitlab.osram.info/osram/deployer/git"
	"gitlab.osram.info/osram/deployer/ruleset"
	"os"
	"gitlab.osram.info/osram/deployer/semantic_version"
	"fmt"
	"github.com/hashicorp/go-version"
)

func DeployQaCommand(
	headerBox *cli_util.TextBox,
	cliContext *cli.Context,
	gitController *git.Controller,
	deploymentContext *ruleset.Context,
) error {

	deploymentContext.CommandReference = ruleset.DeployQaCommand

	headerBox.Clear()
	delete(headerBox.Elements, ruleset.GetSelectionStepPosition(ruleset.CommitSelectionStep))
	delete(headerBox.Elements, ruleset.GetSelectionStepPosition(ruleset.DeploymentTypeSelection))
	headerBox.Elements[ruleset.GetSelectionStepPosition(ruleset.TagSelectionStep)] =
		ruleset.GetSelectionResultBoxes(ruleset.TagSelectionStep)
	headerBox.RenderBox()

	headerBox.SetSuffix(
		deploymentContext.GetCurrentVersion(),
	)

	// Tag Selection
	qaTags, err := gitController.GetTagsByMetadata(ruleset.QaTagMetadata)
	if err != nil {
		return errorWrap(err)
	}

	devTags, err := gitController.GetTagsByMetadata(ruleset.DevTagMetadata)
	if err != nil {
		return errorWrap(err)
	}

	if len(devTags) == 0 {
		headerBox.Elements[ruleset.GetSelectionStepPosition(ruleset.TagSelectionStep)].Item = ruleset.NotAvailablePlaceholder
		headerBox.SetSuffix(
			"No " + cli_util.Bold("Development Server") + " Deployments yet.\n" +
				" Its not possible to deploy to the " + cli_util.Bold("Quality Assurance Server") +
				"\n without existing " + cli_util.Bold("Development Server") + " Deployments.\n",
		)
		os.Exit(0)
	}

	availabeTags, err := semantic_version.GetDiff(devTags, qaTags, false)
	if err != nil {
		return errorWrap(err)
	}

	if len(availabeTags) == 0 {
		headerBox.Elements[ruleset.GetSelectionStepPosition(ruleset.TagSelectionStep)].Item = ruleset.NotAvailablePlaceholder
		headerBox.SetSuffix(
			"No available " + cli_util.Bold("Development Server") + " Deployments to select.\n" +
				" Its not possible to deploy to the " + cli_util.Bold("Quality Assurance Server") +
				"\n without undeployed " + cli_util.Bold("Development Server") + " Deployments \n",
		)
		os.Exit(0)
	}

	//@TODO move into ruleset
	var availabeTagsSlice []cli_util.Selectable
	for i := 0; i < len(availabeTags); i++ {
		tag := availabeTags[i]
		availabeTagsSlice = append(
			availabeTagsSlice,
			semantic_version.Version{
				UnderlyingVersion: tag,
				TextPrefix: "-> " +
					fmt.Sprintf(
						"%v.%v.%v-%v+%v",
						tag.Segments()[0],
						tag.Segments()[1],
						tag.Segments()[2],
						tag.Prerelease(),
						ruleset.QaTagMetadata,
					),
			},
		)
	}

	var selectedTag cli_util.Selectable
	if cliContext.Bool("latest") {
		selectedTag = availabeTagsSlice[0]
	} else {
		selectedTag, err = cli_util.NewSingleSelect(availabeTagsSlice, "Select Version to deploy: ").RenderSelect()
		if err != nil {
			return errorWrap(err)
		}
	}

	headerBox.Elements[ruleset.GetSelectionStepPosition(ruleset.TagSelectionStep)].Item = selectedTag.(cli_util.Showable)
	headerBox.ReRenderBox()

	selectedVersion, err := version.NewVersion(selectedTag.GetId())
	if err != nil {
		return errorWrap(err)
	}

	//@TODO move into ruleset
	targetVersion, err := version.NewVersion(
		fmt.Sprintf(
			"%v.%v.%v-%v+%v",
			selectedVersion.Segments()[0],
			selectedVersion.Segments()[1],
			selectedVersion.Segments()[2],
			selectedVersion.Prerelease(),
			ruleset.QaTagMetadata,
		),
	)
	if err != nil {
		return errorWrap(err)
	}

	headerBox.SetSuffix(
		fmt.Sprintf("Will result in version: %v", targetVersion.String()),
	)

	confirmation := cli_util.NewConfirmation(" Confirm deployment? (y/n):")
	err = confirmation.Render()
	if err != nil {
		return errorWrap(err)
	}

	tagReference, err := gitController.GetTagReference(selectedTag.GetId())
	if err != nil {
		return errorWrap(err)
	}

	err = gitController.AddTag(targetVersion.String(), tagReference.Hash().String())
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
