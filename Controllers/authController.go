package controllers

import (
	middlewears "clores-local/MiddleWears"
	models "clores-local/Models"
	database "clores-local/database"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {

	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}
	if data["password"] != data["password_confirm"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"Message": "passwords do not match",
		})
	}

	user := models.User{
		FirstName:    data["first_name"],
		LastName:     data["last_name"],
		Email:        data["email"],
		IsAmbassador: strings.Contains(c.Path(), "/api/ambassador"),
	}
	user.SetPassword(data["password"])

	database.DB.Create(&user)
	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}
	var user models.User
	database.DB.Where("email =?", data["email"]).First(&user)
	if user.Id == 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Username is incorrect",
		})
	}
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "password is incorrect",
		})
	}

	isAmbassador := strings.Contains(c.Path(), "/api/ambassador")

	var scope string

	if isAmbassador {
		scope = "ambassador"
	} else {
		scope = "admin"
	}

	token, err := middlewears.GenerateJtw(user.Id, scope)

	if err != nil {
		c.Status(fiber.StatusBadRequest)

		return c.JSON(fiber.Map{
			"message": "Username or password is incorrect",
		})
	}
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 2),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})

}

func GetUser(c *fiber.Ctx) error {
	id, _ := middlewears.GetUserID(c)
	var user models.User
	database.DB.Where("id = ?", id).First(&user)
	if strings.Contains(c.Path(), "api/ambassador"){
		ambassador:= models.Ambassador(user)
		ambassador.CalculateRevenue(database.DB)
		return c.JSON(ambassador)
	}
	return c.JSON(user)

}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func UpdateUser(c *fiber.Ctx) error {

	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}
	id, _ := middlewears.GetUserID(c)
	user := models.User{
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Email:     data["email"],
	}
	user.Id = id
	database.DB.Model(&user).Updates(&user)
	return c.JSON(&user)
}

func UpdatePassword(c *fiber.Ctx) error {

	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}
	if data["password"] != data["password_confirm"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"Message": "passwords do not match",
		})
	}
	id, _ := middlewears.GetUserID(c)
	user := models.User{}
	user.Id = id
	user.SetPassword(data["password"])
	database.DB.Model(&user).Updates(&user)
	return c.JSON(user)
}
