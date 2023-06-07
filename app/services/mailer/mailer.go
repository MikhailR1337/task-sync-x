package mailer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/MikhailR1337/task-sync-x/app/application/forms"
	"github.com/sirupsen/logrus"
)

const addr = "http://mailer:3001/email"
const conType = "application/json"

func CheckedHomework(email string, name string, hwName string) {
	notification := forms.Mailer{}
	notification.Email = email
	notification.Subject = "Homework has checked"
	var body bytes.Buffer
	t, err := template.ParseFiles("public/template/email/checked.html")
	if err != nil {
		logrus.WithError(err)
	}
	fmt.Println(t)
	t.Execute(&body, struct {
		Name   string
		HwName string
	}{Name: name, HwName: hwName})
	notification.Template = body.String()
	JsonValue, err := json.Marshal(notification)
	if err != nil {
		logrus.WithError(err)
	}
	sendEmail(JsonValue)
}

func UpdatedHomework(email string, name string, hwName string, status string) {
	notification := forms.Mailer{}
	notification.Email = email
	notification.Subject = "Homework's status was changed"
	var body bytes.Buffer
	t, err := template.ParseFiles("public/template/email/updated.html")
	if err != nil {
		logrus.WithError(err)
	}
	fmt.Println(t)
	t.Execute(&body, struct {
		Name   string
		HwName string
		Status string
	}{Name: name, HwName: hwName, Status: status})
	notification.Template = body.String()
	JsonValue, err := json.Marshal(notification)
	if err != nil {
		logrus.WithError(err)
	}
	sendEmail(JsonValue)
}

func NewHomework(email string, name string, hwName string) {
	notification := forms.Mailer{}
	notification.Email = email
	notification.Subject = "You have a new homework"
	var body bytes.Buffer
	t, err := template.ParseFiles("public/template/email/new.html")
	if err != nil {
		logrus.WithError(err)
	}
	fmt.Println(t)
	t.Execute(&body, struct {
		Name   string
		HwName string
	}{Name: name, HwName: hwName})
	notification.Template = body.String()
	JsonValue, err := json.Marshal(notification)
	if err != nil {
		logrus.WithError(err)
	}
	sendEmail(JsonValue)
}

func sendEmail(body []byte) {
	_, err := http.Post(
		addr,
		conType,
		bytes.NewBuffer(body),
	)
	if err != nil {
		logrus.WithError(err)
	}
}
