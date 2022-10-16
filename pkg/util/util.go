package util

func CheckError(err error) {

	if err != nil {
		panic(err)
	}

}

var Check = CheckError

func SliceRemove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

// Sees if a string splice contains a specified string
func ContainsString(splice []string, str string) bool {

	for _, v := range splice {

		if v == str {
			return true
		}

	}

	return false

}
