package utl

func In(needle interface{}, haystack ...interface{}) bool {
	for _, straw := range haystack {
		if straw == needle {
			return true
		}
	}
	return false
}
