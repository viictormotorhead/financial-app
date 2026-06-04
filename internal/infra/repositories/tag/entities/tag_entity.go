package entities

type TagEntity struct {
	ID          uint
	UserID      *string
	Name        string
	Description string
}
