package command

import (
	"testing"
	"github.com/go-macaron/binding"
	"fmt"
	"gopkg.in/macaron.v1"
)

type Contact struct {
	Name           string `form:"name" binding:"Required"`
	Email          string `form:"email"`
	Message        string `form:"message" binding:"Required"`
	MailingAddress string `form:"mailing_address"`
}

func TestPost(t *testing.T) {

	m := macaron.Classic()

	m.Post("/contact/submit", binding.BindIgnErr(Contact{}), func(contact Contact) string {
		return fmt.Sprintf("Name: %s\nEmail: %s\nMessage: %s\nMailing Address: %v",
			contact.Name, contact.Email, contact.Message, contact.MailingAddress)
	})

	m.Run()

}
