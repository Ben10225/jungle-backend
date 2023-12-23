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
	Titles    []string
	Time      string
	Cost      string
	HourIndex int
	WholeHour int
}
