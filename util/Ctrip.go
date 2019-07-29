package util

func GetPackageId(timeLength int, distanceLength int) (int) {
	return timeLength*1000000 + distanceLength
}
