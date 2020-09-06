package watchman

func Pcre(regexp string) []string {
	return []string{"pcre", regexp}
}
