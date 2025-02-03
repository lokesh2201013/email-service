Email Service 🚀
Overview
A robust REST API built with Golang and Fiber for comprehensive email communication management.
🌟 Features

User Authentication
Email Identity Verification
Template-based Email Sending
Advanced Email Analytics
Secure JWT Authentication
Rate-limited Email Sending

🛠 Technologies

Golang
Fiber Web Framework
GORM (PostgreSQL ORM)
Docker & Docker Compose
bcrypt
JWT Authentication

📋 Prerequisites

Golang 1.23.5+
Docker
Docker Compose
PostgreSQL

🚀 Quick Setup
1. Clone Repository
bashCopygit clone https://github.com/lokesh2201013/email-service.git
cd email-service
2. Launch Application
bashCopydocker-compose up --build -d
🔐 Authentication Endpoints
Register
bashCopycurl -X POST http://localhost:3000/register \
     -H "Content-Type: application/json" \
     -d '{"username":"testuser", "password":"securepassword"}'
Login
bashCopycurl -X POST http://localhost:3000/login \
     -H "Content-Type: application/json" \
     -d '{"username":"testuser", "password":"securepassword"}'
📧 Email Management
Verify Email Identity
bashCopycurl -X POST http://localhost:3000/verify-email-identity \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     -d '{
         "admin_name":"admin", 
         "email":"test@example.com", 
         "smtp_host":"smtp.gmail.com", 
         "smtp_port":587, 
         "username":"youremail@gmail.com", 
         "password":"your_app_password"
     }'
List Verified Emails
bashCopycurl -X GET http://localhost:3000/list-verified-identities \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"
Send Email
bashCopycurl -X POST http://localhost:3000/send-email \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     -d '{
         "from":"verified_email@example.com", 
         "to":["recipient@example.com"], 
         "subject":"Test Email", 
         "body":"Hello, this is a test email", 
         "format":"text"
     }'
📊 Metrics Endpoints
Sender Metrics
bashCopycurl -X GET http://localhost:3000/email-metrics/sender@example.com \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"
Admin Metrics
bashCopycurl -X GET http://localhost:3000/admin-email-metrics/admin_username \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"
🔒 Security Features

JWT Authentication
bcrypt Password Hashing
SMTP Verification
Email Sending Rate Limits

📝 Limitations

200 emails/day for new accounts
Sandbox period of 7 days for new users
