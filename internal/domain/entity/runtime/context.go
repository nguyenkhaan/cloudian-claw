package runtime

import "cloudian/cloudian-claw/internal/domain/entity"

//Tong hop toan bo thong tin cho Agent
type ExecutionContext struct {
	ExecutionID string //Xac nhan xem thuoc ve execution nao trong database 
	Messages []entity.Message 
	Model string 
	Provider string //Nha cung cap 
	Tools []entity.Tool 
	Skills []entity.Skill 
	Rules []entity.Rule
	Request ExecutionRequest
	SystemPrompt string //Chinh xuat ra va dua vao agent o lan cuoi cung 
	Workspace string 

}