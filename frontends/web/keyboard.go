//go:build js

package main

import (
	"sccreeper/goputer/pkg/constants"
	"syscall/js"
)

var keyboardMap map[string]constants.KeyCode

func init() {
	keyboardMap = make(map[string]constants.KeyCode)

	keyboardMap["KeyA"] = constants.KeyA
	keyboardMap["KeyB"] = constants.KeyB
	keyboardMap["KeyC"] = constants.KeyC
	keyboardMap["KeyD"] = constants.KeyD
	keyboardMap["KeyE"] = constants.KeyE
	keyboardMap["KeyF"] = constants.KeyF
	keyboardMap["KeyG"] = constants.KeyG
	keyboardMap["KeyH"] = constants.KeyH
	keyboardMap["KeyI"] = constants.KeyI
	keyboardMap["KeyJ"] = constants.KeyJ
	keyboardMap["KeyK"] = constants.KeyK
	keyboardMap["KeyL"] = constants.KeyL
	keyboardMap["KeyM"] = constants.KeyM
	keyboardMap["KeyN"] = constants.KeyN
	keyboardMap["KeyO"] = constants.KeyO
	keyboardMap["KeyP"] = constants.KeyP
	keyboardMap["KeyQ"] = constants.KeyQ
	keyboardMap["KeyR"] = constants.KeyR
	keyboardMap["KeyS"] = constants.KeyS
	keyboardMap["KeyT"] = constants.KeyT
	keyboardMap["KeyU"] = constants.KeyU
	keyboardMap["KeyV"] = constants.KeyV
	keyboardMap["KeyW"] = constants.KeyW
	keyboardMap["KeyX"] = constants.KeyX
	keyboardMap["KeyY"] = constants.KeyY
	keyboardMap["KeyZ"] = constants.KeyZ

	keyboardMap["Digit1"] = constants.KeyDigit1
	keyboardMap["Digit2"] = constants.KeyDigit2
	keyboardMap["Digit3"] = constants.KeyDigit3
	keyboardMap["Digit4"] = constants.KeyDigit4
	keyboardMap["Digit5"] = constants.KeyDigit5
	keyboardMap["Digit6"] = constants.KeyDigit6
	keyboardMap["Digit7"] = constants.KeyDigit7
	keyboardMap["Digit8"] = constants.KeyDigit8
	keyboardMap["Digit9"] = constants.KeyDigit9
	keyboardMap["Digit0"] = constants.KeyDigit0

	keyboardMap["Space"] = constants.KeySpace

	keyboardMap["Minus"] = constants.KeyMinus
	keyboardMap["Equal"] = constants.KeyEquals
	keyboardMap["BracketLeft"] = constants.KeyBracketLeft
	keyboardMap["BracketRight"] = constants.KeyBracketRight
	keyboardMap["Semicolon"] = constants.KeySemicolon
	keyboardMap["Quote"] = constants.KeyQuote
	keyboardMap["Comma"] = constants.KeyComma
	keyboardMap["Period"] = constants.KeyPeriod
	keyboardMap["Slash"] = constants.KeySlash
	keyboardMap["Backslash"] = constants.KeyBackslash

	keyboardMap["ControlLeft"] = constants.KeyLeftCtrl
	keyboardMap["ControlRight"] = constants.KeyRightCtrl
	keyboardMap["AltLeft"] = constants.KeyLeftAlt
	keyboardMap["AltRight"] = constants.KeyRightAlt
	keyboardMap["ShiftLeft"] = constants.KeyLeftShift
	keyboardMap["ShiftRight"] = constants.KeyRightShift

	keyboardMap["CapsLock"] = constants.KeyCapsLock
	keyboardMap["Tab"] = constants.KeyTab
	keyboardMap["Escape"] = constants.KeyEscape
	keyboardMap["Backquote"] = constants.KeyBackquote
	keyboardMap["Backspace"] = constants.KeyBackspace
	keyboardMap["Enter"] = constants.KeyReturn
	keyboardMap["MetaLeft"] = constants.KeySuperLeft
	keyboardMap["MetaRight"] = constants.KeySuperRight
	keyboardMap["ContextMenu"] = constants.KeyMenu

	keyboardMap["PrintScreen"] = constants.KeyPrintScreen
	keyboardMap["ScrollLock"] = constants.KeyScrollLock
	keyboardMap["Pause"] = constants.KeyPauseBreak

	keyboardMap["Insert"] = constants.KeyInsert
	keyboardMap["Home"] = constants.KeyHome
	keyboardMap["Delete"] = constants.KeyDelete
	keyboardMap["End"] = constants.KeyEnd
	keyboardMap["PageUp"] = constants.KeyPageUp
	keyboardMap["PageDown"] = constants.KeyPageDown
	keyboardMap["ArrowUp"] = constants.KeyArrowUp
	keyboardMap["ArrowRight"] = constants.KeyArrowRight
	keyboardMap["ArrowDown"] = constants.KeyArrowDown
	keyboardMap["ArrowLeft"] = constants.KeyArrowLeft

	keyboardMap["NumLock"] = constants.KeyNumLock
	keyboardMap["NumpadDivide"] = constants.KeyNPDivide
	keyboardMap["NumpadMultiply"] = constants.KeyNPMultiply
	keyboardMap["NumpadSubtract"] = constants.KeyNPMinus
	keyboardMap["NumpadAdd"] = constants.KeyNPPlus
	keyboardMap["NumpadEnter"] = constants.KeyNPReturn
	keyboardMap["Numpad1"] = constants.KeyNP1
	keyboardMap["Numpad2"] = constants.KeyNP2
	keyboardMap["Numpad3"] = constants.KeyNP3
	keyboardMap["Numpad4"] = constants.KeyNP4
	keyboardMap["Numpad5"] = constants.KeyNP5
	keyboardMap["Numpad6"] = constants.KeyNP6
	keyboardMap["Numpad7"] = constants.KeyNP7
	keyboardMap["Numpad8"] = constants.KeyNP8
	keyboardMap["Numpad9"] = constants.KeyNP9
	keyboardMap["Numpad0"] = constants.KeyNP0
	keyboardMap["NumpadDecimal"] = constants.KeyNPDecimal

	keyboardMap["F1"] = constants.KeyF1
	keyboardMap["F2"] = constants.KeyF2
	keyboardMap["F3"] = constants.KeyF3
	keyboardMap["F4"] = constants.KeyF4
	keyboardMap["F5"] = constants.KeyF5
	keyboardMap["F6"] = constants.KeyF6
	keyboardMap["F7"] = constants.KeyF7
	keyboardMap["F8"] = constants.KeyF8
	keyboardMap["F9"] = constants.KeyF9
	keyboardMap["F10"] = constants.KeyF10
	keyboardMap["F11"] = constants.KeyF11
	keyboardMap["F12"] = constants.KeyF12
}

func mapKeycode(this js.Value, args []js.Value) any {

	if _, exists := keyboardMap[args[0].String()]; exists {
		return js.ValueOf(int(keyboardMap[args[0].String()]))
	} else {
		return js.ValueOf(0)
	}

}
