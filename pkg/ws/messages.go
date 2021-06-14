package ws

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

type FilterData struct {
	DFS                int    `json:"dfs"`
	DFSFactor          int    `json:"dfsFaktor"`
	EndTimeNightMode   int    `json:"end_time_night_mode"`
	MinimumFrequency   int    `json:"minFreq"`
	NightModeSollDay   int    `json:"nm_dfs_soll_day"`
	NightModeSollNight int    `json:"nm_dfs_soll_night"`
	PulseModeSollHigh  int    `json:"pm_dfs_soll_high"`
	PulseModeSollLow   int    `json:"pm_dfs_soll_low"`
	PulseModeTimeHigh  int    `json:"pm_time_high"`
	PulseModeTimeLow   int    `json:"pm_time_low"`
	PumpMode           int    `json:"pumpMode"`
	RotationSpeed      int    `json:"rotSpeed"`
	RunTime            int    `json:"runTime"`
	ServiceHour        int    `json:"serviceHour"`
	SollStep           int    `json:"sollStep"`
	StartTimeNightMode int    `json:"start_time_night_mode"`
	TurnOffTime        int    `json:"turn_off_time"`
	Version            int    `json:"version"`
	From               string `json:"from"`
	To                 string `json:"to"`
}
