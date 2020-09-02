package controllers

import (
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"github.com/labstack/echo"
	"gitlab.com/accubits/mapclub_multitenancy_admin/internal/models"
	"gitlab.com/accubits/mapclub_multitenancy_admin/session"

	"github.com/asaskevich/govalidator"
)

//userHandler - controller for users
type userHandler struct{}

//registerUserHandler - register handlers for user api
func registerUserHandler(r *echo.Group) {

	u := userHandler{}

	//group user
	userRoute := r.Group("/user")
	//set api
	userRoute.POST("/create", u.createUser)
	userRoute.GET("/:id", u.getUser)
	userRoute.POST("/edit/:id", u.editUser)

}

//createUser - user and generate token, refresh token
func (h *userHandler) createUser(c echo.Context) error {

	var user models.User

	if err := c.Bind(&user); err != nil {
		return err
	}

	//validator
	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, session.InternalServerError(c, err))
	}

	//get context by value
	user.CreatedBy = c.Get("userId").(int)

	//hash the password)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.PasswordHash = string(hashedPassword)

	//create user
	newUser, err := user.CreateUser(&user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, session.InternalServerError(c, err))
	}

	//show user
	return c.JSON(200, newUser)
}

//getUser - get user information based on id
func (h *userHandler) getUser(c echo.Context) error {

	var (
		userModel models.User
	)

	id := c.Param("id")

	userID, _ := strconv.Atoi(id)
	user, err := userModel.GetUserByID(uint(userID))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, session.InternalServerError(c, err))
	}

	//show user
	return c.JSON(200, user)
}

func (h *userHandler) editUser(c echo.Context) error {

	//Editable Infromation
	//full_name, email.password,account_status,job_title

	var (
		userEdit *models.User
		err      error
	)

	id := c.Param("id")

	userID, _ := strconv.Atoi(id)

	//bind user
	if err = c.Bind(&userEdit); err != nil {
		return err
	}

	userEdit, err = userEdit.EditUser(userID, userEdit)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, session.InternalServerError(c, err))
	}

	return c.JSON(200, userEdit)
}
