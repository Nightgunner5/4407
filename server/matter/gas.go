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

func (g Gas) SpecificHeat() float64 {
	switch g {
	case Oxygen, Nitrogen, NitrousOxide:
		return 20
	case CarbonDioxide:
		return 30
	case Plasma:
		return 200
	}
	return 5
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

	R          = 8.3145  // joules per mole Δkelvin (PV = nRT)
	ATM        = 101.325 // joules per liter (one atmosphere)
	TileVolume = 2500    // liters

	TileContentsTotal    = TileVolume * ATM / (RoomTemperature * R) // moles (2.5 cubic meters at standard pressure and room temperature)
	TileContentsOxygen   = 0.21 * TileContentsTotal
	TileContentsNitrogen = 0.79 * TileContentsTotal
)
