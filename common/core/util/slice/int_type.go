package slice

func InArray(val int, list []int) bool {
	res := false
	for _, item :=  range list {
		if item == val {
			res = true
			break
		}
	}
	return res
}
