package controllers

import (
	"github.com/lokesh2201013/email-service/database"
	"github.com/lokesh2201013/email-service/models"
	"github.com/gofiber/fiber/v2"
	"gopkg.in/gomail.v2"
)

// Send email
func SendEmail(c *fiber.Ctx) error {
	// Define the email request structure
	type EmailRequest struct {
		From    string   `json:"from"`
		To      []string `json:"to"`
		Subject string   `json:"subject"`
		Body    string   `json:"body"`
		Format  string   `json:"format"` // "text" or "html"
	}

	// Parse the incoming request body into the EmailRequest struct
	var req EmailRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request content"})
	}

	// Check if the sender email is verified
	var sender models.Sender
	if err := database.DB.Where("email=? AND verified=?", req.From, true).First(&sender).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Sender not found"})
	}

	// Create a new Gomail message
	mail := gomail.NewMessage()
	mail.SetHeader("From", sender.Email)
	mail.SetHeader("To", req.To...)
	mail.SetHeader("Subject", req.Subject)

	// Set the body based on the email format
	switch req.Format {
	case "html":
		mail.SetBody("text/html", req.Body) // Set body as HTML
	case "text":
		mail.SetBody("text/plain", req.Body) // Set body as plain text
	default:
		return c.Status(400).JSON(fiber.Map{"error": "Invalid format. Use 'html' or 'text'"})
	}

	// Set up the SMTP dialer with the sender's credentials
	d := gomail.NewDialer(sender.SMTPHost, sender.SMTPPort, sender.Username, sender.AppPassword)

	// Try to send the email
	if err := d.DialAndSend(mail); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to send email"})
	}

	// Respond with success message
	return c.JSON(fiber.Map{"message": "Email sent successfully"})
}
