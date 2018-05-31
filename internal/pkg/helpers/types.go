package helpers

import (
	core "github.com/calle-gunnarsson/libdivecomputer-go/core"
)

type backendTable struct {
	Name  string
	Type  core.Family
	Model uint
}

var backends []backendTable = []backendTable{
	{"solution", core.DcFamilySuuntoSolution, 0},
	{"eon", core.DcFamilySuuntoEon, 0},
	{"vyper", core.DcFamilySuuntoVyper, 0x0A},
	{"vyper2", core.DcFamilySuuntoVyper2, 0x10},
	{"d9", core.DcFamilySuuntoD9, 0x0E},
	{"eonsteel", core.DcFamilySuuntoEonsteel, 0},
	{"aladin", core.DcFamilyUwatecAladin, 0x3F},
	{"memomouse", core.DcFamilyUwatecMemomouse, 0},
	{"smart", core.DcFamilyUwatecSmart, 0x10},
	{"meridian", core.DcFamilyUwatecMeridian, 0x20},
	{"sensus", core.DcFamilyReefnetSensus, 1},
	{"sensuspro", core.DcFamilyReefnetSensuspro, 2},
	{"sensusultra", core.DcFamilyReefnetSensusultra, 3},
	{"vtpro", core.DcFamilyOceanicVtpro, 0x4245},
	{"veo250", core.DcFamilyOceanicVeo250, 0x424C},
	{"atom2", core.DcFamilyOceanicAtom2, 0x4342},
	{"nemo", core.DcFamilyMaresNemo, 0},
	{"puck", core.DcFamilyMaresPuck, 7},
	{"darwin", core.DcFamilyMaresDarwin, 0},
	{"iconhd", core.DcFamilyMaresIconhd, 0x14},
	{"ostc", core.DcFamilyHwOstc, 0},
	{"frog", core.DcFamilyHwFrog, 0},
	{"ostc3", core.DcFamilyHwOstc3, 0x0A},
	{"edy", core.DcFamilyCressiEdy, 0x08},
	{"leonardo", core.DcFamilyCressiLeonardo, 1},
	{"n2ition3", core.DcFamilyZeagleN2ition3, 0},
	{"cobalt", core.DcFamilyAtomicsCobalt, 0},
	{"predator", core.DcFamilyShearwaterPredator, 2},
	{"petrel", core.DcFamilyShearwaterPetrel, 3},
	{"nitekq", core.DcFamilyDiveriteNitekq, 0},
	{"aqualand", core.DcFamilyCitizenAqualand, 0},
	{"idive", core.DcFamilyDivesystemIdive, 0x03},
	{"cochran", core.DcFamilyCochranCommander, 0},
}
