package models

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"

	"github.com/astaxie/beego/orm"
	"github.com/ndcinfra/platform/libs"
)

// Event ...
// 이벤트 관리 테이블
type Event struct {
	ID   int    `orm:"column(ID);pk;auto" json:"id"`
	Desc string `orm:"size(500);" json:"desc"` // 설명

	/*
		EventFrom time.Time `orm:"type(datetime);null" json:"event_from"` // 이벤트 시작 일자
		EventTo   time.Time `orm:"type(datetime);null" json:"event_to"`   // 이벤트 끝 일자.

		// event_reward_type
		// 1: send email with cououpon serial
		// 2: reward game item to game directly.
		// 3: reward free coin to platform
		EventRewardType int `json:"event_reward_type"` // 이벤트 리워드 타입
	*/

	Admin    string    `orm:"size(30);" json:"admin"`                       // admin aacount
	CreateAt time.Time `orm:"type(datetime);auto_now_add" json:"create_at"` // 이 테이블 등록 일자
}

// EventRewards ...
//
type EventRewards struct {
	ID int `orm:"column(ID);pk;auto" json:"id"`

	// 필수 항목
	EventID int `orm:"column(EventID);" json:"event_id"` // 위 이벤트 관리 테이블 ID
	// event_reward_type
	// 1: send email with cououpon serial
	// 2: reward game item to game directly.
	// 3: reward free coin to platform
	EventRewardType int    `json:"event_reward_type"`               // 이벤트 리워드 타입
	UID             string `orm:"column(UID);size(50);" json:"uid"` // user id

	// event_reward_type이 1일때 해당 되는 항목
	Displayname string    `orm:"size(30);null" json:"displayname"`  // 4 ~ 16 letters for local,
	Email       string    `orm:"size(100);null" json:"email"`       // max 100 letters
	MID         int       `orm:"column(MID);default(0)" json:"mid"` // email message id
	Serial      string    `orm:"size(50);null" json:"serial"`
	CouponType  int       `json:"coupon_type"`                            // coupon type
	CouponStart time.Time `orm:"type(datetime);null" json:"coupon_start"` // coupon start date
	CouponEnd   time.Time `orm:"type(datetime);null" json:"coupon_end"`   // coupon end date
	SendAt      time.Time `orm:"type(datetime);null" json:"send_at"`      // 희망 발송 시작 일자.

	ItemID    int `orm:"column(ItemID);default(0)" json:"itemid"`       // game item id
	PayItemID int `orm:"column(PayItemID);default(0)" json:"payitemid"` // payment item id. 무료 코인 쿠폰용

	Admin    string    `orm:"size(30);" json:"admin"`                       // admin aacount
	CreateAt time.Time `orm:"type(datetime);auto_now_add" json:"create_at"` // 이 테이블 등록 일자
	UpdateAt time.Time `orm:"type(datetime);auto_now_add" json:"update_at"` // 처리 일자.
	IsDone   bool      `orm:"null;default(false)" json:"is_done"`           // 처리 여부
}

// GetSendEmailReward ...
func GetSendEmailReward() ([]EventRewards, error) {
	var eventRewards []EventRewards

	o := orm.NewOrm()
	sql := "SELECT " +
		" \"ID\" , " +
		" \"EventID\", " +
		" event_reward_type, " +
		" \"UID\", " +
		" displayname, " +
		" email, " +
		" \"MID\", " +
		" serial, " +
		" coupon_type, " +
		" coupon_start, " +
		" coupon_end, " +
		" \"ItemID\", " +
		" \"PayItemID\" " +
		" FROM event_rewards " +
		" WHERE is_done = false " +
		" AND event_reward_type = 1 " + // 1: send email with cououpon serial
		" order by create_at asc " // +
		// " limit 50 "

	_, err := o.Raw(sql).QueryRows(&eventRewards)
	return eventRewards, err
}

// Make Coupon
//
func MakeCoupon(er []EventRewards) ([]int, error) {
	rewardsIDs := make([]int, len(er), len(er))

	o := orm.NewOrm()
	err := o.Begin()
	var sql string

	// init CID
	CID := libs.GenerateID("C") // platform 과 동일 버젼
	initEventID := er[0].EventID

	// for loop
	//		insert coupon_use_history
	for i := 0; i < len(er); i++ {
		// save id for bulk update done
		rewardsIDs[i] = er[i].ID

		serial := makeSerial()

		sql = "INSERT INTO coupon_use_history " +
			"(\"CID\", serial, \"from\", \"to\", type, \"ItemID\", \"PayItemID\", create_at, \"UID\", is_used, is_fixed) " +
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)"

		_, err = o.Raw(sql, CID, serial, er[i].CouponStart, er[i].CouponEnd, er[i].CouponType, er[i].ItemID, er[i].PayItemID, time.Now(), er[i].UID, false, true).Exec()
		if err != nil {
			logs.Error("error", "Error insert coupon_use_history fixed Coupon: ", err)
			_ = o.Rollback()
			return nil, err
		}

		// update serial in event_rewards
		sql = "UPDATE event_rewards SET serial = ? where \"ID\"= ? "
		_, err = o.Raw(sql, serial, er[i].ID).Exec()
		if err != nil {
			logs.Error("error", "Error update event_rewards: ", err)
			_ = o.Rollback()
			return nil, err
		}

		if initEventID != er[i].EventID {
			// new event
			CID = libs.GenerateID("C") // platform 과 동일 버젼
			initEventID = er[i].EventID
		}

	}

	err = o.Commit()

	return rewardsIDs, err

}

func UpdateEventRewardsDone(ids []int) error {
	var inSQL string

	for i := 0; i < len(ids); i++ {
		if i != len(ids)-1 {
			inSQL += strconv.Itoa(ids[i]) + ", "
		} else {
			inSQL += strconv.Itoa(ids[i])
		}
	}

	o := orm.NewOrm()
	sql := "UPDATE event_rewards SET is_done = true where \"ID\" in ( " + inSQL + " )"
	_, err := o.Raw(sql).Exec()
	if err != nil {
		logs.Error("error", "Error update done event_rewards: ", err)
		return err
	}

	return err

}

func makeSerial() string {

	b := make([]byte, 6) //equals 8 charachters
	rand.Read(b)
	s := hex.EncodeToString(b)
	s = strings.ToUpper(s)

	return s
}
