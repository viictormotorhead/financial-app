package scopes

import "gorm.io/gorm"

func UserID(db *gorm.DB, userID string) *gorm.DB {
	return db.Where("user_id = ?", userID)
}
