package stringutils

func StrRevers(s string) string {

	runes := []rune(s)
	length := len(runes)

	reversed := make([]rune, length)

	for i := 0; i < length; i++ {
		reversed[i] = runes[length-1-i]

	}
	return string(reversed)

}
