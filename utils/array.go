package utils

func IsIn(item string,arr []string) bool{
	for _,v := range arr{
		if item == v{
			return true
		}
	}
	return false
}