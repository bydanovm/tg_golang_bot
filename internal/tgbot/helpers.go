package tgbot

func FormatFloatToString(number float32) (format string) {
	// Установим формат общей длиной в 7 знаков
	if number >= 100000 {
		format = "%.1f"
	} else if number >= 10000 {
		format = "%.2f"
	} else if number >= 1000 {
		format = "%.3f"
	} else if number >= 100 {
		format = "%.4f"
	} else if number >= 10 {
		format = "%.5f"
	} else {
		format = "%.6f"
	}
	return format
}
