package util

func CheckError(err error) {

	if err != nil {
		panic(err)
	}

}

func SliceRemove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}
