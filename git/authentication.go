package git

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/howeyc/gopass"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"

	"github.com/fgehrlicher/deployer/cli_util"
	"github.com/fgehrlicher/deployer/config"
)

const HttpsType = "https"
const SshType = "ssh"

func GetAuth(config config.Config) (transport.AuthMethod, error) {
	var authenticator transport.AuthMethod
	var linesToDelete int

	switch config.CloneMethod {
	case SshType:
		err := StartSshAgent()
		if err != nil {
			return nil, err
		}
		keyFile, err := LoadSshKey(config.KeyFilePath)
		if err != nil {
			return nil, err
		}
		authenticator, err = GetSshAuthentication(keyFile, "git", config.SshPassPhrase)
		if err != nil {
			return nil, err
		}
	case HttpsType:
		if len(config.HttpsUserName) == 0 || len(config.HttpsPassword) == 0 {
			reader := bufio.NewReader(os.Stdin)
			if len(config.HttpsUserName) == 0 {
				fmt.Print("Enter Git Username: ")
				input, _ := reader.ReadString('\n')
				input = strings.TrimRight(input, cli_util.LF)
				input = strings.TrimRight(input, cli_util.CR)
				if len(input) == 0 {
					return nil, InvalidUsername
				}
				config.HttpsUserName = input
				linesToDelete++
			}
			if len(config.HttpsPassword) == 0 {
				fmt.Print("Enter Password: ")
				password, err := gopass.GetPasswdMasked()
				if err != nil {
					return nil, err
				}
				input := string(password)
				if len(input) == 0 {
					return nil, InvalidPassword
				}
				config.HttpsPassword = input
				linesToDelete++
			}
		}
		authenticator = GetHttpsAuthentication(config.HttpsUserName, config.HttpsPassword)
	default:
		return nil, InvalidCloneMethod
	}

	for i := 0; i < linesToDelete; i++ {
		fmt.Print(cli_util.DeleteCurrentLine + cli_util.CursorUp)
	}

	return authenticator, nil
}

func GetHttpsAuthentication(userName string, password string) *http.BasicAuth {
	return &http.BasicAuth{
		Username: userName,
		Password: password,
	}
}

func GetSshAuthentication(file []byte, user string, passPhrase string) (*ssh.PublicKeys, error) {
	return ssh.NewPublicKeys(
		user,
		file,
		passPhrase,
	)
}

func LoadSshKey(keyFilePath string) ([]byte, error) {
	return ioutil.ReadFile(keyFilePath)
}
