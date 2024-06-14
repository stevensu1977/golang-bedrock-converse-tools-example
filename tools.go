package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/document"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
)

type Tools interface {
	
	GenerateToolSchema() *types.ToolMemberToolSpec
	GenerateToolResult(string , map[string]interface{}) (*types.Message,error)
	Invoke(string, string, map[string]interface{}) (*types.Message,error)
}

type WeatherTool struct {
	Name        string
	Description string
}

type LocationTool struct {
	Name        string
	Description string
}

var weatherTool=&WeatherTool{
	Name: "get_weather",
	Description: "Returns weather data for a given latitude and longitude",
}

var locationTool=&LocationTool{
	Name: "get_lat_long",
	Description: "Returns weather data for a given latitude and longitude",
}


func (w *WeatherTool) GenerateToolSchema() *types.ToolMemberToolSpec {
	getWeather := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"latitude": map[string]interface{}{
				"type":        "string",
				"description": "The latitude coordinate as a string.",
			},
			"longitude": map[string]interface{}{
				"type":        "string",
				"description": "The longitude coordinate as a string.",
			},
		},
		"required": []interface{}{"latitude","longitude"},
	}

	getWeatherDoc := document.NewLazyDocument(getWeather)
	return 	&types.ToolMemberToolSpec{
				Value: types.ToolSpecification{
					Name:        aws.String(w.Name),
					Description: aws.String(w.Description),
					InputSchema: &types.ToolInputSchemaMemberJson{
						Value: getWeatherDoc,
					},
				},
			}
}

func (w *WeatherTool) GenerateToolResult(toolUseID string, result map[string]interface{}) (*types.Message,error) {
	
	data := make(map[string]interface{})
	data["weather"]=result
	content := document.NewLazyDocument(data)

	return &types.Message{
		Role: "user",
		Content: []types.ContentBlock{
			&types.ContentBlockMemberToolResult{
				Value: types.ToolResultBlock{
					ToolUseId: &toolUseID,
					Content: []types.ToolResultContentBlock{
						&types.ToolResultContentBlockMemberJson{
							Value: content,
						},
					},
				},
			},
		},
	},nil
	
}

func (w *WeatherTool) Invoke(toolUseId string, toolName string, parameters map[string]interface{}) (*types.Message,error) {
	// 在这里实现工具的主要逻辑
	
	if toolName!=w.Name {
		return nil,fmt.Errorf("%s not match %s",toolName,w.Name)
	}
	latitude, hasLatitude := parameters["latitude"].(string)
	longitude, hasLongitude := parameters["longitude"].(string)

	
	if !hasLatitude || !hasLongitude {
		return nil,fmt.Errorf("parameter not correct, %+v",parameters)
	}
	result,err:=w.GetWeather(latitude,longitude)
	printJSON(result)
	if err!=nil {
		return nil,err
	}
	return w.GenerateToolResult(toolUseId,result)

}


func (w *WeatherTool) GetWeather(latitude, longitude string) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&current_weather=true", latitude, longitude)
	fmt.Println(url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return data, nil
}


func (w *LocationTool) GenerateToolSchema() *types.ToolMemberToolSpec {
	getLatLong := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"place": map[string]interface{}{
				"type":        "string",
				"description": "The place name to geocode and get coordinates for.",
			},
		},
		"required": []interface{}{"place"},
	}

	getLatLongDoc := document.NewLazyDocument(getLatLong)
	return 	&types.ToolMemberToolSpec{
				Value: types.ToolSpecification{
					Name:        aws.String(w.Name),
					Description: aws.String(w.Description),
					InputSchema: &types.ToolInputSchemaMemberJson{
						Value: getLatLongDoc,
					},
				},
			}
}

func (w *LocationTool) GenerateToolResult(toolUseID string, result map[string]interface{}) (*types.Message,error) {
	
	data := make(map[string]interface{})
	data["weather"]=result
	content := document.NewLazyDocument(data)

	return &types.Message{
		Role: "user",
		Content: []types.ContentBlock{
			&types.ContentBlockMemberToolResult{
				Value: types.ToolResultBlock{
					ToolUseId: &toolUseID,
					Content: []types.ToolResultContentBlock{
						&types.ToolResultContentBlockMemberJson{
							Value: content,
						},
					},
				},
			},
		},
	},nil
	
}

