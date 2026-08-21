package main

import (
	"cloudian/cloudian-claw/internal/domain/entity"
	"fmt"
)

func main() {
	fmt.Println("Hello world. Coding Agent build with Cloudian Love Cloud")
	//Using for testing only 
	agent := entity.Agent{
		ID: "agent-001", 
		Name: "Cloudian Agent", 
		IsDefault: true, 
		CreatedBy: "cloudian-001",
	} 
	agentError := agent.Validate()
	if  agentError == nil {
		fmt.Println("Thong tin hop le") 
		fmt.Println(agent) 
	} else {
		fmt.Println(agent, " has : " , agentError)
	}
}