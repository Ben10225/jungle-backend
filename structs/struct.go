package structs

type Available struct {
	Date          string `json:"date,omitempty" validate:"required"`
	SureTimeArray []bool `json:"sureTimeArray,omitempty" validate:"required"`
	Yymm          string `json:"yymm,omitempty" validate:"required"`
}

type TimeTable struct {
	ThisMonth []Available
	NextMonth []Available
}
