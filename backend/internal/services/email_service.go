package services

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"gopkg.in/mail.v2"
)

// EmailService handles sending emails
type EmailService struct {
	dialer    *mail.Dialer
	from      string
	to        string
	sanitizer *bluemonday.Policy
}

// NewEmailService creates a new email service
func NewEmailService() *EmailService {
	// Read configuration from environment variables
	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		smtpHost = "smtp.gmail.com" // Default to Gmail
	}

	smtpPort := 587 // Default SMTP port

	smtpUser := os.Getenv("SMTP_USER")
	if smtpUser == "" {
		smtpUser = "joshuabgudgel@gmail.com"
	}

	smtpPass := os.Getenv("SMTP_PASSWORD")
	if smtpPass == "" {
		log.Println("WARNING: SMTP password not set in environment variables")
	}

	toEmail := os.Getenv("NOTIFICATION_EMAIL")
	if toEmail == "" {
		toEmail = "joshuabgudgel@gmail.com"
	}

	// Create the dialer
	dialer := mail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

	// Configure TLS
	dialer.TLSConfig = &tls.Config{
		ServerName: smtpHost,
	}

	sanitizer := bluemonday.StrictPolicy()
	return &EmailService{
		dialer:    dialer,
		from:      smtpUser,
		to:        toEmail,
		sanitizer: sanitizer,
	}
}

// sanitizeInput sanitizes user input to prevent XSS attacks
func (s *EmailService) sanitizeInput(input string) string {
	return s.sanitizer.Sanitize(input)
}

// SendBookingConfirmation sends an email notification for a successful booking
func (s *EmailService) SendBookingConfirmation(bookingID int, name, date, time, location string, people int, pkg string) error {
	// Sanitize all user inputs
	name = s.sanitizeInput(name)
	date = s.sanitizeInput(date)
	time = s.sanitizeInput(time)
	location = s.sanitizeInput(location)
	pkg = s.sanitizeInput(pkg)

	m := mail.NewMessage()

	// Set headers
	m.SetHeader("From", fmt.Sprintf("Toasted Coffee Co Support <%s>", s.from))
	m.SetHeader("To", s.to)
	m.SetHeader("Subject", fmt.Sprintf("New Booking: %s on %s", name, date))

	// Set email body with HTML
	m.SetBody("text/html", fmt.Sprintf(`
        <h2>New Booking Received</h2>
        <p>A new booking has been created successfully.</p>
        <h3>Booking Details:</h3>
        <ul>
            <li><strong>Booking ID:</strong> %d</li>
            <li><strong>Client:</strong> %s</li>
            <li><strong>Date:</strong> %s</li>
            <li><strong>Time:</strong> %s</li>
            <li><strong>Location:</strong> %s</li>
            <li><strong>People:</strong> %d</li>
            <li><strong>Package:</strong> %s</li>
        </ul>
        <p>Please check the admin dashboard for complete details.</p>
    `, bookingID, name, date, time, location, people, pkg))

	// Send the email
	return s.dialer.DialAndSend(m)
}

// SendBookingFailureAlert sends an email notification for a failed booking attempt
func (s *EmailService) SendBookingFailureAlert(name, email, phone string, errorDetails string) error {
	// Sanitize all user inputs
	name = s.sanitizeInput(name)
	email = s.sanitizeInput(email)
	phone = s.sanitizeInput(phone)
	errorDetails = s.sanitizeInput(errorDetails)

	m := mail.NewMessage()

	// Set headers
	m.SetHeader("From", fmt.Sprintf("Toasted Coffee Co Support <%s>", s.from))
	m.SetHeader("To", s.to)
	m.SetHeader("Subject", "ALERT: Failed Booking Attempt")

	// Build contact info section
	var contactInfo string
	if email != "" {
		contactInfo += fmt.Sprintf("<li><strong>Email:</strong> %s</li>", email)
	}
	if phone != "" {
		contactInfo += fmt.Sprintf("<li><strong>Phone:</strong> %s</li>", phone)
	}

	// Set email body with HTML
	m.SetBody("text/html", fmt.Sprintf(`
        <h2>Failed Booking Attempt</h2>
        <p>A customer attempted to make a booking but encountered an error.</p>
        <h3>Customer Information:</h3>
        <ul>
            <li><strong>Name:</strong> %s</li>
            %s
        </ul>
        <h3>Error Details:</h3>
        <p style="color: red; background-color: #ffeeee; padding: 10px; border-left: 4px solid #cc0000;">
            %s
        </p>
        <p>You may want to contact the customer to resolve this issue.</p>
    `, name, contactInfo, errorDetails))

	// Send the email
	return s.dialer.DialAndSend(m)
}

// SendInquiry sends an email notification for customer inquiries or contact form submissions
func (s *EmailService) SendInquiry(name, email, phone, message string) error {
	// Sanitize all user inputs
	name = s.sanitizeInput(name)
	email = s.sanitizeInput(email)
	phone = s.sanitizeInput(phone)
	message = s.sanitizeInput(message)

	// Recover from panic
	defer func() {
		if r := recover(); r != nil {
			log.Printf("RECOVERED from email panic: %v", r)
		}
	}()

	m := mail.NewMessage()

	// Set headers
	m.SetHeader("From", fmt.Sprintf("Toasted Coffee Co Support <%s>", s.from))
	m.SetHeader("To", s.to)
	m.SetHeader("Subject", fmt.Sprintf("New Inquiry from %s", name))

	// Build contact info section
	var contactInfo string
	if email != "" {
		contactInfo += fmt.Sprintf("<li><strong>Email:</strong> %s</li>", email)
	}
	if phone != "" {
		contactInfo += fmt.Sprintf("<li><strong>Phone:</strong> %s</li>", phone)
	}

	// Set email body with HTML
	m.SetBody("text/html", fmt.Sprintf(`
        <h2>New Customer Inquiry</h2>
        <p>A customer has submitted an inquiry or contact form.</p>
        <h3>Customer Information:</h3>
        <ul>
            <li><strong>Name:</strong> %s</li>
            %s
        </ul>
        <h3>Message:</h3>
        <div style="background-color: #f9f9f9; padding: 15px; border-left: 4px solid #4a6f8a; margin: 10px 0;">
            %s
        </div>
        <p style="color: #666; font-style: italic; margin-top: 20px;">
            Sent on: %s
        </p>
    `, name, contactInfo, message, time.Now().Format("January 2, 2006 at 3:04 PM")))

	// Send the email
	return s.dialer.DialAndSend(m)
}
