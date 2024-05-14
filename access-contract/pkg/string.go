package pkg

// UniqueStrings Function to make a list of strings unique
func UniqueStrings(input []string) []string {
	uniqueMap := make(map[string]bool)
	var uniqueList []string

	// Iterate through the input list
	for _, str := range input {
		// If the string is not already in the map, add it
		if _, ok := uniqueMap[str]; !ok {
			uniqueMap[str] = true
			uniqueList = append(uniqueList, str)
		}
	}

	return uniqueList
}
