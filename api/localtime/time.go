package localtime

import (
	"os"
	"time"
)

func NowTime() time.Time {
	loc, err := time.LoadLocation(os.Getenv("TZ"))
	if err != nil {
		loc, _ = time.LoadLocation("Asia/Seoul")
	}

	return time.Now().In(loc)
}
