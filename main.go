package main

import (
	"context"
	"fmt"
)
func worker(ctx context.Context) {
	username := ctx.Value("username") 
	fmt.Println(username) 
}
func main() {
	ctx := context.Background() 
	fmt.Println(ctx.Done()) 
	fmt.Print(ctx.Deadline())
	
	fmt.Println(ctx.Err())

}
