package controllers

import (
	"github.com/lokesh2201013/email-service/database"
	"github.com/lokesh2201013/email-service/models"
	"github.com/gofiber/fiber/v2"
	"gopkg.in/gomail.v2"
	"strings"
	"github.com/lokesh2201013/email-service/metrics"
)


type EmailRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
	Format  string   `json:"format"`
}

// Helper function to create the email message
func createEmailMessage(sender models.Sender, req *EmailRequest) (*gomail.Message, error) {
	mail := gomail.NewMessage()
	mail.SetHeader("From", sender.Email)
	mail.SetHeader("To", req.To...)
	mail.SetHeader("Subject", req.Subject)

	switch req.Format {
	case "html":
		mail.SetBody("text/html", req.Body)
	case "text":
		mail.SetBody("text/plain", req.Body)
	default:
		return nil, fiber.NewError(400, "Invalid format")
	}

	return mail, nil
}

// Handle the error and update analytics
func handleEmailError(err error, analytics *models.Analytics) error {
	errMsg := err.Error()
	if strings.Contains(errMsg, "550") || strings.Contains(errMsg, "551") || strings.Contains(errMsg, "552") || strings.Contains(errMsg, "553") {
		analytics.Bounced++

		/*if strings.Contains(errMsg, "550 ") {
			parts := strings.Split(errMsg, "550 ")
			if len(parts) > 1 {
				invalidEmail := strings.TrimSpace(parts[1])

				if analytics.InvalidEmails == nil {
					analytics.InvalidEmails = []string{}
				}

				analytics.InvalidEmails = append(analytics.InvalidEmails, invalidEmail)
			}
		}*/
	} else if strings.Contains(errMsg, "421") || strings.Contains(errMsg, "530") || strings.Contains(errMsg, "521") || strings.Contains(errMsg, "554") {
		analytics.Rejected++
	}

	// Calculate and update metrics
	metricsWrapper := &metrics.AnalyticsWrapper{*analytics}
	metricsWrapper.CalculateMetrics()

	// Save updated analytics (you can put this in your database logic)
	database.DB.Save(&analytics)
	return fiber.NewError(500, "Failed to send email: " + errMsg)
}

// Send email
func SendEmail(c *fiber.Ctx) error {
	var req EmailRequest
	
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request content"})
	}

	// Fetch sender
	var sender models.Sender
	if err := database.DB.Where("email=? AND verified=?", req.From, true).First(&sender).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Sender not found"})
	}

	// Create email
	mail, err := createEmailMessage(sender, &req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Setup SMTP dialer
	d := gomail.NewDialer(sender.SMTPHost, sender.SMTPPort, sender.Username, sender.AppPassword)

	// Send email and capture the error
	err = d.DialAndSend(mail)

	// Fetch Analytics
	var analytics models.Analytics
	if err := database.DB.Where("admin_name = ? AND sender_id = ?", sender.AdminName, sender.ID).First(&analytics).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Analytics record not found"})
	}

	analytics.TotalEmails++

	// Handle email sending error
	if err != nil {
		return handleEmailError(err, &analytics)
	}

	// If no error, email was delivered
	analytics.Delivered++

	// Calculate and update metrics
	metricsWrapper := &metrics.AnalyticsWrapper{analytics}
	metricsWrapper.CalculateMetrics()

	// Save updated analytics
	database.DB.Save(&analytics)

	return c.JSON(fiber.Map{"message": "Email sent successfully"})
}
