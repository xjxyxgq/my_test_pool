package services

import (
	"io"
	"log"
	"net/http"
	"time"

	"crypto/tls"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"gopkg.in/mail.v2"
)

type EmailService struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
}

type EmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}

func (s *EmailService) SendEmail(c *gin.Context) {
	var req EmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	m := mail.NewMessage()
	m.SetHeader("From", s.SMTPUser)
	m.SetHeader("To", req.To)
	m.SetHeader("Subject", req.Subject)
	m.SetBody("text/html", req.Content)

	d := mail.NewDialer(s.SMTPHost, s.SMTPPort, s.SMTPUser, s.SMTPPassword)
	d.SSL = true
	d.Timeout = 20 * time.Second // 增加超时时间

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Error sending email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	log.Println("Email sent successfully")
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (s *EmailService) SendEmailWithAttachment(to, subject, body string, attachment []byte) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.SMTPUser)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	m.Attach("screenshot.png", gomail.SetCopyFunc(func(w io.Writer) error {
		_, err := w.Write(attachment)
		return err
	}))

	d := gomail.NewDialer(s.SMTPHost, s.SMTPPort, s.SMTPUser, s.SMTPPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return d.DialAndSend(m)
}
