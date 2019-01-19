package cli_util

import (
	"fmt"
	"bufio"
	"os"
	"strconv"
	"errors"
	"strings"
)

type Selectable interface {
	GetId() string
	GetFormattedText() string
}

type SingleSelect struct {
	SelectText string
	Items      []Selectable
}

func NewSingleSelect(items []Selectable, selectText string) *SingleSelect {
	return &SingleSelect{
		Items:      items,
		SelectText: selectText,
	}
}

func (this *SingleSelect) RenderSelect() (Selectable, error) {
	var selectedItem Selectable

	for {
		this.renderItems()
		itemKey, err := this.selectItem()
		this.clear()
		if err == nil {
			selectedItem = this.Items[itemKey]
			break
		}
	}

	return selectedItem, nil
}

func (this *SingleSelect) renderItems() {
	itemsString := LF
	for i := 0; i < len(this.Items); i++ {
		itemsString = itemsString + " "
		itemsString = itemsString + fmt.Sprintf("[%v]", i+1)
		if len(this.Items) >= 10 {
			itemsString = itemsString + " "
		}
		if i+1 < 10 {
			itemsString = itemsString + " "
		}
		itemsString = itemsString + fmt.Sprintf("%v"+LF, this.Items[i].GetFormattedText())
	}
	itemsString = itemsString + LF
	fmt.Print(itemsString)
}

func (this *SingleSelect) selectItem() (int, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(" " + this.SelectText + " ")
	text, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}
	text = strings.TrimRight(text, LF)
	text = strings.TrimRight(text, CR)

	key, err := strconv.Atoi(text)
	if err != nil {
		return 0, err
	}
	key = key - 1

	if key < 0 || key >= len(this.Items) {
		return 0, errors.New("invalid choice")
	}

	return key, nil
}

func (this *SingleSelect) clear() {
	additionLines := 3
	for i := 0; i < len(this.Items)+additionLines; i++ {
		fmt.Print(DeleteCurrentLine + CursorUp)
	}
}
