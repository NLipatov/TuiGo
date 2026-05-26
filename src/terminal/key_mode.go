package terminal

type KeyMod uint8

const (
	ModNone  KeyMod = 0 // 0000_0000
	ModCtrl  KeyMod = 1 // 0000_0001
	ModAlt   KeyMod = 2 // 0000_0010
	ModShift KeyMod = 4 // 0000_0100
)
