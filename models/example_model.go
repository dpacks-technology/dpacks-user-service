package models

type ExampleModel struct {
	Column1 string `json:"column1"`
	Column2 string `json:"column2"`
	Column3 int    `json:"column3"`
}

type UpdateModel struct {
	Column1 string `json:"column1"`
	Column2 string `json:"column2"`
}
