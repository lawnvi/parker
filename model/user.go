package model

import (
	"github.com/jinzhu/gorm"
	"parker/config/database"
	"parker/config/log"
)

type User struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Field string `json:"field"`
}

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) UserDao {
	database.HasTable(User{})
	return UserDao{
		db:   db,
	}
}

func (m *UserDao) Append(user *User) {
	m.db.Create(user)
	log.I("db", "append user ok"+user.Name)
}

func (m *UserDao) FindByEmail(email string) User{
	var user = User{}
	m.db.Where("email= ?", email).Find(&user)
	return user
}

func (m *UserDao) Update(user *User) {
	m.db.Save(user)
}
