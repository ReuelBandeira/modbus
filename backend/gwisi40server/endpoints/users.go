package endpoints

import (
	"os"
	"regexp"
	"strconv"
	"time"

	"gwisi40server/database"
	"gwisi40server/globals"
	"gwisi40server/models"
	"gwisi40server/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// Login @Summary		Login
// @Description	Login
// @Tags			Login
// @Accept			json
// @Produce		json
// @Success		200		{object}	models.Token
// @Failure		400		Invalid		request
// @Failure		401		Invalid		password
// @Failure		404		User		not				registered
// @Param			type	body		models.Login	true	"Login credentials"
// @Router			/api/v1/login [post]
func Login(c *fiber.Ctx) error {

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return fiber.ErrBadRequest
	}

	var user models.User

	database.DB.Where("email = ?", data["email"]).First(&user)

	if user.Id == 0 {
		return fiber.ErrNotFound
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		return fiber.ErrUnauthorized
	}

	type customClaims struct {
		Userid string `json:"user"`
		jwt.StandardClaims
	}

	expiration := time.Duration(globals.Conf.JWT.Expiration) * time.Second

	tok := customClaims{
		Userid: strconv.Itoa(int(user.Id)),
		StandardClaims: jwt.StandardClaims{
			Issuer:    strconv.Itoa(int(user.Id)),
			ExpiresAt: time.Now().Add(expiration).Unix(),
		},
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, tok)

	token, err := claims.SignedString([]byte(os.Getenv("API_SECRET")))
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(fiber.Map{
		"id":       user.Id,
		"email":    user.Email,
		"name":     user.Name,
		"language": user.Language,
		"token":    token,
	})
}

// CreateUser @Summary		Add user
// @Description	Add user
// @Tags			Login
// @Accept			json
// @Produce		json
//
// @Param			Authorization	header		string	true	"Authentication header"
//
// @Success		200				{object}	models.User
// @Failure		400				"Invalid	request"
// @Failure		404				"User not registered"
// @Failure		409				"User already	registered"
// @Param			type			body		models.User	true	"User data"
// @Router			/api/v1/register [post]
func CreateUser(c *fiber.Ctx) error {

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return fiber.ErrBadRequest
	}

	if data["name"] == "" || data["email"] == "" || data["password"] == "" {
		return fiber.ErrBadRequest
	}

	user := models.User{}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if !emailRegex.MatchString(data["email"]) {
		return fiber.ErrBadRequest
	}

	database.DB.Where("email = ?", data["email"]).First(&user)

	if user.Id != 0 {
		return fiber.ErrConflict
	}

	passwd, _ := utils.HashPassword(data["password"])

	user = models.User{
		Name:     data["name"],
		Email:    data["email"],
		Language: data["language"],
		Password: passwd,
	}

	database.DB.Create(&user)

	if user.Id == 0 {
		return fiber.ErrNotAcceptable
	}

	return c.JSON(user)
}

// UpdateUser @Summary		Update user password
// @Description	Update user password
// @Tags			Login
// @Accept			json
// @Produce		json
//
// @Param			Authorization	header		string	true	"Authentication header"
//
// @Success		200				{object}	models.User
// @Failure		400				Invalid		request
// @Failure		404				User		not				registered
// @Param			type			body		models.User	true	"User data"
// @Router			/api/v1/user/{email} [put]
func UpdateUser(c *fiber.Ctx) error {

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return fiber.ErrBadRequest
	}

	if data["password"] == "" {
		return fiber.ErrBadRequest
	}

	passwd, _ := utils.HashPassword(data["password"])

	user := models.User{}

	database.DB.Where("email = ?", c.Params("id")).First(&user)

	result := database.DB.Model(&user).Where("email = ?", c.Params("id")).Update("password", passwd)

	if result.RowsAffected == 0 {
		return fiber.ErrNotFound
	}

	return c.JSON(user)
}

// UpdateUserData @Summary		Update user data (name, email and password)
// @Description	Update user data (name, email and password)
// @Tags			Login
// @Accept			json
// @Produce		json
//
// @Param			Authorization	header		string	true	"Authentication header"
//
// @Success		200				{object}	models.User
// @Failure		400				Invalid		request
// @Failure		404				User		not				registered
// @Param			type			body		models.User	true	"User data"
// @Router			/api/v1/userdata/{id} [patch]
func UpdateUserData(c *fiber.Ctx) error {

	var data map[string]string
	var passwd []byte

	if err := c.BodyParser(&data); err != nil {
		return fiber.ErrBadRequest
	}

	user := models.User{}
	language := models.Language{}

	result := database.DB.Where("id = ?", c.Params("id")).First(&user)

	if result.RowsAffected == 0 {
		return fiber.ErrNotFound
	}

	if data["password"] != "" {
		passwd, _ = utils.HashPassword(data["password"])
		database.DB.Model(&user).Where("id = ?", c.Params("id")).Update("password", passwd)
	}

	if data["name"] != "" {
		database.DB.Model(&user).Where("id = ?", c.Params("id")).Update("name", data["name"])
	}

	if data["email"] != "" {
		database.DB.Model(&user).Where("id = ?", c.Params("id")).Update("email", data["email"])
	}

	if data["language"] != "" {
		database.DB.Model(&user).Where("id = ?", c.Params("id")).Update("language", data["language"])
		database.DB.Model(&language).Where("id = ?", 1).Update("language", data["language"])
	}

	return c.JSON(user)
}
