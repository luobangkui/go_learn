package main
import (
	"text/template"
	"bytes"
	"net/smtp"
	"strconv"
)
type EmailMessage struct {
	From, Subject, Body string
	To []string
}
type EmailCredentials struct {
	Username, Password, Server string
	Port int
}
const emailTemplate = `From: {{.From}}
To: {{.To}}
Subject {{.Subject}}
{{.Body}}
`
var t *template.Template

func init() {
	t = template.New("email")
	t.Parse(emailTemplate)
}
func main() {
	message := &EmailMessage{
		From: "1234555",
		To: []string{"1015957634@qq.com"},
		Subject: "A test",
		Body: "Just saying hi",
	}
	var body bytes.Buffer
	t.Execute(&body, message)
	authCreds := &EmailCredentials{
		Username: "luobangkui",
		Password: "xiaoluo1370",
		Server: "smtp.example.com",
		Port: 25,
	}
	auth := smtp.PlainAuth("",
		authCreds.Username,
		authCreds.Password,
		authCreds.Server,
	)
	smtp.SendMail(authCreds.Server+":"+strconv.Itoa(authCreds.Port),
		auth,
		message.From,
		message.To,
		body.Bytes())
}