func (w *LocationTool) Invoke(toolUseId string, toolName string, parameters map[string]interface{}) (*types.Message,error) {
	// 在这里实现工具的主要逻辑
	if toolName!=w.Name {
		return nil,fmt.Errorf("%s not match %s",toolName,w.Name)
	}
	place, hasPlace := parameters["place"].(string)
	if !hasPlace  {
		return nil,fmt.Errorf("parameter not correct, %+v",parameters)
	}
	result,err:=w.GetLatLong(place)
	if err!=nil {
		return nil,err
	}
	return w.GenerateToolResult(toolUseId,result)

}

func  (w *LocationTool) GetLatLong(place string) (map[string]interface{}, error) {
	baseURL := "https://nominatim.openstreetmap.org/search"

	params := url.Values{}
	params.Add("q", place)
	params.Add("format", "json")
	params.Add("limit", "1")

	// Create a custom HTTP client with a custom User-Agent
	client := &http.Client{
			//Timeout: time.Second * 10, // Set a timeout if needed
	}

	req, err := http.NewRequest("GET", baseURL+"?"+params.Encode(), nil)
	if err != nil {
			return nil, err
	}

	// Set the User-Agent header
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")

	resp, err := client.Do(req)

	if err != nil {
			fmt.Println("err",err)
			return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
			return nil, err
	}
	fmt.Println(string(body))
	var data []map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
			return nil, err
	}
	
	if len(data) > 0 {
			return map[string]interface{}{"latitude": fmt.Sprintf("%s,",data[0]["lat"]), "longitude": fmt.Sprintf("%s,",data[0]["lon"])}, nil
	}

	return nil, nil
}




func handleToolUseSteram(toolUseID, toolName, parameter string, messages *[]types.Message) error {
	data, err := parseJSONParameter(parameter)
	if err != nil {
		return err
	}

	var message *types.Message
	switch toolName {
	case "get_weather":
		message, err = weatherTool.Invoke(toolUseID, toolName, data)
	case "get_lat_long":
		message, err = locationTool.Invoke(toolUseID, toolName, data)
	default:
		return fmt.Errorf("unsupported tool: %s", toolName)
	}

	if err != nil {
		return fmt.Errorf("error invoking tool: %v", err)
	}

	*messages = append(*messages, *message)
	return nil
}

func parseJSONParameter(parameter string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	err := json.Unmarshal([]byte(parameter), &data)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON parameter: %v", err)
	}
	return data, nil
}


// handleToolUse 处理工具使用情况
func handleToolUse(output types.ConverseOutput, messages *[]types.Message) {
	switch v := output.(type) {
	case *types.ConverseOutputMemberMessage:
		*messages = append(*messages, v.Value)
		
		for _, item := range v.Value.Content {
			switch d := item.(type) {
			case *types.ContentBlockMemberText:
				
			case *types.ContentBlockMemberToolUse:
				
				if *d.Value.Name == "get_lat_long" {
					data := make(map[string]interface{})
					err := d.Value.Input.UnmarshalSmithyDocument(&data)
					
					if err == nil {
						message, err := locationTool.Invoke(*d.Value.ToolUseId, *d.Value.Name, data)
						if err != nil {
							fmt.Printf("Error invoking tool: %v\n", err)
						} else {
							print(*message)
							*messages = append(*messages,*message)
						}
					}
				}

				if *d.Value.Name == "get_weather" {
					data := make(map[string]interface{})
					err := d.Value.Input.UnmarshalSmithyDocument(&data)
					
					if err == nil {
						message, err := weatherTool.Invoke(*d.Value.ToolUseId, *d.Value.Name, data)
						if err != nil {
							fmt.Printf("Error invoking tool: %v\n", err)
						} else {
							print(*message)
							*messages = append(*messages,*message)
						}
					}
				}
			}
		}
	default:
		fmt.Println("Response is nil or unknown type")
	}
}
