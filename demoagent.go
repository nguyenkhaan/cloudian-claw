package main 
// // https://stackoverflow.com/questions/67248551/go-get-values-from-env-file

// package main

// import (
// 	"bufio"
// 	"context"
// 	"fmt"
// 	"log"
// 	"os"
// 	"strings"

// 	"github.com/joho/godotenv"
// 	"github.com/openai/openai-go/v3"
// 	"github.com/openai/openai-go/v3/option"
// )
// func main() {
// 	err := godotenv.Load() 
// 	if err != nil {
// 		fmt.Println("Error while loading .env")
// 		os.Exit(1) 
// 	}
// 	apiKey := os.Getenv("OPENROUTER_API_KEY")
// 	if apiKey == "" {
// 		log.Fatal("API key is required")
// 		os.Exit(1)
// 	}
// 	url := "https://openrouter.ai/api/v1"

// 	client := openai.NewClient(
// 		option.WithBaseURL(url),  
// 		option.WithAPIKey(apiKey), 
// 	)
// 	messages := []openai.ChatCompletionMessageParamUnion{}
// 	model := "z-ai/glm-5.2:free"
// 	ctx := context.Background() //Khong cho dung cong bviec giua chung, context.Background() noi rang cong viec nay la vo han thoi gian 
// 	//hay thuc hien cong viec nay cho den khi nao hoan thanh xong cong viec 
// 	params := openai.ChatCompletionNewParams{
// 		Model: model, 
// 		Messages: messages,
// 	}
// 	var count = 0 
// 	scanner := bufio.NewScanner(os.Stdin)
// 	for {
// 		fmt.Print("> ") 
// 		if !scanner.Scan() {
// 			break 
// 		}
// 		//Su dung goi thu vien strings de xu ly chuoi 
// 		input := strings.TrimSpace(
// 			scanner.Text(), 
// 		)
// 		if input == "" {
// 			continue
// 		}
// 		if input == "clear" {
// 			//Clear the terminal = Print the spcial key 
// 			fmt.Print("\033[H\033[2J")
// 			continue 
// 		}
// 		if input == "exit" {
// 			fmt.Println("Goodnight") 
// 			break 
// 		}
// 		//Append the messages to the param's messages 
// 		params.Messages = append(params.Messages , openai.UserMessage(input))
		
// 		response, err := client.Chat.Completions.New(ctx , params)
// 		if err != nil {
// 			log.Fatal(err)
// 		} 
// 		output := response.Choices[0].Message.Content 

// 		fmt.Println("[AGENT]: ", output) 
// 		params.Messages = append(params.Messages , openai.AssistantMessage(output)) 
// 		//Build the agent loop fuck u :v 
// 		count++ 
// 		if count >= 10 {
// 			params.Messages = params.Messages[(count-4) : ]
// 			count = 0; 
// 		}
// 	}

// }