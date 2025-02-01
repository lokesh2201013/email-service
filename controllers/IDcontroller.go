package controllers

import (
	"github.com/lokesh2201013/email-service/database"
	"github.com/lokesh2201013/email-service/models"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/gomail.v2"
)

func verifySTMP(sender models.Sender) error{

	d:=gomail.NewDialer(sender.SMTPHost,sender.SMTPPort,sender.Username,sender.Password)

	_,err:=d.Dial()
	return err
}
//verification of email id
func VerifyEmailID(c *fiber.Ctx) error{
	var sender models.Sender

	if err:=c.BodyParser(&sender); err!=nil{
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request content"})
	}
    
	var existingSender models.Sender

	if err:=database.DB.Where("email = ?", sender.Email).First(&existingSender).Error; err!=nil{
		return c.Status(400).JSON(fiber.Map{"error": "Sender not found"})
	}

	if err:= verifySTMP(sender);err!=nil{
		return c.Status(400).JSON(fiber.Map{"error": "SMTP verification failure"})
	}
     
    sender.Verified = true

	if err := database.DB.Create(&sender).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to add sender"})
	}

	return c.JSON(fiber.Map{"message": "Email verified and added to sender list", "email": sender.Email})
}

//list all verified sender

func ListVerifiedIdentities(c *fiber.Ctx) error {
	var senders []models.Sender

	database.DB.Where("verified = ?", true).Find(&senders)
     
	var emails []string

	for _,sender:=range senders{
		emails=append(emails,sender.Email)
	}
	return c.JSON(fiber.Map{"idmail": emails})
}

//list all unverified sender

func ListunVerifiedIdentities(c *fiber.Ctx) error {
	var senders []models.Sender

	database.DB.Where("verified = ?", false).Find(&senders)
     
	var emails []string

	for _,sender:=range senders{
		emails=append(emails,sender.Email)
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

