package controllers

import (
	"github.com/lokesh2201013/email-service/database"
	"github.com/lokesh2201013/email-service/models"
	"github.com/gofiber/fiber/v2"
	"gopkg.in/gomail.v2"
)

func verifySMTP(sender models.Sender) error {
	d := gomail.NewDialer(sender.SMTPHost, sender.SMTPPort, sender.Username, sender.AppPassword)
	_, err := d.Dial()
	return err
}

func VerifyEmailID(c *fiber.Ctx) error {
	var sender models.Sender
	if err := c.BodyParser(&sender); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request content"})
	}

	var existingSender models.Sender
	if err := database.DB.Where("email = ?", sender.Email).First(&existingSender).Error; err != nil {
		if err := verifySMTP(sender); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "SMTP verification failure"})
		}
		sender.Verified = true
		if err := database.DB.Create(&sender).Error; err != nil {
			var Analytics models.Analytics
			if err := database.DB.Create(&Analytics).Error; err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Could not create analytics"})
			}
			return c.Status(500).JSON(fiber.Map{"error": "Failed to add sender"})
		}
		return c.JSON(fiber.Map{"message": "Email verified and added to sender list", "email": sender.Email})
	} else {
		if existingSender.Verified {
			return c.Status(400).JSON(fiber.Map{"error": "Sender already verified"})
		}
		if err := verifySMTP(existingSender); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "SMTP verification failure for existing sender"})
		}
		existingSender.Verified = true
		if err := database.DB.Save(&existingSender).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update sender"})
		}
		return c.JSON(fiber.Map{"message": "Email verified and updated", "email": existingSender.Email})
	}
}

func ListVerifiedIdentities(c *fiber.Ctx) error {
	var senders []models.Sender
	adminUsername := c.Locals("admin").(string)
	if err := database.DB.Where("verified = ? AND admin_name = ?", true, adminUsername).Find(&senders).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve senders"})
	}
	var emails []string
	for _, sender := range senders {
		emails = append(emails, sender.Email)
	}
	return c.JSON(fiber.Map{"idmail": emails})
}

func ListunVerifiedIdentities(c *fiber.Ctx) error {
	var senders []models.Sender
	adminUsername := c.Locals("admin").(string)
	if err := database.DB.Where("verified = ? AND admin_name = ?", false, adminUsername).Find(&senders).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve senders"})
	}
	var emails []string
	for _, sender := range senders {
		emails = append(emails, sender.Email)
	}
	return c.JSON(fiber.Map{"idmail": emails})
}

func DeleteIdentity(c *fiber.Ctx) error {
	email := c.Params("email")
	if err := database.DB.Where("email = ?", email).Delete(&models.Sender{}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete sender"})
	}
	return c.JSON(fiber.Map{"message": "Sender deleted successfully"})
}
