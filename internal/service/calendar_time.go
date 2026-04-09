package service

import "time"

var taipeiLoc = func() *time.Location {
	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		return time.FixedZone("CST", 8*3600)
	}
	return loc
}()

// resolveCalendarQueryTime converts a query date into a stable civil-day sampling time.
// We use Asia/Taipei 12:00 to avoid boundary-day drift around solar-term transition instants.
func resolveCalendarQueryTime(dateStr string) (time.Time, error) {
	if dateStr == "" {
		now := time.Now().In(taipeiLoc)
		return time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, taipeiLoc), nil
	}

	d, err := time.ParseInLocation("2006-01-02", dateStr, taipeiLoc)
	if err != nil {
		return time.Time{}, err
	}

	return time.Date(d.Year(), d.Month(), d.Day(), 12, 0, 0, 0, taipeiLoc), nil
}
