package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/joho/godotenv"
	"github.com/ndcinfra/eventreward/libs"
	"github.com/ndcinfra/eventreward/models"
	"github.com/vanng822/go-premailer/premailer"

	// "gopkg.in/gomail.v2"
	"gopkg.in/mail.v2"
	gomail "gopkg.in/mail.v2"

	"github.com/astaxie/beego/orm"
)

// 프로그램의 기능
// event_reward_type
// 1: send email with cououpon serial
//		insert coupon
//		insert coupon_send_history and get serial
//		send email
// 2:
//
//

// 우선 for loop를 사용 해서 개발
// TODO: 향후 channel을 이용한다. !!!
// TODO: 포인터 array 이용 !!!

// GetSendEmail ...
func GetSendEmail() {
	// Step 1
	// 대상자 가져오기
	var eventRewards []models.EventRewards
	eventRewards, err := models.GetSendEmailReward()
	if err != nil {
		logs.Error("Error GetSendEmailReward: ", err)
		return
	}

	logs.Info("Success GetSendEmail: ", len(eventRewards))

	// event_reward_type에 따라 분기
	// event_reward_type:
	//		1 : 대상자에 쿠폰을 만들어 메일을 발송한다.
	//		2 : 마케팅용 이메일 이다. (inline CSS를 이용한 배너가 이메일로 발송 된다.)

	for _, r := range eventRewards {
		switch r.EventRewardType {
		case 1:
			//
			rewardsIDs, err := MakeCoupon(eventRewards)
			if err != nil {
				logs.Error("Error MakeCoupon: ", err)
				return
			}

			logs.Info("Success MakeCoupon: ", len(rewardsIDs))

			// Step 2
			// get GetSendEmailReward again with serial.
			/*
				var eventRewards []models.EventRewards
				eventRewards, err := models.GetSendEmailReward()
				if err != nil {
					logs.Error("Error ReGetSendEmailReward: ", err)
					return
				}
				logs.Info("Success ReGetSendEmailReward: ", len(eventRewards))
			*/

			// TODO: for loop with go routine.
			libs.MakeEmail(eventRewards)

			// bulk update
			// TODO: need to change 건 바이 건 ???
			/*
				err = models.UpdateEventRewardsDone(rewardsIDs)
				if err != nil {
					logs.Error("Error UpdateEventRewardsDone: ", err)
					return
				}

				logs.Info("Success UpdateEventRewardsDone. IDs: ", rewardsIDs)
			*/
		case 2:
			go libs.MakeEmailMarketing(r)
		}
	}

}

func TestHerems() {
	inputFile := "./index.html"

	prem, err := premailer.NewPremailerFromFile(inputFile, premailer.NewOptions())
	if err != nil {
		log.Fatal(err)
	}

	html, err := prem.Transform()
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(html)

	// send email

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
	m.SetHeader("To", "youngtip@gmail.com", "quybv90@gmail.com", "hankyeol@naddic.com", "kim.dokyung@naddic.com", "marmarhan2@gmail.com", "skskunk2@gmail.com")
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

// MakeCoupon ...
func MakeCoupon(er []models.EventRewards) ([]int, error) {
	rewardsIDs, err := models.MakeCoupon(er)
	return rewardsIDs, err
}

// GetGameReward ...
func GetGameReward() {

}

func init() {
	models.RegisterDB()
}

func main() {
	orm.Debug = true

	//orm.RunSyncdb("default", false, true)

	//logging
	logs.SetLogger(logs.AdapterFile, `{"filename":"./logs/project.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":30,"color":true}`)

	// logs.Error("test")

	// sendtype 따라 분기
	// TODO: 채널 이용
	// TODO: 포인터 이용

	//i := 1
	//for {

	//logs.Info("start: ", i)
	start := time.Now()

	// herems test
	TestHerems()

	// GetSendEmail()

	logs.Info("Total Time: ", time.Since(start))

	//i++
	//time.Sleep(30 * time.Minute)
	//}

}
