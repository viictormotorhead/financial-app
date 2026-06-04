package entities

type Tag struct {
	ID          uint    `gorm:"primaryKey"`
	UserID      *string `gorm:"type:varchar(21);uniqueIndex:idx_tags_user_name"`
	Name        string  `gorm:"type:varchar(255);not null;uniqueIndex:idx_tags_user_name"`
	Description string  `gorm:"type:varchar(500)"`
}

func (Tag) TableName() string {
	return "tags"
}
