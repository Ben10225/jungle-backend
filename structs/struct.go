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
