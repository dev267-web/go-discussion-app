// mailer helper 
// pkg/mailer/mailer.go

package mailer

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"os"
	"strings"
)

// Config holds SMTP server configuration.
type Config struct {
	Host     string // SMTP server host (e.g. "smtp.gmail.com")
	Port     string // SMTP server port (e.g. "587")
	Username string // SMTP auth username (often the full email address)
	Password string // SMTP auth password (or app password)
	From     string // From email address (e.g. "noreply@example.com")
}

// loadConfig reads required environment variables into a Config struct.
// It panics if any required var is missing.
func loadConfig() *Config {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USERNAME")
	pass := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("FROM_EMAIL")

	missing := []string{}
	if host == "" {
		missing = append(missing, "SMTP_HOST")
	}
	if port == "" {
		missing = append(missing, "SMTP_PORT")
	}
	if user == "" {
		missing = append(missing, "SMTP_USERNAME")
	}
	if pass == "" {
		missing = append(missing, "SMTP_PASSWORD")
	}
	if from == "" {
		missing = append(missing, "FROM_EMAIL")
	}
	if len(missing) > 0 {
		panic(fmt.Sprintf("missing required environment variables: %s", strings.Join(missing, ", ")))
	}

	return &Config{
		Host:     host,
		Port:     port,
		Username: user,
		Password: pass,
		From:     from,
	}
}

// buildAuth returns an smtp.Auth object for PLAIN auth.
func (cfg *Config) buildAuth() smtp.Auth {
	// Use PLAIN authentication (most SMTP servers on port 587 support this).
	return smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
}

// dialTLS creates a TLS connection to the SMTP server (for port 465).
// Note: If you use port 587 with STARTTLS, use smtp.Dial + StartTLS instead.
func (cfg *Config) dialTLS() (*smtp.Client, error) {
	addr := net.JoinHostPort(cfg.Host, cfg.Port)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         cfg.Host,
	}
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to dial TLS: %w", err)
	}

	client, err := smtp.NewClient(conn, cfg.Host)
	if err != nil {
		return nil, fmt.Errorf("failed to create SMTP client: %w", err)
	}
	return client, nil
}

// SendMail sends a plaintext email to one or more recipients.
// - to: slice of recipient email addresses.
// - subject: email subject.
// - body: plaintext body (no HTML).
func SendMail(to []string, subject, body string) error {
	cfg := loadConfig()
	
	// 1) Build the email headers
	headers := make(map[string]string)
	headers["From"] = cfg.From
	headers["To"] = strings.Join(to, ", ")
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/plain; charset=\"utf-8\""
	
	// 2) Concatenate headers and body
	var msgBuilder strings.Builder
	for k, v := range headers {
		fmt.Fprintf(&msgBuilder, "%s: %s\r\n", k, v)
	}
	msgBuilder.WriteString("\r\n" + body)

	// 3) Connect to SMTP server
	addr := net.JoinHostPort(cfg.Host, cfg.Port)

	// Use STARTTLS on port 587 (typical). If your provider requires port 465, swap to dialTLS().
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("smtp dial error: %w", err)
	}
	defer client.Quit()

	// 4) StartTLS (if needed)
	if ok, _ := client.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         cfg.Host,
		}
		if err := client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("failed to start TLS: %w", err)
		}
	}

	// 5) Authenticate
	if err := client.Auth(cfg.buildAuth()); err != nil {
		return fmt.Errorf("smtp auth error: %w", err)
	}

	// 6) Set the sender and recipients
	if err := client.Mail(cfg.From); err != nil {
		return fmt.Errorf("failed to set MAIL FROM: %w", err)
	}
	for _, recipient := range to {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to add RCPT TO %s: %w", recipient, err)
		}
	}

	// 7) Write the message data
	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get Data writer: %w", err)
	}
	_, err = wc.Write([]byte(msgBuilder.String()))
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("failed to close Data writer: %w", err)
	}

	return nil
}

// SendMailHTML sends an HTML email to one or more recipients.
// - to: slice of recipient email addresses.
// - subject: email subject.
// - htmlBody: HTML content; headers will be set accordingly.
func SendMailHTML(to []string, subject, htmlBody string) error {
	cfg := loadConfig()
	
	// 1) Build email headers
	headers := make(map[string]string)
	headers["From"] = cfg.From
	headers["To"] = strings.Join(to, ", ")
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"utf-8\""
	
	// 2) Concatenate headers and HTML body
	var msgBuilder strings.Builder
	for k, v := range headers {
		fmt.Fprintf(&msgBuilder, "%s: %s\r\n", k, v)
	}
	msgBuilder.WriteString("\r\n" + htmlBody)

	// 3) Connect to SMTP server
	addr := net.JoinHostPort(cfg.Host, cfg.Port)
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("smtp dial error: %w", err)
	}
	defer client.Quit()

	// 4) StartTLS if available
	if ok, _ := client.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         cfg.Host,
		}
		if err := client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("failed to start TLS: %w", err)
		}
	}

	// 5) Authenticate
	if err := client.Auth(cfg.buildAuth()); err != nil {
		return fmt.Errorf("smtp auth error: %w", err)
	}

	// 6) Set MAIL FROM and RCPT TO
	if err := client.Mail(cfg.From); err != nil {
		return fmt.Errorf("failed to set MAIL FROM: %w", err)
	}
	for _, recipient := range to {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to add RCPT TO %s: %w", recipient, err)
		}
	}

	// 7) Write the HTML message data
	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get Data writer: %w", err)
	}
	_, err = wc.Write([]byte(msgBuilder.String()))
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("failed to close Data writer: %w", err)
	}

	return nil
}
