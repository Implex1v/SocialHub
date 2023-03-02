package messaging

func strFirstCharacters(string string, n int) string {
	l := min(len(string), n)
	return string[:l]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
