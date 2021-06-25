package data

import (
	"fmt"
)

type GetFilterData struct {
	Title string `json:"title"`
	To    string `json:"to"`
	From  string `json:"string"`
}

func NewGetFilterDataMessage(target string) GetFilterData {
	return GetFilterData{
		Title: "GET_FILTER_DATA",
		To:    target,
		From:  "USER",
	}
}

type UserData struct {
	Aquarium          string `json:"aqName"`
	Dst               int    `json:"dst"`
	FirmwareAvailable int    `json:"firmwareAvailable"`
	Host              string `json:"host"`
	Name              string `name:"name"`
	TankConfig        string `json:"tankConfig"`
	NetMode           string `json:"netmode"`
	Power             string `json:"power"`
	Unit              int    `json:"unit"`
	Version           int    `json:"version"`
	From              string `json:"from"`
	To                string `json:"to"`
}

type NetworkDevice struct {
	DHCP    int    `json:"dhcp"`
	Gateway []int  `json:"gateway"`
	IP      []int  `json:"ip"`
	SSID    string `json:"stSSID"`
	Subnet  []int  `json:"subaddress"`
	StPower string `json:"stPW"`
	From    string `json:"from"`
	To      string `json:"to"`
}

type AccessPoint struct {
	Power string `json:"apPW"`
	SSID  string `json:"apSSID"`
	From  string `json:"from"`
	To    string `json:"to"`
}
type PumpMode int

func (p PumpMode) String() string {
	switch p {
	case ConstantFlowMode:
		return "Constant Flow"
	case BioMode:
		return "Bio"
	case PulseFlowMode:
		return "Pulse"
	case ManualMode:
		return "Manual"
	default:
		return fmt.Sprintf("Unknown(%d)", p)
	}
}

var (
	// ConstantFlowMode indicates flow is at a constant rate.
	ConstantFlowMode PumpMode = 1
	// BioMode indicates the flow rate is different during the day and night.
	BioMode PumpMode = 4
	// PulseFlowMode indicates the flow fluctuates at a constant rate.
	PulseFlowMode PumpMode = 8
	// ManualMode indicates there is no flow control.
	ManualMode PumpMode = 16
)

type FilterData struct {
	DFS                int      `json:"dfs"`
	DFSFactor          int      `json:"dfsFaktor"`
	EndTimeNightMode   int      `json:"end_time_night_mode"`
	MinimumFrequency   int      `json:"minFreq"`
	Frequency          int      `json:"freq"`
	NightModeSollDay   int      `json:"nm_dfs_soll_day"`
	NightModeSollNight int      `json:"nm_dfs_soll_night"`
	PulseModeSollHigh  int      `json:"pm_dfs_soll_high"`
	PulseModeSollLow   int      `json:"pm_dfs_soll_low"`
	PulseModeTimeHigh  int      `json:"pm_time_high"`
	PulseModeTimeLow   int      `json:"pm_time_low"`
	PumpMode           PumpMode `json:"pumpMode"`
	RotationSpeed      int      `json:"rotSpeed"`
	RunTime            int      `json:"runTime"`
	ServiceHour        int      `json:"serviceHour"`
	SollStep           int      `json:"sollStep"`
	StartTimeNightMode int      `json:"start_time_night_mode"`
	TurnOffTime        int      `json:"turn_off_time"`
	Version            int      `json:"version"`
	From               string   `json:"from"`
	To                 string   `json:"to"`
}
