package keyboard

import (
	"sccreeper/goputer/pkg/constants"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var keyMappings map[int32]constants.KeyCode

func init() {
	keyMappings = make(map[int32]constants.KeyCode)

	keyMappings[rl.KeyA] = constants.KeyA
	keyMappings[rl.KeyB] = constants.KeyB
	keyMappings[rl.KeyC] = constants.KeyC
	keyMappings[rl.KeyD] = constants.KeyD
	keyMappings[rl.KeyE] = constants.KeyE
	keyMappings[rl.KeyF] = constants.KeyF
	keyMappings[rl.KeyG] = constants.KeyG
	keyMappings[rl.KeyH] = constants.KeyH
	keyMappings[rl.KeyI] = constants.KeyI
	keyMappings[rl.KeyJ] = constants.KeyJ
	keyMappings[rl.KeyK] = constants.KeyK
	keyMappings[rl.KeyL] = constants.KeyL
	keyMappings[rl.KeyM] = constants.KeyM
	keyMappings[rl.KeyN] = constants.KeyN
	keyMappings[rl.KeyO] = constants.KeyO
	keyMappings[rl.KeyP] = constants.KeyP
	keyMappings[rl.KeyQ] = constants.KeyQ
	keyMappings[rl.KeyR] = constants.KeyR
	keyMappings[rl.KeyS] = constants.KeyS
	keyMappings[rl.KeyT] = constants.KeyT
	keyMappings[rl.KeyU] = constants.KeyU
	keyMappings[rl.KeyV] = constants.KeyV
	keyMappings[rl.KeyW] = constants.KeyW
	keyMappings[rl.KeyX] = constants.KeyX
	keyMappings[rl.KeyY] = constants.KeyY
	keyMappings[rl.KeyZ] = constants.KeyZ

	keyMappings[rl.KeySpace] = constants.KeySpace

	keyMappings[rl.KeyOne] = constants.KeyDigit1
	keyMappings[rl.KeyTwo] = constants.KeyDigit2
	keyMappings[rl.KeyThree] = constants.KeyDigit3
	keyMappings[rl.KeyFour] = constants.KeyDigit4
	keyMappings[rl.KeyFive] = constants.KeyDigit5
	keyMappings[rl.KeySix] = constants.KeyDigit6
	keyMappings[rl.KeySeven] = constants.KeyDigit7
	keyMappings[rl.KeyEight] = constants.KeyDigit8
	keyMappings[rl.KeyNine] = constants.KeyDigit9
	keyMappings[rl.KeyZero] = constants.KeyDigit0

	keyMappings[rl.KeyMinus] = constants.KeyMinus
	keyMappings[rl.KeyEqual] = constants.KeyEquals

	keyMappings[rl.KeyLeftBracket] = constants.KeyBracketLeft
	keyMappings[rl.KeyRightBracket] = constants.KeyBracketRight

	keyMappings[rl.KeySemicolon] = constants.KeySemicolon
	keyMappings[rl.KeyApostrophe] = constants.KeyQuote

	keyMappings[rl.KeyComma] = constants.KeyComma
	keyMappings[rl.KeyPeriod] = constants.KeyPeriod

	keyMappings[rl.KeySlash] = constants.KeySlash
	keyMappings[rl.KeyBackSlash] = constants.KeyBackslash

	keyMappings[rl.KeyLeftControl] = constants.KeyLeftCtrl
	keyMappings[rl.KeyRightControl] = constants.KeyRightCtrl

	keyMappings[rl.KeyLeftAlt] = constants.KeyLeftAlt
	keyMappings[rl.KeyRightAlt] = constants.KeyRightAlt

	keyMappings[rl.KeyLeftShift] = constants.KeyLeftShift
	keyMappings[rl.KeyRightShift] = constants.KeyRightShift

	keyMappings[rl.KeyCapsLock] = constants.KeyCapsLock
	keyMappings[rl.KeyTab] = constants.KeyTab
	keyMappings[rl.KeyEscape] = constants.KeyEscape
	keyMappings[rl.KeyGrave] = constants.KeyBackquote

	keyMappings[rl.KeyBackspace] = constants.KeyBackspace
	keyMappings[rl.KeyEnter] = constants.KeyReturn

	keyMappings[rl.KeyLeftSuper] = constants.KeySuperLeft
	keyMappings[rl.KeyRightSuper] = constants.KeySuperRight

	keyMappings[rl.KeyMenu] = constants.KeyMenu

	keyMappings[rl.KeyPrintScreen] = constants.KeyPrintScreen
	keyMappings[rl.KeyScrollLock] = constants.KeyScrollLock
	keyMappings[rl.KeyPause] = constants.KeyPauseBreak

	keyMappings[rl.KeyInsert] = constants.KeyInsert
	keyMappings[rl.KeyHome] = constants.KeyHome
	keyMappings[rl.KeyDelete] = constants.KeyDelete
	keyMappings[rl.KeyEnd] = constants.KeyEnd

	keyMappings[rl.KeyPageUp] = constants.KeyPageUp
	keyMappings[rl.KeyPageDown] = constants.KeyPageDown

	keyMappings[rl.KeyUp] = constants.KeyArrowUp
	keyMappings[rl.KeyRight] = constants.KeyArrowRight
	keyMappings[rl.KeyDown] = constants.KeyArrowDown
	keyMappings[rl.KeyLeft] = constants.KeyArrowLeft

	keyMappings[rl.KeyNumLock] = constants.KeyNumLock
	keyMappings[rl.KeyKpDivide] = constants.KeyNPDivide
	keyMappings[rl.KeyKpMultiply] = constants.KeyNPMultiply
	keyMappings[rl.KeyKpSubtract] = constants.KeyNPMinus
	keyMappings[rl.KeyKpAdd] = constants.KeyNPPlus
	keyMappings[rl.KeyKpEnter] = constants.KeyNPReturn

	keyMappings[rl.KeyKp1] = constants.KeyNP1
	keyMappings[rl.KeyKp2] = constants.KeyNP2
	keyMappings[rl.KeyKp3] = constants.KeyNP3
	keyMappings[rl.KeyKp4] = constants.KeyNP4
	keyMappings[rl.KeyKp5] = constants.KeyNP5
	keyMappings[rl.KeyKp6] = constants.KeyNP6
	keyMappings[rl.KeyKp7] = constants.KeyNP7
	keyMappings[rl.KeyKp8] = constants.KeyNP8
	keyMappings[rl.KeyKp9] = constants.KeyNP9
	keyMappings[rl.KeyKp0] = constants.KeyNP0
	keyMappings[rl.KeyKpDecimal] = constants.KeyNPDecimal

	keyMappings[rl.KeyF1] = constants.KeyF1
	keyMappings[rl.KeyF2] = constants.KeyF2
	keyMappings[rl.KeyF3] = constants.KeyF3
	keyMappings[rl.KeyF4] = constants.KeyF4
	keyMappings[rl.KeyF5] = constants.KeyF5
	keyMappings[rl.KeyF6] = constants.KeyF6
	keyMappings[rl.KeyF7] = constants.KeyF7
	keyMappings[rl.KeyF8] = constants.KeyF8
	keyMappings[rl.KeyF9] = constants.KeyF9
	keyMappings[rl.KeyF10] = constants.KeyF10
	keyMappings[rl.KeyF11] = constants.KeyF11
	keyMappings[rl.KeyF12] = constants.KeyF12

}

func MapKey(key int32) constants.KeyCode {

	if _, exists := keyMappings[key]; exists {
		return keyMappings[key]
	} else {
		return 0
	}

}
