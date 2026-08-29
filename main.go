package main

import (
	"cloudian/cloudian-claw/internal/impl/config"
	"context"
	"fmt"

	"github.com/joho/godotenv"
)
func main() {
	godotenv.Load()  
	ctx := context.Background() 

	container, err := config.NewContainer(ctx) 
	if err != nil {
		fmt.Println("Error while loading container: " , err) 
		return 
	} 
	fmt.Print(container) 
	applicationConfig := container.Config 
	fmt.Println(applicationConfig.ServerHost)
}
