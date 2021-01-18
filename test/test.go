package main

import "fmt"
import "time"

func main() {
	for i := 0; i < 10; i++ {
		j := i
		go func() {
			fmt.Println(j)
		}()
	}
	time.Sleep(time.Duration(time.Second))
}
