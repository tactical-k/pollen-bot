package util

import (
	"fmt"
	"strconv"
	"strings"
)

// FormatDate 日付文字列を日本語形式にフォーマット
// 例: "2024-03-15" -> "2024年3月15日"
func FormatDate(dateStr string) string {
	parts := strings.Split(dateStr, "-")
	if len(parts) != 3 {
		return dateStr
	}

	// 年月日を整数に変換して先頭のゼロを削除
	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return dateStr
	}
	month, err := strconv.Atoi(parts[1])
	if err != nil {
		return dateStr
	}
	day, err := strconv.Atoi(parts[2])
	if err != nil {
		return dateStr
	}

	return fmt.Sprintf("%d年%d月%d日", year, month, day)
}
