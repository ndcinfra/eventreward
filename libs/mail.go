package libs

import (
	"os"
	"strconv"

	"gopkg.in/mail.v2"
	gomail "gopkg.in/mail.v2"

	"github.com/astaxie/beego/logs"
	"github.com/joho/godotenv"

	"github.com/ndcinfra/eventreward/models"
)

type SendEmailInfo struct {
	Mid         int
	Email       string
	Displayname string
	Serial      string
	Title       string
	Body        string
}

// MakeEmail ...
func MakeEmail(er []models.EventRewards) {
	// logs.Info("test: ", e.Address)

	err := godotenv.Load()
	if err != nil {
		logs.Error("Error loading .env file")
	}

	SMTP := os.Getenv("SMTP")
	SMTPPORT, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	SMTPID := os.Getenv("SMTPID")
	SMTPPASS := os.Getenv("SMTPPASS")

	// for loop
	// TODO: 향후 channel로 변경 한다.
	// TODO: 포인터 array 이용
	var sei *SendEmailInfo
	sei = &SendEmailInfo{}

	// init mid
	var mid int
	mid = er[0].MID

	var ms *models.EmailMessage
	ms = &models.EmailMessage{}
	ms.ID = mid
	ms.GetMessage()

	for i := 0; i < len(er); i++ {
		if mid != er[i].MID {
			// get new message from db with new mid
			ms.ID = er[i].MID
			ms.GetMessage()
		}

		sei.Displayname = er[i].Displayname
		sei.Email = er[i].Email
		sei.Serial = er[i].Serial
		sei.Mid = er[i].MID

		makeMessage(sei, ms) // attach Eng

		//sendEmail(sei)

		m := gomail.NewMessage()
		m.SetHeader("From", "th@closerscs.com")
		m.SetHeader("To", sei.Email)
		m.SetHeader("Subject", sei.Title)
		m.SetBody("text/html", sei.Body)

		d := gomail.NewDialer(SMTP, SMTPPORT, SMTPID, SMTPPASS)
		d.StartTLSPolicy = mail.MandatoryStartTLS

		// Send the email to Bob, Cora and Dan.
		if err := d.DialAndSend(m); err != nil {
			logs.Error("send email error: ", err, sei.Email, sei.Displayname)
		} else {
			logs.Info("success send email", sei.Email, sei.Displayname)
		}
	}

}

func makeMessage(e *SendEmailInfo, ms *models.EmailMessage) {

	e.Title = ms.Title

	dear := ""
	preBody := ms.Body1
	postBody := "<br/><br/>" + ms.Body2 + e.Serial + "<br/><br/>" + ms.Body3 + ms.Body4 + ms.Body5

	e.Body = dear + preBody + postBody
}

func sendEmail(e *SendEmailInfo) {
	err := godotenv.Load()
	if err != nil {
		logs.Error("Error loading .env file")
	}

	SMTP := os.Getenv("SMTP")
	SMTPPORT, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	SMTPID := os.Getenv("SMTPID")
	SMTPPASS := os.Getenv("SMTPPASS")

	m := gomail.NewMessage()
	m.SetHeader("From", "th@closerscs.com")
	m.SetHeader("To", e.Email)
	m.SetHeader("Subject", e.Title)
	m.SetBody("text/html", e.Body)

	d := gomail.NewDialer(SMTP, SMTPPORT, SMTPID, SMTPPASS)
	d.StartTLSPolicy = mail.MandatoryStartTLS

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		logs.Error("send email error: ", err, e.Email, e.Displayname)
	} else {
		logs.Info("success send email")
	}

}
