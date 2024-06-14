package main

import (
	"context"
	"fmt"
	"log"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/document"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"

	
)


func processConverseStreamingOutput(
	output *bedrockruntime.ConverseStreamOutput,
	handler ConverseStreamingOutputDeltaHandler,
	start ConverseStreamingOutputStartHandler) (Response, error) {

	resp := Response{}
	
	for event := range output.GetStream().Events() {
		switch v := event.(type) {
		case *types.ConverseStreamOutputMemberContentBlockDelta:
			handler(context.Background(), v.Value.Delta)
		case *types.ConverseStreamOutputMemberContentBlockStart:
			start(context.Background(), v.Value.Start)
			//print(1)
		case *types.ConverseStreamOutputMemberContentBlockStop:
			//print(2)
		case *types.ConverseStreamOutputMemberMetadata:
			//print(3)
		case *types.ConverseStreamOutputMemberMessageStart:
			//print(4)
		case *types.ConverseStreamOutputMemberMessageStop:
			resp.StopReason=string(v.Value.StopReason)

		default:
			fmt.Println("union is nil or unknown type")
		}
	}
	return resp, nil
}

// generateTextStream 使用 Bedrock 模型生成文本
func generateTextStream(ctx context.Context, client *bedrockruntime.Client, modelID, inputText string, toolConfig *types.ToolConfiguration) error {
	
	var outputText string
	var toolName string
	var toolUseIDs = make(map[string]string)
	var parameters = make(map[string]string)
	log.Printf("Generating text with model %s", modelID)

	// 创建初始消息
	initialMessage := types.Message{
		Role: "user",
		Content: []types.ContentBlock{
			&types.ContentBlockMemberText{
				Value: inputText,
			},
		},
	}
	messages := []types.Message{initialMessage}
	
	system:=[]types.SystemContentBlock{
		&types.SystemContentBlockMemberText{
			Value: "You are AI assistant",
		},
	}

	printJSON(messages)

	


 for {

	output, err := client.ConverseStream(context.Background(), &bedrockruntime.ConverseStreamInput{
		System: system,
		Messages:   messages,
		ModelId:    &modelID,
		ToolConfig: toolConfig,
	})

	if err != nil {
		return err
	}
	resp,err:=processConverseStreamingOutput(output,
		func(ctx context.Context, part types.ContentBlockDelta) error {
			switch d := part.(type) {
			case *types.ContentBlockDeltaMemberText:
				var chunk = d.Value
				outputText += chunk
				
			case *types.ContentBlockDeltaMemberToolUse:
				var chunk = d.Value
				parameters[toolName]+=*chunk.Input
				

			}
			return nil
		},
		func(ctx context.Context, start types.ContentBlockStart) error {
			switch d := start.(type) {
			case *types.ContentBlockStartMemberToolUse:
				print(fmt.Sprintf("\n====tool use=====%s, %s\n",*d.Value.ToolUseId,*d.Value.Name))
				toolName=*d.Value.Name
				toolUseIDs[toolName]=*d.Value.ToolUseId
			}
			return nil
		},
	)
	if err!=nil {
		fmt.Println(err)
		break
	}
	
	
	if resp.StopReason== string(types.StopReasonToolUse){
		
		if(outputText==""){
			continue
		}
		

		data := make(map[string]interface{})
		json.Unmarshal([]byte(parameters[toolName]),&data)

		content := document.NewLazyDocument(data)
	
		
		outputMessage := types.Message{
			Role: "assistant",
			Content: []types.ContentBlock{
				&types.ContentBlockMemberText{
					Value: outputText,
				},
				&types.ContentBlockMemberToolUse{
					Value: types. ToolUseBlock{
						ToolUseId: aws.String(toolUseIDs[toolName]),
						Name: aws.String(toolName),
						Input: content,
					},
				},
			},
		}

		

		messages = append(messages, outputMessage)

		
		
		printJSON(messages)

		handleToolUseSteram(toolUseIDs[toolName],toolName,parameters[toolName],&messages)
		
	}

	if resp.StopReason==string(types.StopReasonEndTurn) {
		print(outputText)
		break
	}
 }
	
	
	
	return nil
}


// generateText 使用 Bedrock 模型生成文本
func generateText(ctx context.Context, client *bedrockruntime.Client, modelID, inputText string, toolConfig *types.ToolConfiguration) error {
	log.Printf("Generating text with model %s", modelID)

	// 创建初始消息
	initialMessage := types.Message{
		Role: "user",
		Content: []types.ContentBlock{
			&types.ContentBlockMemberText{
				Value: inputText,
			},
		},
	}
	messages := []types.Message{initialMessage}
	
	system:=[]types.SystemContentBlock{
		&types.SystemContentBlockMemberText{
			Value: "You are AI assistant",
		},
	}

	printJSON(messages)

	// 与模型进行对话
	for {
		response, err := client.Converse(ctx, &bedrockruntime.ConverseInput{
			ModelId:    &modelID,
			Messages:   messages,
			ToolConfig: toolConfig,
			System: system,
		})
		if err != nil {
			return err
		}

		printJSON(response.Output)
		printJSON(response.Usage)
		if response.StopReason == types.StopReasonToolUse {
			// 处理工具使用情况
			handleToolUse(response.Output, &messages)
		} else {
			// 对话结束,打印最终输出
			switch v := response.Output.(type) {
			case *types.ConverseOutputMemberMessage:
				for _,content:=range(v.Value.Content){
					switch c := content.(type) {
						case *types.ContentBlockMemberText:
							fmt.Println(c.Value)

					}
				}
			}	
			break
		}
	}

	return nil
}