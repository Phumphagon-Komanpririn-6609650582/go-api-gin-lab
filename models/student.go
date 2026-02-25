package models

type Student struct {
	Id    string  `json:"id" binding:"required"`
	Name  string  `json:"name" binding:"required"`
	Major string  `json:"major"`
	GPA   float64 `json:"gpa" binding:"min=0.00,max=4.00"`
}
