package main

import (
	"fmt"
	"time"
)

// 判断date是否合法
func dateCheck(date string) bool {
	return len(date) == 8	
}

func GetDate(biasDays int) string {

	p := time.Now().Add(-24 * time.Hour * time.Duration(biasDays))
	year := p.Year()
	month := int(p.Month())
	day := p.Day()
	date := fmt.Sprintf("%04d%02d%02d", year, month, day)
	return date
}
