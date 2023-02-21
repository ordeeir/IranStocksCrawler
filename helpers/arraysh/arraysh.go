package arraysh

func Contains(array []string, search string) bool {
	for _, a := range array {
		if a == search {
			return true
		}
	}
	return false
}
