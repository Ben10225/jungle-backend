package structs

type Available struct {
	Date     string `json:"date,omitempty" bson:"date"`
	WorkTime []int  `json:"workTime,omitempty" bson:"workTime"`
	Yymm     string `json:"yymm,omitempty" bson:"yymm"`
}

type TimeTable struct {
	ThisMonth         []Available
	NextMonth         []Available
	TheMonthAfterNext []Available
}

type ReviseAvailable struct {
	Yymm      string
	Date      string
	HourIndex int
	WholeHour int
}

type ReserveData struct {
	Titles    []string `json:"titles,omitempty"`
	Time      string   `json:"time,omitempty"`
	Cost      string   `json:"cost,omitempty"`
	HourIndex int      `json:"hourIndex,omitempty" bson:"hourIndex"`
	WholeHour int      `json:"wholeHour,omitempty" bson:"wholeHour"`
	Yymm      string   `json:"yymm,omitempty"`
	Date      string   `json:"date,omitempty"`
}

type Booking struct {
	ThisMonth         []ReserveData
	NextMonth         []ReserveData
	TheMonthAfterNext []ReserveData
}
