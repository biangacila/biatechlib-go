package dates

import (
	"fmt"
	"github.com/leekchan/accounting"
	"math/big"
	"strconv"
	"strings"
	"time"
)

func IsDateValidDateFormat(dateString string) bool {
	dateFormat := "2006-01-02"
	_, err := time.Parse(dateFormat, dateString)
	return err == nil
}

// "2006-02-01"
func GetCurrentDateFormat(dateFormat string) string {
	now := time.Now()
	n := fmt.Sprintf("%v", now.Format(dateFormat))
	return n
}
func ChangeDateFormat(dateString, formatName string) string {
	if dateString == "" {
		return dateString
	}
	if strings.Contains(dateString, "/") {
		return dateString
	}
	now := ConvertStringToDateTime(dateString)
	n := fmt.Sprintf("%v", now.Format(formatName))
	return n
}
func FormatCurrency2(strIn float64) string {
	ac := accounting.Accounting{Symbol: "", Precision: 2}
	s := ac.FormatMoneyBigFloat(big.NewFloat(strIn))
	s = strings.Replace(s, ",", " ", 10)
	return s
}
func FormatCurrency3(strIn float64) string {
	ac := accounting.Accounting{Symbol: "", Precision: 0}
	s := ac.FormatMoneyBigFloat(big.NewFloat(strIn))
	s = strings.Replace(s, ",", " ", 10)
	return s
}
func FormatCurrency(strIn float64) string {
	ac := accounting.Accounting{Symbol: "R", Precision: 2}
	s := ac.FormatMoneyBigFloat(big.NewFloat(strIn))
	s = strings.Replace(s, ",", " ", 10)
	return s
}
func DateToDay(dateIn string) string {
	date := ConvertStringToDateTime(dateIn)
	day := fmt.Sprintf("%v", date.Format("Mon"))
	return day
}
func DateToWeekNumber(dateIn string) int {
	date := ConvertStringToDateTime(dateIn)
	_, w := date.ISOWeek()
	return w
}
func WeekStartDate(date time.Time) time.Time {
	offset := (int(time.Monday) - int(date.Weekday()) - 7) % 7
	result := date.Add(time.Duration(offset*24) * time.Hour)
	return result
}
func WeekStart(year, week int) time.Time {
	// Start from the middle of the year:
	t := time.Date(year, 7, 1, 0, 0, 0, 0, time.UTC)

	// Roll back to Monday:
	if wd := t.Weekday(); wd == time.Sunday {
		t = t.AddDate(0, 0, -6)
	} else {
		t = t.AddDate(0, 0, -int(wd)+1)
	}

	// Difference in weeks:
	_, w := t.ISOWeek()
	t = t.AddDate(0, 0, (week-w)*7)

	return t
}
func WeekRange(year, week int) (start, end time.Time) {
	start = WeekStart(year, week)
	end = start.AddDate(0, 0, 6)
	return
}
func ConvertStringToDateTime(dateIn string) time.Time {
	arr := strings.Split(dateIn, "-")
	year, _ := strconv.Atoi(arr[0])
	month, _ := strconv.Atoi(arr[1])
	day, _ := strconv.Atoi(arr[2])
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return date
}
func IsDateGreatThan(targetDate, compareDate string) bool {
	tDate := ConvertStringToDateTime(targetDate)
	tCompare := ConvertStringToDateTime(compareDate)
	return !tCompare.After(tDate)
}
func IsDateBetween(targetDate, startDate, endDate string) bool {
	if targetDate == startDate || targetDate == endDate {
		return true
	}
	check := ConvertStringToDateTime(targetDate)
	start := ConvertStringToDateTime(startDate)
	end := ConvertStringToDateTime(endDate)
	return check.After(start) && check.Before(end)
}
func ConvertMonthFromFullNameToDigitString(inValue string) string {
	var my = make(map[string]string)
	my["January"] = "1"
	my["February"] = "2"
	my["March"] = "3"
	my["April"] = "4"
	my["May"] = "5"
	my["June"] = "6"
	my["July"] = "7"
	my["August"] = "8"
	my["September"] = "9"
	my["October"] = "10"
	my["November"] = "11"
	my["December"] = "12"

	month, _ := my[inValue]
	return month
}
func FloatToInt(inVal float64) int {
	str := fmt.Sprintf("%v", inVal)
	res, _ := strconv.Atoi(str)
	return res
}
func RemoveTrim(val string) string {
	val = strings.Trim(val, " ")
	val = strings.TrimSpace(val)
	val = strings.Replace(val, "  ", " ", 100)
	return val
}
