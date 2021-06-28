package mail_service

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"github.com/xhit/go-simple-mail/v2"
	"html/template"
	"time"
)

type Options struct {
	Host               string
	Port               int
	Username           string
	Password           string
	Encryption         string
	FromHeader         string
	TechSupportSubject string
	DeveloperEmails    []string
	ProjectName        string
}

type MailService struct {
	host               string
	port               int
	username           string
	password           string
	encryption         string
	fromHeader         string
	techSupportSubject string
	developerEmails    []string
	projectName        string
}

func New(options Options) *MailService {
	return &MailService{
		host:               options.Host,
		port:               options.Port,
		username:           options.Username,
		password:           options.Password,
		encryption:         options.Encryption,
		techSupportSubject: options.TechSupportSubject,
		developerEmails:    options.DeveloperEmails,
		projectName:        options.ProjectName,
		fromHeader:         options.FromHeader,
	}
}

func (m *MailService) SendRegistrationConfirmation(emailAddress string, secretUrl string) {
	data := struct {
		SecretUrl   string
		Year        int
		ProjectName string
	}{
		SecretUrl:   secretUrl,
		Year:        time.Now().Year(),
		ProjectName: m.projectName,
	}
	m.sendHtmlEmailFromTemplate("templates/email/registration_confirmation.html", data, emailAddress)
}

func (m *MailService) SendPasswordRequest(emailAddress string, secretUrl string) {
	data := struct {
		SecretUrl   string
		Year        int
		ProjectName string
	}{
		SecretUrl:   secretUrl,
		Year:        time.Now().Year(),
		ProjectName: m.projectName,
	}
	m.sendHtmlEmailFromTemplate("templates/email/password_request.html", data, emailAddress)
}

func (m *MailService) SendUserConfirmed(emailAddress string) {
	data := struct {
		SecretUrl   string
		Year        int
		ProjectName string
	}{
		Year:        time.Now().Year(),
		ProjectName: m.projectName,
	}
	m.sendHtmlEmailFromTemplate("templates/email/user_confirmed.html", data, emailAddress)
}

func (m *MailService) sendHtmlEmailFromTemplate(template string, data interface{}, address ...string) {
	client, email, err := m.newEmail(address...)
	if err != nil {
		return
	}

	body, err := m.getTemplate(template, data)
	if err != nil {
		return
	}

	m.sendHtmlEmail(client, email, body)
}

func (m *MailService) sendHtmlEmail(client *mail.SMTPClient, email *mail.Email, body string) {
	email.SetBody(mail.TextHTML, body)

	if email.Error != nil {
		logrus.Errorf("unable to create mail message (%s)", email.Error.Error())
		return
	}

	err := email.Send(client)
	if err != nil {
		logrus.Errorf("unable to send message to (%s)", err.Error())
	}
}

func (m *MailService) newEmail(addresses ...string) (*mail.SMTPClient, *mail.Email, error) {
	client, err := m.newConnection()
	if err != nil {
		logrus.Errorf("unable to connect to mail server (%s)", err.Error())
		return nil, nil, err
	}

	email := mail.NewMSG()
	email.SetFrom(m.fromHeader)
	email.SetSubject(m.techSupportSubject)
	email.AddTo(addresses...)

	return client, email, nil
}

func (m *MailService) newConnection() (*mail.SMTPClient, error) {
	server := mail.NewSMTPClient()
	server.Host = m.host
	server.Port = m.port
	server.Username = m.username
	server.Password = m.password

	return server.Connect()
}

func (m *MailService) getTemplate(templateFile string, data interface{}) (string, error) {
	tpl, err := template.ParseFiles(templateFile)

	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	tpl.Execute(&buffer, data)

	return buffer.String(), nil
}
