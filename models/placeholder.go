package models

type Placeholder string

func NewPlaceholder(placeholderString string) *Placeholder {
	placeholder := Placeholder(placeholderString)
	return &placeholder
}

func (this Placeholder) GetId() string {
	return string(this)
}

func (this Placeholder) GetDisplayText() string {
	return string(this)
}

func (this Placeholder) GetFormattedText() string {
	return string(this)
}
