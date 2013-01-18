package matter

import (
	"strconv"
)

type Gas uint8

const (
	_ Gas = iota
	Oxygen
	Nitrogen
	Plasma
	CarbonDioxide
	NitrousOxide

	gasCount
)

// Specific heat capacity in joules per kelvin
func (g Gas) SpecificHeat() float64 {
	switch g {
	case Oxygen:
		return 29.4 // http://www.wolframalpha.com/input/?i=specific+heat+capacity+of+1+mol+of+oxygen
	case Nitrogen:
		return 29.1 // http://www.wolframalpha.com/input/?i=specific+heat+capacity+of+1+mol+of+nitrogen
	case NitrousOxide:
		return 38.6 // http://www.wolframalpha.com/input/?i=specific+heat+capacity+of+1+mol+of+nitrous+oxide
	case CarbonDioxide:
		return 37.1 // http://www.wolframalpha.com/input/?i=specific+heat+capacity+of+1+mol+of+carbon+dioxide
	case Plasma:
		return 2000
	}
	panic(g)
}

func (g Gas) GoString() string {
	switch g {
	case Oxygen:
		return "Oxygen"
	case Nitrogen:
		return "Nitrogen"
	case Plasma:
		return "Plasma"
	case CarbonDioxide:
		return "CarbonDioxide"
	case NitrousOxide:
		return "NitrousOxide"
	}
	return strconv.FormatUint(uint64(g), 10)
}

func (g Gas) String() string {
	switch g {
	case Oxygen:
		return "Oxygen"
	case Nitrogen:
		return "Nitrogen"
	case Plasma:
		return "Plasma"
	case CarbonDioxide:
		return "Carbon Dioxide"
	case NitrousOxide:
		return "Nitrous Oxide"
	}
	return "Unknown-" + strconv.FormatUint(uint64(g), 10)
}

const (
	TempSpace       = 2.7
	Temp0C          = 273.15 // The temperature 0 Celsius in Kelvin
	WaterFreezes    = Temp0C
	RoomTemperature = Temp0C + 20
	WaterBoils      = Temp0C + 100

	R          = 8.3145  // joules per mole kelvin (PV = nRT)
	ATM        = 101.325 // joules per liter (one atmosphere)
	TileVolume = 2500    // liters

	TileContentsTotal    = TileVolume * ATM / (RoomTemperature * R) // moles (2.5 cubic meters at standard pressure and room temperature)
	TileContentsOxygen   = 0.21 * TileContentsTotal
	TileContentsNitrogen = 0.79 * TileContentsTotal

	GasMoveFraction = 0.2
)
