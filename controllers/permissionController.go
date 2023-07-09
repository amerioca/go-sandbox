package controllers

import (
	"github.com/antoniopapa/go-admin/database"
	"github.com/antoniopapa/go-admin/models"
	"github.com/gofiber/fiber/v2"
)

func AllPermissions(c *fiber.Ctx) error {
	var permissions []models.Permission

	database.DB.Find(&permissions)

	return c.JSON(permissions)
}
