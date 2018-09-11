package main

import (
	"fmt"

	"gitee.com/goshark/gtimers"
)

func main() {
	timer := gtimers.NewTimer()
	timer.SetInterval("1m")
	timer.SetRepeatcount(3)

	timer.SyncStart(func() {
		//doing sth at this time..
		fmt.Println("doing something for current function..")
	})
}
