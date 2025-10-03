package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Username string `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email    string `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password string `gorm:"type:varchar(255);not null" json:"-"`
	Age      int    `gorm:"check:age > 0 AND age < 150" json:"age"`
	Status   string `gorm:"type:varchar(20);default:'active'" json:"status"`

	// 关联关系
	Orders []Order `gorm:"foreignKey:UserID" json:"orders,omitempty"`
	Posts  []Post  `gorm:"foreignKey:UserID" json:"posts,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate 创建前钩子
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	// 可以在这里添加创建前的逻辑
	// 例如：设置默认值、加密密码等
	return
}

// BeforeUpdate 更新前钩子
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	// 可以在这里添加更新前的逻辑
	return
}

// Order 订单模型
type Order struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	UserID  uint    `gorm:"not null" json:"user_id"`
	Amount  float64 `gorm:"type:decimal(10,2);not null" json:"amount"`
	Status  string  `gorm:"type:varchar(20);default:'pending'" json:"status"`
	Product string  `gorm:"type:varchar(100)" json:"product"`

	// 关联关系
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (Order) TableName() string {
	return "orders"
}

// Post 文章模型
type Post struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	UserID  uint   `gorm:"not null" json:"user_id"`
	Title   string `gorm:"type:varchar(200);not null" json:"title"`
	Content string `gorm:"type:text" json:"content"`

	// 关联关系
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (Post) TableName() string {
	return "posts"
}

// UserProfile 用户资料扩展模型
type UserProfile struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	UserID uint   `gorm:"uniqueIndex;not null" json:"user_id"`
	Bio    string `gorm:"type:text" json:"bio"`
	Avatar string `gorm:"type:varchar(255)" json:"avatar"`

	// 关联关系
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (UserProfile) TableName() string {
	return "user_profiles"
}
