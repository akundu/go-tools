package main

func checkPrime(num int) bool {
	if num < 0 {
		num = -num
	}
	switch {
	case num < 2:
		return false
	default:
		for index := 2; index < num; index++ {
			if num%index == 0 {
				return false
			}
		}
	}
	return true
}

//TODO: save the highest prime numbers - as a way to simulate dynamic programming
