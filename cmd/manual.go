package cmd

import (
	"fmt"

	"github.com/urfave/cli"
)

func ManualCommand(c *cli.Context) error {
	Init(c)

	var (
		binaryName = c.App.Name
	)

	fmt.Println(`DEPLOYMEMT MANUAL:

  ENVIRONMENT VARIABLE CONFIGURATION:
    The deployment tool is configured by environment variables. The current 
    configuration can be validated by the "validate" command.
    e.g "` + binaryName + ` validate"

    GENERAL CONFIGURATION:
      These environment variables are always required.
      +======================+==========+============================================+
      | Environment Variable | Required |                Description                 |
      +======================+==========+============================================+
      | CLONE_METHOD         |    X     | The git clone method which should be used  |
      |                      |          | by the tool. Possible values: https, ssh   |
      +----------------------+----------+--------------------------------------------+
      | REMOTE_URL           |    X     | The git remote url.                        |
      +----------------------+----------+--------------------------------------------+

    SSH CONFIGURATION:
      These environment variables are only used if the clone method is set to "ssh".
      +======================+==========+============================================+
      | Environment variable | Required |                 Description                |
      +======================+==========+============================================+
      | KEY_FILE_PATH        |    X     | The absolute path to the private key       |
      |                      |          | which should be used for authentication.   |
      |                      |          | against the git server.                    |
      +----------------------+----------+--------------------------------------------+
      | SSH_PASSPHRASE       |          | The passphrase for the private key. If     |
      |                      |          | the variable is not set but needed by the  |
      |                      |          | key you will be prompted to enter the      |
      |                      |          | passphrase.                                |
      +----------------------+----------+--------------------------------------------+
      | SSH_KNOWN_HOSTS      |          | The absolute path to the known_hosts file  |
      |                      |          | which contains a valid entry for           |
      |                      |          | "gitlab.osram.info". If the variable is    |
      |                      |          | not set the following file locations       |
      |                      |          | will be used:                              |
      |                      |          | "~/.ssh/known_hosts"                       |
      |                      |          | "/etc/ssh/ssh_known_hosts"                 |
      +----------------------+----------+--------------------------------------------+

    HTTPS CONFIGURATION
      These environment variables are only used if the clone method is set to "https".
      +======================+==========+============================================+
      | Environment variable | Required |                 Description                |
      +======================+==========+============================================+
      | HTTPS_USER           |          | The username used for authentication       |
      |                      |          | against the git server. If the variable    |
      |                      |          | is not set you will be prompted to enter   |
      |                      |          | your username.                             |
      +----------------------+----------+--------------------------------------------+
      | HTTPS_PASSWORD       |          | The password used for authentication       |
      |                      |          | against the git server. If the variable    |
      |                      |          | is not set you will be prompted to enter   |
      |                      |          | the password.                              |
      +----------------------+----------+--------------------------------------------+

  DEPLOYMENT TARGET / SERVER STAGES DEFINITION
    There are three possible deployment targets used by the OSRAM workflow. These 
    deployment targets are preselectable by the "stage" flag.
    e.g "` + binaryName + ` deploy --stage dev"
    +===============================+================================================+
    |            Target             |                   Description                  |
    +===============================+================================================+
    | Development Server (dev)      | The DEV server is available to the agency for  |
    |                               | its own integration and internal acceptance    |
    |                               | testing.                                       |
    +-------------------------------+------------------------------------------------+
    | Quality Assurance Server (qa) | The QA server is available for OSRAM to accept |
    |                               | the release.                                   |
    +-------------------------------+------------------------------------------------+
    | Production Server (prod)      | The prod server is used to host the final.     |
    |                               | released application.                          |
    +-------------------------------+------------------------------------------------+

  VERSIONING
    The OSRAM workflow uses semantic versioning 2.0.0 
    (https://semver.org/spec/v2.0.0.html).

    [major].[minor].[patch]-[iteration]+[stage]
    e.g: "1.3.1-25+dev"

    How to choose which version to increment:
    +=============+===============================================================+
    |   Version   |                             When                              |
    +=============+===============================================================+
    | Patch       | Only backwards compatible bug fixes are introduced. A bug fix |
    |             | is defined as an change that fixes incorrect behavior.        |
    +-------------+---------------------------------------------------------------+
    | Minor       | Only backwards compatible functionality is introduced.        |
    +-------------+---------------------------------------------------------------+
    | Major       | backwards incompatible changes are introduced.                |
    +-------------+---------------------------------------------------------------+

  EXAMPLE RELEASE OF A NEW VERSION:

    INITIAL SITUATION:
    The master and develop branch are currently pointing to "Commit 1", nothing is
    deployed to any server yet, nothing is tagged yet.
    git log:
    +------------------------------------------------------+
    | * 072b23a - Commit 1 (origin/master, origin/develop) |
    +------------------------------------------------------+

    Meanwhile several developers committed and code to the develop branch.
    git log:
    +---------------------------------------+
    | * 39ac05d - Commit 5 (origin/develop) |
    | * 61051c4 - Commit 4                  |
    | * e0067af - Commit 3                  |
    | * 2c6d420 - Commit 2                  |
    | * 072b23a - Commit 1 (origin/master)  |
    +---------------------------------------+

    The development team decides that Commit 1-n (n because it is possible that the
    testing on the dev server results in additional bug fix commits ) should result
    in a new major version (major version zero "0.y.z" is for initial development,
    version 1.0.0 defines the public api). A developer deploys "Commit 3" to dev
    as first iteration of the new major version. 
    He executes: "` + binaryName + ` deploy --stage dev --type major" and selects
    "Commit 3". This results in a development server deployment and the following 
    git log:
    +---------------------------------------+
    | * 39ac05d - Commit 5 (origin/develop) |
    | * 61051c4 - Commit 4                  |
    | * e0067af - Commit 3 (1.0.0-1+dev)    |
    | * 2c6d420 - Commit 2                  |
    | * 072b23a - Commit 1 (origin/master)  |
    +---------------------------------------+

    The testing team/person tests commit 1-3 on the dev server. Meanwhile or after
    the testing a developer deploys the latest commit (Commit 5) to dev as the
    second iteration of the new major version (Note that it is not possible to
    deploy "Commit 1" or "Commit 2" to dev as seperate versions because these
    commits are already included in "1.0.0-1+dev").
    He executes: "` + binaryName + ` deploy --stage dev --latest --type increment".
    This results in a development server deployment and the following git log:
    +----------------------------------------------------+
    | * 39ac05d - Commit 5 (origin/develop, 1.0.0-2+dev) |
    | * 61051c4 - Commit 4                               |
    | * e0067af - Commit 3 (1.0.0-1+dev)                 |
    | * 2c6d420 - Commit 2                               |
    | * 072b23a - Commit 1 (origin/master)               |
    +----------------------------------------------------+

    The testing team/person tests commit 4-5 on the dev server and finds several
    severe bugs. The development team fixes these bugs and pushes "Commit 6" and
    "Commit 7" to develop. A developer deploys the latest commit (Commit 7) to dev
    as the third iteration of the new version.
    He executes: "` + binaryName + ` deploy --stage dev --latest --type increment".
    This results in a development server deployment and the following git log:
    +----------------------------------------------------+
    | * 6737421 - Commit 7 (origin/develop, 1.0.0-3+dev) |
    | * ee4b3dc - Commit 6                               |
    | * 39ac05d - Commit 5 (1.0.0-2+dev)                 |
    | * 61051c4 - Commit 4                               |
    | * e0067af - Commit 3 (1.0.0-1+dev)                 |
    | * 2c6d420 - Commit 2                               |
    | * 072b23a - Commit 1 (origin/master)               |
    +----------------------------------------------------+

    The testing team/person verifies that the bugs have been resolved and approves
    the new version. After the approval from the testing team/person a developer
    deploys the latest dev iteration (1.0.0-3+dev) to the qa server.
    He executes: "` + binaryName + ` deploy --stage qa --latest".
    This results in a qa server deployment and the following git log:
    +----------------------------------------------------------------+
    | * 6737421 - Commit 7 (origin/develop, 1.0.0-3+dev, 1.0.0-3+qa) |
    | * ee4b3dc - Commit 6                                           |
    | * 39ac05d - Commit 5 (1.0.0-2+dev)                             |
    | * 61051c4 - Commit 4                                           |
    | * e0067af - Commit 3 (1.0.0-1+dev)                             |
    | * 2c6d420 - Commit 2                                           |
    | * 072b23a - Commit 1 (origin/master)                           |
    +----------------------------------------------------------------+

    Note that it is not possible to deploy "1.0.0-2+dev" or "1.0.0-1+dev" to the
    qa server anymore because everything up to "1.0.0-3+dev" is already included 
    in the "1.0.0-3+qa" deployment. 
    The OSRAM testing team finds another bug which needs to be resolved before the
    production release of the new version. The development team fixes this bugs,
    pushes the new commit ("Commit 8") to develop and deploys the new commit as
    new iteration to develop.
    He executes: "` + binaryName + ` deploy --stage dev --latest --type increment"
    and after validating that the bug has been resolved on dev he executes:
    "` + binaryName + ` deploy --stage qa --latest". That results in the following 
    git log:
    +----------------------------------------------------------------+
    | * 57761ba - Commit 8 (origin/develop, 1.0.0-4+dev, 1.0.0-4+qa) |
    | * 6737421 - Commit 7 (1.0.0-3+dev, 1.0.0-3+qa)                 |
    | * ee4b3dc - Commit 6                                           |
    | * 39ac05d - Commit 5 (1.0.0-2+dev)                             |
    | * 61051c4 - Commit 4                                           |
    | * e0067af - Commit 3 (1.0.0-1+dev)                             |
    | * 2c6d420 - Commit 2                                           |
    | * 072b23a - Commit 1 (origin/master)                           |
    +----------------------------------------------------------------+

    The OSRAM testing team validates that the bug has been resolved and approves
    the new version. After that final testing a developer deploys the current
    qa iteration to production.
    He executes: "` + binaryName + ` deploy --stage prod --latest". That results in
    the following git log:
    +-----------------------------------------------------------------------+
    | * 57761ba - Commit 8 (origin/develop, 1.0.0-4+dev, 1.0.0-4+qa)        |
    |                      (origin/master, 1.0.0)                           |
    | * 6737421 - Commit 7 (1.0.0-3+dev, 1.0.0-3+qa)                        |
    | * ee4b3dc - Commit 6                                                  |
    | * 39ac05d - Commit 5 (1.0.0-2+dev)                                    |
    | * 61051c4 - Commit 4                                                  |
    | * e0067af - Commit 3 (1.0.0-1+dev)                                    |
    | * 2c6d420 - Commit 2                                                  |
    | * 072b23a - Commit 1                                                  |
    +-----------------------------------------------------------------------+

    Now everything up to "Commit 8" is deployed to the production server as
    version "1.0.0". After the production deployment the version is marked
    as deployed which prevents any additional 1.0.0-n+dev or 1.0.0-n+qa
    incrementation.
`)
	return nil
}
