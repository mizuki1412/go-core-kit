package mathkit

func GroupNum(sum int, group int) int {
	if sum == 0 {
		return 1
	}
	if sum%group == 0 {
		return sum / group
	} else {
		return sum/group + 1
	}
}
