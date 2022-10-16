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
