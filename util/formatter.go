package util

import (
	"fmt"
	"strings"
)

// FormatDate 日付文字列を日本語形式にフォーマット
// 例: "2024-03-15" -> "2024年3月15日"
func FormatDate(dateStr string) string {
	parts := strings.Split(dateStr, "-")
	if len(parts) != 3 {
		return dateStr
	}
	return fmt.Sprintf("%s年%s月%s日", parts[0], parts[1], parts[2])
}
