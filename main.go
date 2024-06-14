// Copyright (c) 2023 suwei007@gmail.com
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
)

const defaultRegion = "us-east-1"
const defaultModel = "claude3-sonnet"

var modelMap = map[string]string{
	"claude3-sonnet": "anthropic.claude-3-sonnet-20240229-v1:0",
	"claude3-haiku":  "anthropic.claude-3-haiku-20240307-v1:0",
}

var verbose *bool

func main() {

	model := flag.String("model", defaultModel, "The model to use (e.g.,claude3-sonnet, claude3-haiku)")
	stream := flag.Bool("stream", false, "Use streaming")
	region := flag.String("region", defaultRegion, "setup default region")
	verbose = flag.Bool("verbose", false, "setting to true will log messages being exchanged with LLM")

	inputText := flag.String("question", "What's weather in Beijing ?", "your question")
	flag.Parse()

	if os.Getenv("AWS_REGION") != "" {
		region = aws.String(os.Getenv("AWS_REGION"))
	}

	modelID := "anthropic.claude-3-sonnet-20240229-v1:0"

	log.Printf("Amazon Bedrock [AWS_REGION: %s, model: %s, modeID: %s ,stream: %v ]", *region, *model, modelMap[*model], *stream)

	// Check if the provided model is valid
	if _, ok := modelMap[*model]; !ok {
		log.Fatalf("Invalid model: %s", *model)
	}

	//build toolconfiguration
	weatherToolSchema := *weatherTool.GenerateToolSchema()
	locationToolSchema := *locationTool.GenerateToolSchema()

	toolConfig := types.ToolConfiguration{
		Tools: []types.Tool{
			&weatherToolSchema,
			&locationToolSchema,
		},
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(*region))
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

	client := bedrockruntime.NewFromConfig(cfg)

	fmt.Printf("Question: %s\n", *inputText)
	if *stream {
		err = generateTextStream(context.TODO(), client, modelID, *inputText, &toolConfig)
		if err != nil {
			log.Fatalf("failed to generate text, %v", err)
		}
	} else {
		err = generateText(context.TODO(), client, modelID, *inputText, &toolConfig)
		if err != nil {
			log.Fatalf("failed to generate text, %v", err)
		}
	}

	fmt.Printf("Finished generating text with model %s.\n", modelID)

}
