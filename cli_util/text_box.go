package cli_util

import (
	"fmt"
	"strings"
)

type Showable interface {
	GetId() string
	GetDisplayText() string
}

type TextBoxElement struct {
	displayName string
	Item        Showable
}

func NewTextBoxElement(displayName string, item Showable) *TextBoxElement {
	return &TextBoxElement{
		displayName: displayName,
		Item:        item,
	}
}

func (this *TextBoxElement) SetDisplayName(displayName string) {
	this.displayName = displayName
}

type TextBox struct {
	Elements map[int]*TextBoxElement
	suffix   string
}

func NewTextBox() *TextBox {
	return &TextBox{Elements: make(map[int]*TextBoxElement)}
}

func (this *TextBox) RenderBox() {
	var (
		leftUpperCorner      = "╔"
		leftBottomCorner     = "╚"
		rightUpperCorner     = "╗"
		rightBottomCorner    = "╝"
		horizontalSideBorder = "║"
		verticalSideBorder   = "═"
		maxKeyLength         = 0
		maxElementLength     = 0
	)

	for _, textBoxElement := range this.Elements {
		key := (*textBoxElement).displayName
		var element string
		if (*textBoxElement).Item != nil {
			element = (*textBoxElement).Item.GetDisplayText()
		}

		currentKeyLength := len(key)
		currentElementLength := len(element)
		if currentElementLength > maxElementLength {
			maxElementLength = currentElementLength
		}
		if currentKeyLength > maxKeyLength {
			maxKeyLength = currentKeyLength
		}
	}
	maxRowWidth := 5 + maxKeyLength + maxElementLength
	fmt.Println(leftUpperCorner + strings.Repeat(verticalSideBorder, maxRowWidth) + rightUpperCorner)

	for i := 0; i < len(this.Elements); i ++ {
		textBoxElement := this.Elements[i]
		key := (*textBoxElement).displayName
		var element string
		if (*textBoxElement).Item != nil {
			element = (*textBoxElement).Item.GetDisplayText()
		}

		keyLenght := len(key)
		elementLenght := len(element)

		row := horizontalSideBorder + " " +
			key + strings.Repeat(" ", maxKeyLength-keyLenght) + " : " +
			element + strings.Repeat(" ", maxElementLength-elementLenght) +
			" " + horizontalSideBorder

		fmt.Println(row)
	}

	fmt.Println(leftBottomCorner + strings.Repeat(verticalSideBorder, maxRowWidth) + rightBottomCorner)

	if this.suffix != "" {
		fmt.Print(this.suffix)
	}
}

func (this TextBox) Clear() {
	clearLines := len(this.Elements) + 2

	if this.suffix != "" {
		lineBreaks := strings.Count(this.suffix, LF)
		clearLines = clearLines + lineBreaks
	}

	for i := 0; i < clearLines; i++ {
		fmt.Print(CursorUp + DeleteCurrentLine)
	}
}

func (this *TextBox) ReRenderBox() {
	this.Clear()
	this.RenderBox()
}

func (this *TextBox) SetSuffix(suffix string) {
	if suffix != "" {
		this.Clear()
		this.suffix = fmt.Sprintf(LF+" %v"+LF, suffix)
		this.RenderBox()
	}
}
