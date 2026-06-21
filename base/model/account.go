// model/account.go のイメージ
package model

import "time"
import "github.com/oapi-codegen/runtime/types"

type Account struct {
    ID        uint      `gorm:"primaryKey"`
    Username  string    `gorm:"not null;unique"`
    Email     types.Email `gorm:"not null;unique"`
    Password  string    `gorm:"not null"` // API側でハッシュ化（bcrypt等）した文字列を格納します
    CreatedAt time.Time
    UpdatedAt time.Time
    Todos     []Todo    `gorm:"foreignKey:AccountID"` // 1対多のリレーション
}

func (Account) TableName() string {
    return "account"
}
