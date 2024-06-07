package user

type User struct {
	Id        int64  `gorm:"primary_key" json:"id"`
	Name      string `gorm:"not null" json:"name"`
	Email     string `gorm:"not null" json:"email"`
	Password  string `gorm:"not null" json:"-"`
	Role      int    `gorm:"not null" json:"role"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt int64  `gorm:"autoCreateTime" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

type Role int

const (
	Admin     Role = 0
	Moderator Role = 1
	Member    Role = 2
)

func (r Role) Value() int {
	return int(r)
}
