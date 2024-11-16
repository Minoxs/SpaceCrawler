package SpaceCrawler

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// KeyView is the structure responsible for rendering out a KeyMap
type KeyView struct {
	box  *tview.Flex
	keys *KeyMap
}

// NewKeyView returns a new KeyView for a given KeyMap
func NewKeyView(keys *KeyMap) *KeyView {
	return &KeyView{
		box:  tview.NewFlex(),
		keys: keys,
	}
}

// Model returns the tview.Primitive to be used for rendering
func (k *KeyView) Model() tview.Primitive {
	for _, key := range k.keys.Keys {
		k.box.AddItem(
			tview.NewTextView().
				SetText(key.String()).
				SetTextAlign(tview.AlignCenter).
				SetTextColor(tcell.ColorGrey),
			0,
			1,
			false,
		)
	}
	return k.box
}
