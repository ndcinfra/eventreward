package main

import (
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/ndcinfra/eventreward/libs"
	"github.com/ndcinfra/eventreward/models"

	"github.com/astaxie/beego/orm"
)

// 프로그램의 기능
// event_reward_type
// 1: send email with cououpon serial
//		insert coupon
//		insert coupon_send_history and get serial
//		send email
// 2: reward game item to game directly.
//
// 3: reward free coin to platform
//		TODO: 나중에 개발.
//

// 우선 for loop를 사용 해서 개발
// TODO: 향후 channel을 이용한다. !!!
// TODO: 포인터 array 이용 !!!

// GetSendEmailReward ...
func GetSendEmailReward() {
	var eventRewards []models.EventRewards
	eventRewards, err := models.GetSendEmailReward()
	if err != nil {
		logs.Error("Error GetSendEmailReward: ", err)
		return
	}

	logs.Info("Success GetSendEmailReward: ", len(eventRewards))

	// make coupon and coupon_send_history
	if len(eventRewards) > 0 {
		rewardsIDs, err := MakeCoupon(eventRewards)
		if err != nil {
			logs.Error("Error MakeCoupon: ", err)
			return
		}

		logs.Info("Success MakeCoupon: ", len(rewardsIDs))

		// Step 2

		// get GetSendEmailReward again with serial.
		eventRewards, err = models.GetSendEmailReward()
		if err != nil {
			logs.Error("Error ReGetSendEmailReward: ", err)
			return
		}

		logs.Info("Success ReGetSendEmailReward: ", len(eventRewards))

		// bulk update
		/*
			err = models.UpdateEventRewardsDone(rewardsIDs)
			if err != nil {
				logs.Error("Error UpdateEventRewardsDone: ", err)
				return
			}

			logs.Info("Success UpdateEventRewardsDone. IDs: ", rewardsIDs)
		*/

		// send email
		// go libs.MakeEmail(eventRewards)
		libs.MakeEmail(eventRewards)

	} else {
		logs.Info("no data GetSendEmailReward")
		return
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
	orm.RunSyncdb("default", false, true)

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
	GetSendEmailReward()
	logs.Info("Total Time: ", time.Since(start))

	//i++
	//time.Sleep(30 * time.Minute)
	//}

}
