package g

import (
	"database/sql/driver"
	"fmt"
	"time"
)

/*
   @File: timestamp.go
   @Author: khaosles
   @Time: 2023/6/16 14:25
   @Desc: json格式化返回 YYYY-mm-dd HH:MM:SS 格式的时间字段
*/

// 时间格式
const timeLayout = time.DateTime

type Timestamp time.Time

func (ts *Timestamp) UnmarshalJSON(data []byte) error {
	// 加载上海时区
	location, _ := time.LoadLocation("Asia/Shanghai")
	t, err := time.ParseInLocation(`"`+timeLayout+`"`, string(data), location)
	if err != nil {
		return err
	}
	*ts = Timestamp(t)
	return nil
}

func (ts Timestamp) MarshalJSON() ([]byte, error) {
	t := time.Time(ts)
	formatted := fmt.Sprintf(`"%s"`, t.Format(timeLayout))
	return []byte(formatted), nil
}

func (ts Timestamp) Value() (driver.Value, error) {
	return time.Time(ts), nil
}

func (ts *Timestamp) Scan(value interface{}) error {
	if value == nil {
		*ts = Timestamp(time.Time{})
		return nil
	}
	if t, ok := value.(time.Time); ok {
		*ts = Timestamp(t)
		return nil
	}
	return fmt.Errorf("failed to scan CustomTime value")
}
