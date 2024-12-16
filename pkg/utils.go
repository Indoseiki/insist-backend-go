package pkg

import (
	"fmt"
	"strconv"
	"time"
)

func CalculateServiceDuration(employeeNumber string) string {
	joinYear, _ := strconv.Atoi(employeeNumber[:4])
	joinMonth, _ := strconv.Atoi(employeeNumber[4:6])
	currentDate := time.Now()
	currentYear := currentDate.Year()
	currentMonth := int(currentDate.Month())

	years := currentYear - joinYear
	months := currentMonth - joinMonth

	if months < 0 {
		years--
		months += 12
	}

	return fmt.Sprintf("%d Years, %d Months", years, months)
}
