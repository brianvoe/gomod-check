package main

func strPadding(str string, length int) string {
	curLength := len(str)
	if curLength == length {
		return str
	}

	difference := length - curLength
	for i := 0; i < difference; i++ {
		str += " "
	}

	return str
}
