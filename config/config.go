package config

import (
	"fmt"
	"os"
	"text/tabwriter"
)

type Config struct {
	InputLineEndings string
	CloneMethod      string
	RemoteUrl        string
	HttpsUserName    string
	HttpsPassword    string
	KeyFilePath      string
	SshPassPhrase    string
}

func LoadConfig() Config {
	return Config{
		CloneMethod:   os.Getenv("CLONE_METHOD"),
		RemoteUrl:     os.Getenv("REMOTE_URL"),
		HttpsUserName: os.Getenv("HTTPS_USER"),
		HttpsPassword: os.Getenv("HTTPS_PASSWORD"),
		KeyFilePath:   os.Getenv("KEY_FILE_PATH"),
		SshPassPhrase: os.Getenv("SSH_PASSPHRASE"),
	}
}

func (this *Config) PrintConfig() {
	tabWriter := new(tabwriter.Writer)
	tabWriter.Init(os.Stdout, 0, 4, 0, '\t', 0)
	defer tabWriter.Flush()

	var buffer []byte

	buffer = append(buffer, fmt.Sprintf("\n %s\t%s", "name", "Value")...)
	buffer = append(buffer, fmt.Sprintf("\n %s\t%s", "----", "----")...)
	buffer = append(buffer, fmt.Sprintf("\n %s\t%s", "CLONE_METHOD", this.CloneMethod)...)

	switch this.CloneMethod {
	case "https":
		if this.HttpsUserName != "" {
			buffer = append(buffer, fmt.Sprintf("\n %s\t%s", "HTTPS_USER", this.HttpsUserName)...)
		}
		if this.HttpsPassword != "" {
			buffer = append(buffer, fmt.Sprintf("\n %s\t%s", "HTTPS_PASSWORD", "************")...)
		}
	case "ssh":
		if this.KeyFilePath != "" {
			buffer = append(buffer, fmt.Sprintf("\n %s\t%s", "KEY_FILE_PATH", this.KeyFilePath)...)
		}
		if this.SshPassPhrase != "" {
			buffer = append(buffer, fmt.Sprintf("\n %s\t%s", "SSH_PASSPHRASE", this.SshPassPhrase)...)
		}
	}
	buffer = append(buffer, fmt.Sprintf("\n %s\t%s\n", "REMOTE_URL", this.RemoteUrl)...)

	tabWriter.Write(buffer)
}
