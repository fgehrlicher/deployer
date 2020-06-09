# Deployer

This tool can help you to properly increment tags for this semver (https://semver.org/spec/v2.0.0.html) based schema:

[major].[minor].[patch]-[iteration]+[stage]  
e.g: "1.3.1-25+dev"

<p align="center"><img src="/render1591687974072.gif?raw=true"/></p>

## Prerequisites
(Older versions of the following tools might work but are not tested)
* [Docker](https://docs.docker.com/get-started/#download-and-install-docker-desktop) >= 19.03.5 
* GNU Make >= 3.81
* Bash >= 5.0.0

## Build 
Run ```make``` to get a quick overview of the available commands.  
eg. to build the docker image:
```sh
make build-container
```

If you dont want to build the image yourself, you can get it from the github docker registry:  
https://github.com/fgehrlicher/deployer/packages/259619   
eg:
```sh
docker pull docker.pkg.github.com/fgehrlicher/deployer/deployer:latest
```

## Usage
```
NAME:
   deployer - deployer command [command options]

USAGE:
   deployer [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
     validate, v  validates the current settings
     manual, m    shows the deployment manual
     deploy, d    marks an commit for deployment
     help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug        enable debug mode
   --help, -h     show help
   --version, -v  print the version

```

## Run

You must mount your private key which is used to access the git repository.
``REMOTE_URL`` must be set to the ssh remote of the repository to be tagged.  
e.g:
```
docker run -it --rm \
-v ~/.ssh/id_rsa:/root/.ssh/id_rsa \
-v ~/.ssh/known_hosts:/root/ssh/known_hosts \
-e REMOTE_URL="git@github.com:fgehrlicher/deployer.git" \
docker.pkg.github.com/fgehrlicher/deployer/deployer:latest deploy
```
