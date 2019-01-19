package models

type Target struct {
	id   string
	name string
}

func NewTarget(id string, name string) *Target {
	return &Target{
		id:   id,
		name: name,
	}
}

func (this Target) GetId() string {
	return this.id
}

func (this Target) GetDisplayText() string {
	return this.name
}

func (this Target) GetFormattedText() string {
	return this.name
}
