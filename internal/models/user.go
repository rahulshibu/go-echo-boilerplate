package models

import (
	"time"

	"github.com/labstack/gommon/log"
	"github.com/rahulshibu/go-echo-boilerplate/internal/database"
	"golang.org/x/crypto/bcrypt"
)

// User - user basic
type User struct {
	ID               uint       `gorm:"primary_key" json:"id" `
	FullName         string     `gorm:"type:varchar(40);" json:"full_name" valid:"required,length(3|100)"`
	Email            string     `gorm:"type:varchar(40);unique_index; not null" json:"email" valid:"email,required"`
	JobTitle         string     `json:"job_title" valid:"required,length(1|100)"`
	PasswordHash     string     `json:"-"`
	Status           Status     `gorm:"foreignkey:id;status_foreignkey:name" json:"status"`
	AccountStatusID  int        `json:"account_status_id,omitempty" valid:"required,numeric"`
	IsAdmin          bool       `json:"is_admin" valid:"required"`
	CreatedBy        int        `json:"created_by"`
	PasswordRestHash string     `json:"-"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `sql:"index" json:"deleted_at,omitempty"`
	Password         string     `gorm:"-" json:"password,omitempty" valid:"required,length(8|100)"`
}

//CreateUser - create user model
func (um *User) CreateUser(u *User) (*User, error) {

	//database connection
	dbconn := database.GetSharedConnection()

	tx := dbconn.Begin()
	if err := tx.Create(u).Error; err != nil {

		//log error and return
		log.Error(u)
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	//get the user with status
	nUser, err := um.GetUserByID(u.ID)

	if err != nil {
		return nil, err
	}

	return &nUser, nil
}

//GetUserByID - get user by id
func (um *User) GetUserByID(userID uint) (User, error) {

	var (
		user User
		//status *Status
	)

	//get user by id
	db := database.GetSharedConnection()

	row := db.Debug().Select(`
			users.id,
			users.full_name,
			users.email,
			users.job_title,
			users.is_admin,
			users.created_by,
			users.created_at,
			users.updated_at,
			status.id,
			status.status_name,
			status.status_description
	`).Table("users").
		Where(`users.id = ?`, userID).
		Joins(`JOIN status ON status.id = users.account_status_id`).
		Row()

	err := row.Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.JobTitle,
		&user.IsAdmin,
		&user.CreatedBy,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Status.ID,
		&user.Status.StatusName,
		&user.Status.StatusDescription,
	)

	if err != nil {
		return user, err
	}

	return user, nil
}

//EditUser - Edit User Information
func (um *User) EditUser(userID int, u *User) (*User, error) {

	//database connection
	dbconn := database.GetSharedConnection()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.PasswordHash = string(hashedPassword)

	err := dbconn.Debug().Model(&u).Where("id = ? ", userID).Updates(
		map[string]interface{}{
			"full_name":         u.FullName,
			"email":             u.Email,
			"job_title":         u.JobTitle,
			"password_hash":     u.PasswordHash,
			"is_admin":          u.IsAdmin,
			"account_status_id": u.AccountStatusID,
		}).Error

	if err != nil {
		return nil, err
	}

	//get user
	user, err := um.GetUserByID(uint(userID))

	if err != nil {
		return nil, err
	}

	return &user, nil
}
