package constants

type KeyCode uint32

// Some names are taken from https://developer.mozilla.org/en-US/docs/Web/API/UI_Events/Keyboard_event_key_values
// Keys are taken from 104 key US QWERTY layout + some more F keys
// https://en.wikipedia.org/wiki/Keyboard_layout#/media/File:Qwerty.svg
// I'm not going to add more, because in theory you can have infinitely many keys and I've had enough

const (
	KeyA KeyCode = 1 // 0 is reserved for "unknown" keys
	KeyB KeyCode = 2
	KeyC KeyCode = 3
	KeyD KeyCode = 4
	KeyE KeyCode = 5
	KeyF KeyCode = 6
	KeyG KeyCode = 7
	KeyH KeyCode = 8
	KeyI KeyCode = 9
	KeyJ KeyCode = 10
	KeyK KeyCode = 11
	KeyL KeyCode = 12
	KeyM KeyCode = 13
	KeyN KeyCode = 14
	KeyO KeyCode = 15
	KeyP KeyCode = 16
	KeyQ KeyCode = 17
	KeyR KeyCode = 18
	KeyS KeyCode = 19
	KeyT KeyCode = 20
	KeyU KeyCode = 21
	KeyV KeyCode = 22
	KeyW KeyCode = 23
	KeyX KeyCode = 24
	KeyY KeyCode = 25
	KeyZ KeyCode = 26

	KeySpace KeyCode = 27

	KeyDigit1 KeyCode = 28
	KeyDigit2 KeyCode = 29
	KeyDigit3 KeyCode = 30
	KeyDigit4 KeyCode = 31
	KeyDigit5 KeyCode = 32
	KeyDigit6 KeyCode = 33
	KeyDigit7 KeyCode = 34
	KeyDigit8 KeyCode = 35
	KeyDigit9 KeyCode = 36
	KeyDigit0 KeyCode = 37

	KeyMinus  KeyCode = 38
	KeyEquals KeyCode = 39

	KeyBracketLeft  KeyCode = 40
	KeyBracketRight KeyCode = 41

	KeySemicolon KeyCode = 42
	KeyQuote     KeyCode = 43

	KeyComma  KeyCode = 44
	KeyPeriod KeyCode = 45

	KeySlash     KeyCode = 46
	KeyBackslash KeyCode = 47

	KeyLeftCtrl  KeyCode = 48
	KeyRightCtrl KeyCode = 49

	KeyLeftAlt  KeyCode = 50
	KeyRightAlt KeyCode = 51

	KeyLeftShift  KeyCode = 52
	KeyRightShift KeyCode = 53

	KeyCapsLock  KeyCode = 54
	KeyTab       KeyCode = 55
	KeyEscape    KeyCode = 56
	KeyBackquote KeyCode = 57

	KeyBackspace KeyCode = 58
	KeyReturn    KeyCode = 59

	KeySuperLeft  KeyCode = 60
	KeySuperRight KeyCode = 61

	KeyMenu KeyCode = 62

	KeyPrintScreen KeyCode = 63
	KeyScrollLock  KeyCode = 64
	KeyPauseBreak  KeyCode = 65

	KeyInsert KeyCode = 66
	KeyHome   KeyCode = 67
	KeyDelete KeyCode = 68
	KeyEnd    KeyCode = 69

	KeyPageUp   KeyCode = 70
	KeyPageDown KeyCode = 71

	KeyArrowUp    KeyCode = 72
	KeyArrowRight KeyCode = 73
	KeyArrowDown  KeyCode = 74
	KeyArrowLeft  KeyCode = 75

	KeyNumLock    KeyCode = 76
	KeyNPDivide   KeyCode = 77
	KeyNPMultiply KeyCode = 78
	KeyNPMinus    KeyCode = 79
	KeyNPPlus     KeyCode = 80
	KeyNPReturn   KeyCode = 81

	KeyNP1       KeyCode = 82
	KeyNP2       KeyCode = 83
	KeyNP3       KeyCode = 84
	KeyNP4       KeyCode = 85
	KeyNP5       KeyCode = 86
	KeyNP6       KeyCode = 87
	KeyNP7       KeyCode = 88
	KeyNP8       KeyCode = 89
	KeyNP9       KeyCode = 90
	KeyNP0       KeyCode = 91
	KeyNPDecimal KeyCode = 92

	KeyF1  KeyCode = 93
	KeyF2  KeyCode = 94
	KeyF3  KeyCode = 95
	KeyF4  KeyCode = 96
	KeyF5  KeyCode = 97
	KeyF6  KeyCode = 98
	KeyF7  KeyCode = 99
	KeyF8  KeyCode = 100
	KeyF9  KeyCode = 101
	KeyF10 KeyCode = 102
	KeyF11 KeyCode = 103
	KeyF12 KeyCode = 104
)
