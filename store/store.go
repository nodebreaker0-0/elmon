package store

var GlobalState GlobalStateType

func init() {
	// I don't want but
	// It has dependency to main.go
	GlobalState = GlobalStateType{
		ELs: make(map[string]*ELType),
	}
}
