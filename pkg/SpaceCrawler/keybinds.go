package SpaceCrawler

import "fmt"

// KeyBind represents a single key that can do an action
type KeyBind struct {
	Name        string
	Keys        []rune
	Description string
	Action      func()
}

// NewKeyBind returns a KeyBind with the given name and keys
func NewKeyBind(name string, keys []rune) *KeyBind {
	return &KeyBind{
		Name:        name,
		Keys:        keys,
		Description: "",
		Action:      func() {},
	}
}

// String returns string representation of key
func (k *KeyBind) String() string {
	var keys = ""
	for i, key := range k.Keys {
		var sep = "/"
		if i == len(k.Keys)-1 {
			sep = ""
		}

		keys += string(key) + sep
	}

	if len(k.Description) > 0 {
		return fmt.Sprintf("%s %s\n%s", keys, k.Name, k.Description)
	} else {
		return fmt.Sprintf("%s %s", keys, k.Name)
	}
}

// WithDescription adds a description to the KeyBind
func (k *KeyBind) WithDescription(description string) *KeyBind {
	k.Description = description
	return k
}

// WithAction adds an action that can be called to the KeyBind
func (k *KeyBind) WithAction(action func()) *KeyBind {
	k.Action = action
	return k
}

// KeyMap is a structure for holding multiple KeyBind
// Only access properties for reading
// We're all grown-ups here
type KeyMap struct {
	Keys     []*KeyBind
	Bindings map[rune]*KeyBind
}

// NewKeyMap returns a new KeyMap
func NewKeyMap() *KeyMap {
	return &KeyMap{
		Keys:     make([]*KeyBind, 0),
		Bindings: make(map[rune]*KeyBind),
	}
}

// Add a new KeyBind to the KeyMap
func (k *KeyMap) Add(bind *KeyBind) *KeyMap {
	k.Keys = append(k.Keys, bind)
	for _, key := range bind.Keys {
		k.Bindings[key] = bind
	}
	return k
}

// Action will call the action associated with the given key
func (k *KeyMap) Action(key rune) bool {
	var bind, ok = k.Bindings[key]
	if !ok {
		return false
	}

	bind.Action()
	return true
}
