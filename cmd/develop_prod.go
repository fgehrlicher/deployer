package cmd

import (
	"gitlab.osram.info/osram/deployer/cli_util"
	"gitlab.osram.info/osram/deployer/git"
	"github.com/urfave/cli"
	"gitlab.osram.info/osram/deployer/ruleset"
	"github.com/hashicorp/go-version"
	"fmt"
	"os"
	"gitlab.osram.info/osram/deployer/semantic_version"
)

func DeployProdCommand(
	headerBox *cli_util.TextBox,
	cliContext *cli.Context,
	gitController *git.Controller,
	deploymentContext *ruleset.Context,
) error {

	deploymentContext.CommandReference = ruleset.DeployProdCommand

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
	prodTags, err := gitController.GetProductionTags()
	if err != nil {
		return errorWrap(err)
	}

	qaTags, err := gitController.GetTagsByMetadata(ruleset.QaTagMetadata)
	if err != nil {
		return errorWrap(err)
	}

	if len(qaTags) == 0 {
		headerBox.Elements[ruleset.GetSelectionStepPosition(ruleset.TagSelectionStep)].Item = ruleset.NotAvailablePlaceholder
		headerBox.SetSuffix(
			"No " + cli_util.Bold("Quality Assurance Server") + " Deployments yet.\n" +
				" Its not possible to deploy to the " + cli_util.Bold("Production Server") +
				"\n without existing " + cli_util.Bold("Quality Assurance Server") + " Deployments.\n",
		)
		os.Exit(0)
	}

	availableTags, err := semantic_version.GetDiff(qaTags, prodTags, true)
	if err != nil {
		return errorWrap(err)
	}

	if len(availableTags) == 0 {
		headerBox.Elements[ruleset.GetSelectionStepPosition(ruleset.TagSelectionStep)].Item = ruleset.NotAvailablePlaceholder
		headerBox.SetSuffix(
			"No available " + cli_util.Bold("Quality Assurance Server") + " Deployments to select.\n" +
				" Its not possible to deploy to the " + cli_util.Bold("Production Server") +
				"\n without undeployed " + cli_util.Bold("Quality Assurance Server") + " Deployments \n",
		)
		os.Exit(0)
	}

	semantic_version.StripIterationsOfSameVersion(availableTags)

	//@TODO move into ruleset
	var availabeTagsSlice []cli_util.Selectable
	for i := 0; i < len(availableTags); i++ {
		tag := availableTags[i]
		availabeTagsSlice = append(
			availabeTagsSlice,
			semantic_version.Version{
				UnderlyingVersion: tag,
				TextPrefix: "-> " +
					fmt.Sprintf(
						"%v.%v.%v",
						tag.Segments()[0],
						tag.Segments()[1],
						tag.Segments()[2],
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
			"%v.%v.%v",
			selectedVersion.Segments()[0],
			selectedVersion.Segments()[1],
			selectedVersion.Segments()[2],
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
