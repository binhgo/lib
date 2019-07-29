package util

import "math"

func CalcGuestNumber( adultNumber int,childNumber int) int {
	guestNumber := 0
	if (childNumber <= 1) {
		guestNumber = childNumber + adultNumber;
	} else {
		guestNumber = int(math.Floor(float64(childNumber)*1.5 + float64(adultNumber) + 0.5));
	}

	return guestNumber
}
