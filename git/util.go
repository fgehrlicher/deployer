package git

import (
	"os/exec"
	"regexp"
	"errors"
	"os"
	)

var (
	invalidSshAuthSocketError = errors.New("invalid SSH_AUTH_SOCK")
	invalidAgentPidError      = errors.New("invalid SSH_AGENT_PID")
)

func StartSshAgent() error {
	output, err := exec.Command("/usr/bin/ssh-agent").Output()
	if err != nil {
		return err
	}
	outputString := string(output)

	authSocketMatch := regexp.MustCompile("SSH_AUTH_SOCK=(.*?);")
	if !authSocketMatch.MatchString(outputString) {
		return invalidSshAuthSocketError
	}

	agentPidMatch := regexp.MustCompile("SSH_AGENT_PID=(.*?);")
	if !agentPidMatch.MatchString(outputString) {
		return invalidAgentPidError
	}

	subMatchAuthSocket := authSocketMatch.FindAllStringSubmatch(outputString, -1)
	authSocketEnv := subMatchAuthSocket[0][1]
	if authSocketEnv == "" {
		return invalidSshAuthSocketError
	}

	subMatchAgentPid := agentPidMatch.FindAllStringSubmatch(outputString, -1)
	agentPidEnv := subMatchAgentPid[0][1]
	if agentPidEnv == "" {
		return invalidAgentPidError
	}

	err = os.Setenv("SSH_AUTH_SOCK", authSocketEnv)
	if err != nil {
		return err
	}

	return os.Setenv("SSH_AGENT_PID", agentPidEnv)
}
