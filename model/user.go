package model

type User struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Age  int64  `json:"age" db:"age"`
}
