package application

import "time"

// DisplayItem wraps the file stat info and string to be printed
type DisplayItem struct {
	FileInfo
	Time    *time.Time
	Display []string
}

type DisplayItemList []*DisplayItem

func (list DisplayItemList) Len() int {
	return len(list)
}

func (list DisplayItemList) Get(index int) []string {
	return list[index].Display
}
