package main

import (
	"github.com/fatih/color"
)

// Linker variables
var Commit string

// CLI flags

// build & run flags
var UseJson bool
var JsonPath string
var OutputPath string
var BeVerbose bool = false
var IsStandalone bool
var FrontendToUse string
var GPExec string

//Colours

var GreenBoldUnderline = color.New([]color.Attribute{color.FgGreen, color.Bold, color.Underline}...)
var Bold = color.New([]color.Attribute{color.Bold}...)
var Underline = color.New([]color.Attribute{color.FgWhite, color.Underline}...)
var Grey = color.New([]color.Attribute{color.FgHiBlack}...)

//Other

var PluginExt string
