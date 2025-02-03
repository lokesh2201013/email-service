# email-service
Email Service

Overview

The Email Service is a REST API built with Golang and Fiber, designed for managing email communications. It supports user authentication, email sending, template creation, and email analytics.

Features

User Registration and Authentication

Email Identity Verification

Sending Emails with Templates

Email Analytics

Secure API with Authentication Middleware

Technologies Used

Golang

Fiber (Web Framework)

GORM (ORM for PostgreSQL)

Docker & Docker Compose

bcrypt for Password Hashing

JWT for Authentication

Alpine Linux for Lightweight Containerization

Installation

Prerequisites

Golang 1.23.5 and above

Docker & Docker Compose

PostgreSQL Database

teps

Clone the repository:

git clone https://github.com/lokesh2201013/email-service.git
cd email-service

Build and run the Docker containers:

docker-compose up --build -d

Based on the provided code, here are some valid curl endpoints for testing the email service API:

Register User:

curl -X POST http://localhost:3000/register \
     -H "Content-Type: application/json" \
     -d '{"username":"testuser", "password":"securepassword"}'

Login:

curl -X POST http://localhost:3000/login \
     -H "Content-Type: application/json" \
     -d '{"username":"testuser", "password":"securepassword"}'

Verify Email Identity (requires authentication token):

curl -X POST http://localhost:3000/verify-email-identity \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     -d '{"admin_name":"admin", "email":"test@example.com", "smtp_host":"smtp.gmail.com", "smtp_port":587, "username":"youremail@gmail.com", "password":"your_app_password"}'

List Verified Identities (requires authentication):

curl -X GET http://localhost:3000/list-verified-identities \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"

Delete Identity (requires authentication):

curl -X DELETE http://localhost:3000/delete-identity/test@example.com \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"

Create Template:

curl -X POST http://localhost:3000/create-template \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     -d '{"name":"welcome_template", "subject":"Welcome", "body":"Welcome to our service", "format":"text"}'

Send Email (authenticated):

curl -X POST http://localhost:3000/send-email \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     -d '{
         "from":"your_verified_email@example.com", 
         "to":["recipient@example.com"], 
         "subject":"Test Email", 
         "body":"Hello, this is a test email", 
         "format":"text"
     }'

Get Email Metrics for a Specific Sender:

curl -X GET http://localhost:3000/email-metrics/your_sender_email@example.com \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"

Get Admin Email Metrics:

curl -X GET http://localhost:3000/admin-email-metrics/your_admin_username \
     -H "Authorization: Bearer YOUR_JWT_TOKEN" 