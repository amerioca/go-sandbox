package controllers

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/antoniopapa/go-admin/database"
	"github.com/antoniopapa/go-admin/models"
	"github.com/antoniopapa/go-admin/util"
	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
	// return c.SendString("sent")
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		fmt.Printf("Parse ERR: %v\n", err)
		return err
	}
	/* Check for email */
	fmt.Println(data)
	if data["email"] == "" {
		c.Status(203)
		return c.JSON(fiber.Map{
			"msg": "email required",
		})
	}
	/* Check for Duplicate Emails */
	var chkemail models.User
	if rows := database.DB.Debug().Raw("SELECT id FROM users WHERE email = ? LIMIT 1", data["email"]).
		Scan(&chkemail).RowsAffected; rows != 0 {
		c.Status(203)
		return c.JSON(fiber.Map{
			"msg": "email exists",
		})
	}

	// Make Sure you have a Role for Id 1 as admin
	if data["password"] != data["password_confirm"] {
		c.Status(203)
		return c.JSON(fiber.Map{
			"message": "passwords do not match",
		})
	}

	lat, _ := strconv.ParseFloat(data["lat"], 64)
	lng, _ := strconv.ParseFloat(data["lng"], 64)
	fmt.Printf("lat: %+v %+v lng: %+v type: %+v\n", lat, data["lat"], lng, reflect.TypeOf(lat))

	user := models.User{
		/* basic info */
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Email:     data["email"],
		RoleId:    7, // Registered User
		/* extended info */
		LastSeen:    time.Now(),
		MemberAt:    time.Now(),
		Experiences: "[]",
		Map:         "[]",
		Lat:         lat,
		Lng:         lng,
	}

	user.SetPassword(data["password"])
	// fmt.Printf("CreateUser: %+v \n", user)

	// os.Exit(12)
	result := database.DB.Create(&user)
	user.Telephone = "(21)" + user.SiteId
	database.DB.Updates(&user)
	if result.Error != nil {
		fmt.Printf("%v\n", result.Error)
	}

	return c.JSON(fiber.Map{
		"msg":  "success",
		"data": user,
	})
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}
	fmt.Println(data)
	var user models.User

	database.DB.Where("email = ?", data["email"]).Preload("Role.Permissions").First(&user)

	if user.ID == 0 {
		c.Status(203)
		return c.JSON(fiber.Map{
			"msg": "email not found",
		})
	}

	if err := user.ComparePassword(data["password"]); err != nil {
		c.Status(203)
		return c.JSON(fiber.Map{
			"msg": "incorrect password",
		})
	}

	token, err := util.GenerateJwt(strconv.Itoa(int(user.ID)))

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24 * 365),
		HTTPOnly: true,
		// Secure:   true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"msg":  "success",
		"user": user,
	})
}

func User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	id, _ := util.ParseJwt(cookie)

	var user models.User

	database.DB.Where("id = ?", id).First(&user)

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

func UpdateInfo(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	cookie := c.Cookies("jwt")

	id, _ := util.ParseJwt(cookie)

	userId, _ := strconv.Atoi(id)

	user := models.User{
		ID:        uint(userId),
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Email:     data["email"],
	}

	database.DB.Model(&user).Updates(user)

	return c.JSON(user)
}

func UpdatePassword(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	if data["password"] != data["password_confirm"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "passwords do not match",
		})
	}

	cookie := c.Cookies("jwt")

	id, _ := util.ParseJwt(cookie)

	userId, _ := strconv.Atoi(id)

	user := models.User{
		ID: uint(userId),
	}

	user.SetPassword(data["password"])

	database.DB.Model(&user).Updates(user)

	return c.JSON(user)
}
