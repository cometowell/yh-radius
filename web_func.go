package main

import "time"

const PaginationBarSize int = 10

func FormatDateTime(datatime time.Time) string {
	return datatime.Format(DateTimeFormat)
}

func PaginationBar(current int64) []int {
	_current := int(current)
	base := _current / PaginationBarSize
	pbs := [PaginationBarSize]int{}

	for i:=0;i<PaginationBarSize; i++ {
		pbs[i] = base * 10 + i
	}
	return pbs[:]
}
