package main

var attributes map[string]interface{}

// Instructions
const (
	SysIFieldQuery  uint8 = 0x00
	SysIDeviceQuery uint8 = 0x01
)

// Instruction args
const (
	FieldName          uint8 = 0x00
	FieldDisplayWidth  uint8 = 0x01
	FieldDisplayHeight uint8 = 0x02
	FieldMemory        uint8 = 0x03
)
