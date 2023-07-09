package middlewares

import (
	"errors"
	"strconv"

	"github.com/antoniopapa/go-admin/database"
	"github.com/antoniopapa/go-admin/models"
	"github.com/antoniopapa/go-admin/util"
	"github.com/gofiber/fiber/v2"
)

func IsAuthorized(c *fiber.Ctx, page string) error {
	cookie := c.Cookies("jwt")

	Id, err := util.ParseJwt(cookie)

	// fmt.Println("Id", Id)

	if err != nil {
		return err
	}

	userId, _ := strconv.Atoi(Id)

	user := models.User{
		ID: uint(userId),
	}

	database.DB.Preload("Role").Find(&user)

	role := models.Role{
		ID: user.RoleId,
	}

	database.DB.Preload("Permissions").Find(&role)

	// fmt.Println("role.Perms", role.Permissions)
	if c.Method() == "GET" {
		for _, permission := range role.Permissions {
			if permission.Name == "view_"+page || permission.Name == "edit_"+page {
				return nil
			}
		}
	} else {
		for _, permission := range role.Permissions {
			// fmt.Println("perms.Name", permission.Name)
			if permission.Name == "edit_"+page {
				return nil
			}
		}
	}

	c.Status(fiber.StatusUnauthorized)
	// c.Status(409)
	return errors.New("unauthorizzed")
}
