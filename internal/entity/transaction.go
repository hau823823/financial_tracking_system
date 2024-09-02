package entity

import "time"

type Transaction struct {
	ID          string    `gorm:"primaryKey"`
	UserID      string    `gorm:"index"`
	Date        time.Time `gorm:"index"`
	Amount      float64
	Category    string
	Description string
	Source      string `gorm:"type:enum('MANUAL', 'BANK', 'CREDIT_CARD')"`
	Reconciled  bool
}
