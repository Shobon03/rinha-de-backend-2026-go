package util

var monthOffsets = [12]int{0, 3, 2, 5, 0, 3, 5, 1, 4, 6, 2, 4}
var daysInMonth = [13]int{0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

func ParseFastDate(dateStr string) (hour int, mappedWeekday int, absMinutes int) {
	if len(dateStr) < 16 {
		return 0, 0, 0
	}

	// Manual parsing of YYYY-MM-DDTHH:mm to avoid the overhead of time.Parse
	y := int(dateStr[0]-'0')*1000 + int(dateStr[1]-'0')*100 + int(dateStr[2]-'0')*10 + int(dateStr[3]-'0')
	m := int(dateStr[5]-'0')*10 + int(dateStr[6]-'0')
	d := int(dateStr[8]-'0')*10 + int(dateStr[9]-'0')
	h := int(dateStr[11]-'0')*10 + int(dateStr[12]-'0')
	min := int(dateStr[14]-'0')*10 + int(dateStr[15]-'0')

	hour = h

	// Sakamoto's algorithm to calculate the day of the week efficiently
	// 0=Sunday, 1=Monday, ..., 6=Saturday
	ySakamoto := y
	if m < 3 {
		ySakamoto -= 1
	}
	weekday := (ySakamoto + ySakamoto/4 - ySakamoto/100 + ySakamoto/400 + monthOffsets[m-1] + d) % 7

	// Map weekday so it starts on Monday (0=Monday, ..., 6=Sunday)
	mappedWeekday = (weekday + 6) % 7

	// Calculate absolute days since year 0 to determine a total minutes counter
	// This helps in calculating time differences without full timestamp objects
	leapDays := y/4 - y/100 + y/400
	days := y*365 + leapDays
	for i := 1; i < m; i++ {
		days += daysInMonth[i]
	}
	// Add an extra day if it's a leap year and we are past February
	if m > 2 && (y%4 == 0 && (y%100 != 0 || y%400 == 0)) {
		days++
	}
	days += d

	// Convert everything to absolute minutes for distance calculations
	absMinutes = days*1440 + h*60 + min
	return
}
