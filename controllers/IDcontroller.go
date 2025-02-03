package controllers

import (
	"github.com/lokesh2201013/email-service/database"
	"github.com/lokesh2201013/email-service/models"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/gomail.v2"
)

func verifySMTP(sender models.Sender) error{

	d:=gomail.NewDialer(sender.SMTPHost,sender.SMTPPort,sender.Username,sender.AppPassword)

	_,err:=d.Dial()
	return err
}
//verification of email id
func VerifyEmailID(c *fiber.Ctx) error {
	var sender models.Sender

	// Parse the request body to get the sender details
	if err := c.BodyParser(&sender); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request content"})
	}

	// Check if the sender already exists in the database
	var existingSender models.Sender
	if err := database.DB.Where("email = ?", sender.Email).First(&existingSender).Error; err != nil {
		// If the sender does not exist, verify SMTP and create a new sender
		if err := verifySMTP(sender); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "SMTP verification failure"})
		}

		// Mark sender as verified and create a new entry in the database
		sender.Verified = true
		if err := database.DB.Create(&sender).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to add sender"})
		}
		return c.JSON(fiber.Map{"message": "Email verified and added to sender list", "email": sender.Email})
	} else {
		// If the sender exists and is not verified, verify SMTP and update the sender
		if existingSender.Verified {
			return c.Status(400).JSON(fiber.Map{"error": "Sender already verified"})
		}

		// Verify SMTP for the existing sender
		if err := verifySMTP(existingSender); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "SMTP verification failure for existing sender"})
		}

		// Update the existing sender record
		existingSender.Verified = true
		if err := database.DB.Save(&existingSender).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update sender"})
		}
		return c.JSON(fiber.Map{"message": "Email verified and updated", "email": existingSender.Email})
	}
}

//list all verified sender

func ListVerifiedIdentities(c *fiber.Ctx) error {
	var senders []models.Sender

	// Retrieve the admin username from the context
	adminUsername := c.Locals("admin").(string)

	// Filter by admin username
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

	// Retrieve the admin username from the context
	adminUsername := c.Locals("admin").(string)

	// Filter by admin username
	if err := database.DB.Where("verified = ? AND admin_name = ?", false, adminUsername).Find(&senders).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve senders"})
	}

	var emails []string
	for _, sender := range senders {
		emails = append(emails, sender.Email)
	}
	return c.JSON(fiber.Map{"idmail": emails})
}

//delete sender

func DeleteIdentity(c *fiber.Ctx) error {
	email := c.Params("email")

	if err := database.DB.Where("email = ?", email).Delete(&models.Sender{}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete sender"})
	}

	return c.JSON(fiber.Map{"message": "Sender deleted successfully"})
}

