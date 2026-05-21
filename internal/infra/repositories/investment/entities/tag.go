package entities

type Tag struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"type:varchar(255);not null"`
	Description string `gorm:"type:varchar(500)"`
}

func (Tag) TableName() string {
	return "tags"
}
