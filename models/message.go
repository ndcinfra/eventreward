package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

// EmailMessage ...
// 이메일 발송 메세지 관리
type EmailMessage struct {
	ID   int    `orm:"column(ID);pk;auto" json:"id"`
	Desc string `orm:"size(100);" json:"desc"`

	Title string `orm:"size(1000);" json:"title"`
	Body1 string `orm:"size(4000);" json:"body1"`
	Body2 string `orm:"size(4000);" json:"body2"`
	Body3 string `orm:"size(4000);" json:"body3"`
	Body4 string `orm:"size(4000);" json:"body4"`
	Body5 string `orm:"size(4000);" json:"body5"`

	CreateAt time.Time `orm:"type(datetime);auto_now_add" json:"create_at"` // 이 테이블 등록 일자
}

// GetMessage ...
func (e *EmailMessage) GetMessage() error {
	o := orm.NewOrm()
	sql := "SELECT " +
		" \"ID\" , " +
		" \"desc\", " +
		" Title, " +
		" Body1, " +
		" Body2, " +
		" Body3, " +
		" Body4, " +
		" Body5 " +
		" FROM email_message " +
		" WHERE \"ID\" = ? "

	err := o.Raw(sql, e.ID).QueryRow(&e)
	return err
}
