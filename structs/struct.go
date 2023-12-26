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
	Titles []string `json:"titles"`
	Detail struct {
		Time  string `json:"time"`
		Cost  string `json:"cost"`
		State int    `json:"state"`
	} `json:"detail"`
	Hour struct {
		Index int `json:"index"`
		Whole int `json:"whole"`
	} `json:"hour"`

	Yymm string `json:"yymm"`
	Date string `json:"date"`
	User struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
	} `json:"user"`
}

type Booking struct {
	ThisMonth         []ReserveData
	NextMonth         []ReserveData
	TheMonthAfterNext []ReserveData
}
