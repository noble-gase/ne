package ne

import "time"

// GMT8 东八区时区
var GMT8 = time.FixedZone("CST", 8*3600)

// WeekAround 返回给定时间戳所在周的「周一」和「周日」时间字符串
func WeekAround(layout string, now time.Time) (monday, sunday string) {
	weekday := now.Weekday()

	// monday
	offset := int(time.Monday - weekday)
	if offset > 0 {
		offset = -6
	}
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	monday = today.AddDate(0, 0, offset).Format(layout)

	// sunday
	offset = int(time.Sunday - weekday)
	if offset < 0 {
		offset += 7
	}
	sunday = today.AddDate(0, 0, offset).Format(layout)

	return
}
