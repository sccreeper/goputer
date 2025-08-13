package main

import (
	"github.com/fatih/color"
)

// Linker variables
var Commit string

// CLI flags

// build & run flags
var useJson bool
var jsonPath string

var programCompileOut string
var beVerbose bool = false
var isStandalone bool

var frontendToUse string
var gpExec string
var useProfiler bool
var profilerOut string

//Colours

var greenBoldUnderline = color.New([]color.Attribute{color.FgGreen, color.Bold, color.Underline}...)
var bold = color.New([]color.Attribute{color.Bold}...)
var underline = color.New([]color.Attribute{color.FgWhite, color.Underline}...)
var grey = color.New([]color.Attribute{color.FgHiBlack}...)

//Other

var pluginExt string
