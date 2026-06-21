package model

import "time"

type Todo struct {
    ID          uint       `gorm:"primaryKey"`
    AccountID   uint       
    Title       string     
    Description string     
    IsCompleted bool       
    DueDate     *time.Time 
    CreatedAt   time.Time  
    UpdatedAt   time.Time  
}

func (Todo) TableName() string {
    return "todo" // 小文字単数形のテーブル名を明示
}
