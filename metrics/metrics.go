package metrics

import (
	"github.com/lokesh2201013/email-service/models"
	"github.com/gofiber/fiber/v2"
	"github.com/lokesh2201013/email-service/database"
)

// analyticswrapper wraps the analytics model
type AnalyticsWrapper struct {
	models.Analytics
}


func (a *AnalyticsWrapper) CalculateMetrics() {
	a.DeliveryRate = float64(a.Delivered) / float64(a.TotalEmails) * 100
	a.BounceRate = float64(a.Bounced) / float64(a.TotalEmails) * 100
	a.ComplaintRate = float64(a.Complaints) / float64(a.TotalEmails) * 100
	a.RejectRate = float64(a.Rejected) / float64(a.TotalEmails) * 100
}

func GetEmailMetrics(c *fiber.Ctx) error {
	// Get the sender email from the URL parameter
	senderEmail := c.Params("senderEmail")

	// Fetch the sender by email
	var sender models.Sender
	if err := database.DB.Where("email = ?", senderEmail).First(&sender).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Sender not found"})
	}

	// Fetch the associated analytics record using the sender's ID
	var analytics models.Analytics
	if err := database.DB.Where("sender_id = ?", sender.ID).First(&analytics).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Analytics not found for sender"})
	}

	// Create metrics wrapper
	metricsWrapper := &AnalyticsWrapper{
		Analytics: analytics,
	}

	// Calculate the metrics
	metricsWrapper.CalculateMetrics()

	// Return the metrics data
	return c.JSON(fiber.Map{
		"delivery_rate":   metricsWrapper.DeliveryRate,
		"bounce_rate":     metricsWrapper.BounceRate,
		"complaint_rate":  metricsWrapper.ComplaintRate,
		"reject_rate":     metricsWrapper.RejectRate,
		"total_emails":    analytics.TotalEmails,
		"delivered":       analytics.Delivered,
		"bounced":         analytics.Bounced,
		"complaints":      analytics.Complaints,
		"rejected":        analytics.Rejected,
	})
}

func GetAdminEmailMetrics(c *fiber.Ctx) error {
	// Get the admin name from the URL parameter
	adminName := c.Params("adminName")

	// Fetch all senders associated with the admin
	var senders []models.Sender
	if err := database.DB.Where("admin_name = ?", adminName).Find(&senders).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "No senders found for the specified admin"})
	}

	// Initialize a slice to hold metrics for all senders
	var allSenderMetrics []fiber.Map

	// Iterate over all senders and fetch their analytics and calculate metrics
	for _, sender := range senders {
		// Fetch the associated analytics for the sender
		var analytics models.Analytics
		if err := database.DB.Where("sender_id = ?", sender.ID).First(&analytics).Error; err != nil {
			continue // If no analytics found, skip this sender
		}

		// Create the analytics wrapper and calculate metrics
		metricsWrapper := AnalyticsWrapper{
			Analytics: analytics,
		}
		metricsWrapper.CalculateMetrics()

		// Append the metrics for this sender
		allSenderMetrics = append(allSenderMetrics, fiber.Map{
			"sender_email":   sender.Email,
			"delivery_rate":  metricsWrapper.DeliveryRate,
			"bounce_rate":    metricsWrapper.BounceRate,
			"complaint_rate": metricsWrapper.ComplaintRate,
			"reject_rate":    metricsWrapper.RejectRate,
			"total_emails":   analytics.TotalEmails,
			"delivered":      analytics.Delivered,
			"bounced":        analytics.Bounced,
			"complaints":     analytics.Complaints,
			"rejected":       analytics.Rejected,
		})
	}

	// return admin name and all sender metrics in form of array of type fiber.Map
	return c.JSON(fiber.Map{
		"admin": adminName,
		"senders_metrics": allSenderMetrics,
	})
}

