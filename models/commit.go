package models

import "strings"

const PrintMessageLen = 40

type Commit struct {
	Author  string
	Message string
	Hash    string
	Tag     string
}

func (this Commit) GetId() string {
	return this.Hash
}

func (this Commit) GetFormattedText() string {
	return this.Hash[0:10] + " " + this.GetFormattedMessage() + " <" + this.Author + ">"
}

func (this Commit) GetDisplayText() string {
	return this.Hash[0:10] + " " + this.Message + " <" + this.Author + ">"
}

func (this Commit) GetFormattedMessage() string {
	messageLen := len(this.Message)
	if messageLen < PrintMessageLen {
		return this.Message + strings.Repeat(" ", PrintMessageLen-messageLen)
	} else {
		return this.Message[0:PrintMessageLen]
	}
}
