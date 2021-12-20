package utils

// DupCount creates a map of item:sum for each duplicate item in a slice
func DupCount(list []string) map[string]int {
	dupFreq := make(map[string]int)

	for _, item := range list {
		_, exist := dupFreq[item]

		if exist {
			dupFreq[item]++
		} else {
			dupFreq[item] = 1
		}
	}
	return dupFreq
}
