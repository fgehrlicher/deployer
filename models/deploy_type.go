package models

type DeployType struct {
	id          string
	description string
}

func NewDeployType(id string, description string) *DeployType {
	return &DeployType{
		id:          id,
		description: description,
	}
}

func (this DeployType) GetId() string {
	return this.id
}

func (this DeployType) GetDisplayText() string {
	return this.description
}

func (this DeployType) GetFormattedText() string {
	return this.description
}
