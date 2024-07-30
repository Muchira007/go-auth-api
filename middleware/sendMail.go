package middleware

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"

	"gopkg.in/gomail.v2"
)

func SendMailSimple(subject string, body string, to []string) {
	auth := smtp.PlainAuth(
		"",
		"stivmicah@gmail.com",
		"irsw zplm xzkz gymx",
		"smtp.gmail.com",
	)

	msg := "Subject:" + subject + "\n" + body

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"muchirasteve2@gmail.com",
		to,
		[]byte(msg),
	)

	if err != nil {
		fmt.Println("Error sending email:", err)
	} else {
		fmt.Println("Email sent successfully")
	}
}

// func SendMailSimpleHTML(subject string, templatePath string, to []string) {
// 	//Get html
// 	var body bytes.Buffer
// 	t, err := template.ParseFiles(templatePath)

// 	t.Execute(&body, struct{ Name string }{Name: "Steve"})

// 	if err != nil {
// 		fmt.Println("Error parsing template:", err)
// 		return
// 	}

// 	auth := smtp.PlainAuth(
// 		"",
// 		"stivmicah@gmail.com",
// 		"irsw zplm xzkz gymx",
// 		"smtp.gmail.com",
// 	)

// 	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"

// 	msg := "Subject:" + subject + "\n" + headers + "\n \n" + body.String()

// 	err = smtp.SendMail(
// 		"smtp.gmail.com:587",
// 		auth,
// 		"muchirasteve2@gmail.com",
// 		to,
// 		[]byte(msg),
// 	)

//		if err != nil {
//			fmt.Println("Error sending email:", err)
//		} else {
//			fmt.Println("Email sent successfully")
//		}
//	}
func SendGoMail(toEmail, subject, bodyText string) {
	// Prepare the email template
	var body bytes.Buffer
	t, err := template.New("email").Parse(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>{{.Subject}}</title>
		</head>
		<body>
			{{.BodyText}}
		</body>
		</html>
	`)

	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	err = t.Execute(&body, struct {
		Subject  string
		BodyText string
	}{
		Subject:  subject,
		BodyText: bodyText,
	})

	if err != nil {
		fmt.Println("Error executing template:", err)
		return
	}

	// Send with gomail
	m := gomail.NewMessage()
	m.SetHeader("From", "muchirasteve2@gmail.com")
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer("smtp.gmail.com", 587, "stivmicah@gmail.com", "irsw zplm xzkz gymx")

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
