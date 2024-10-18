package utils

import (
	"fmt"
	"math"
)

// ConvertBytesToReadable converts an integer byte value to a human-readable string with units (KB, MB, GB, TB, PB).
func ConvertBytesToReadable(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%dB", bytes) // 小于1024字节时，直接返回字节数
	}
	// 存储单位字符 'K', 'M', 'G', 'T', 'P'
	sizes := []string{"KB", "MB", "GB", "TB", "PB"}
	div, exp := float64(unit), 0
	for n := float64(bytes) / unit; n >= unit && exp < len(sizes)-1; n /= unit {
		div *= unit
		exp++
	}
	// 使用 math.Floor 截断为两位小数而非四舍五入
	result := float64(bytes) / div
	truncatedResult := math.Floor(result*100) / 100
	return fmt.Sprintf("%.2f %s", truncatedResult, sizes[exp])
}
