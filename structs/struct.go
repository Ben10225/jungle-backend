package structs

type Available struct {
	Date     string `json:"date,omitempty" validate:"required"`
	WorkTime []int  `json:"workTime,omitempty" validate:"required"`
	Yymm     string `json:"yymm,omitempty" validate:"required"`
}

type TimeTable struct {
	ThisMonth []Available
	NextMonth []Available
}
