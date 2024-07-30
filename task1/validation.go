package main

func validateGrade(grade float64) (string, bool) {
	if grade < 0 || grade > 100 {
		return "grade range should be between 0 - 100", false
	}
	return "", true
}
