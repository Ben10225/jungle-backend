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
	Titles []string       `json:"titles"`
	Detail BookingDetails `json:"detail"`
	Hour   BookingHour    `json:"hour"`
	Yymm   string         `json:"yymm"`
	Date   string         `json:"date"`
	User   BookingUser    `json:"user"`
}

type BookingUser struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type BookingHour struct {
	Index int `json:"index"`
	Whole int `json:"whole"`
}

type BookingDetails struct {
	Time  string `json:"time"`
	Cost  string `json:"cost"`
	State int    `json:"state"`
}

type Booking struct {
	ThisMonth         []ReserveData
	NextMonth         []ReserveData
	TheMonthAfterNext []ReserveData
}
