package libs

import (
	"log"
	"os"
	"strconv"

	"gopkg.in/mail.v2"
	gomail "gopkg.in/mail.v2"

	"github.com/astaxie/beego/logs"
	"github.com/joho/godotenv"

	"github.com/ndcinfra/eventreward/models"
	"github.com/vanng822/go-premailer/premailer"
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
			// DB update
			err := models.UpdateEventRewardsDoneOne(er[i].ID)
			if err != nil {
				logs.Error("update done error", er[i].ID, sei.Email, sei.Displayname)
			}

			logs.Info("success send email", sei.Email, sei.Displayname)
		}
	}

}

// 우선 개발
// 향후 수정
func MakeEmailMarketing(er models.EventRewards) {
	inputFile := "./index.html"

	prem, err := premailer.NewPremailerFromFile(inputFile, premailer.NewOptions())
	if err != nil {
		log.Fatal(err)
	}

	html, err := prem.Transform()
	if err != nil {
		log.Fatal(err)
	}

	err = godotenv.Load()
	if err != nil {
		logs.Error("Error loading .env file")
	}

	SMTP := os.Getenv("SMTP")
	SMTPPORT, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	SMTPID := os.Getenv("SMTPID")
	SMTPPASS := os.Getenv("SMTPPASS")

	m := gomail.NewMessage()
	m.SetHeader("From", "th@closerscs.com")
	m.SetHeader("To", er.Email)
	m.SetHeader("Subject", "[Closers Thailand] อัพเดทสนามบินนานาชาติ และ ชุดว่ายน้ำสุดคูล")
	m.SetBody("text/html", html)

	d := gomail.NewDialer(SMTP, SMTPPORT, SMTPID, SMTPPASS)
	d.StartTLSPolicy = mail.MandatoryStartTLS

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		logs.Error("send email error: ", err)
	} else {
		// DB update
		logs.Info("success send email")
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
