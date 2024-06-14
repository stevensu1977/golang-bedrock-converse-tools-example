# golang-bedrock-converse-tools-example
This is an AWS Go SDK v2 sample project for how to use bedrock converse API (tools). 

We create 2 tools :

* location tool , get latitude and longitude

* weather tool, get weather from latitude and longitude

  


Install :

```bash
git clone https://github.com/stevensu1977/golang-bedrock-converse-tools-example
cd golang-bedrock-converse-tools-example
go get 
go build
./weather-tools


Usage of ./weather-tools:
  -model string
        The model to use (e.g.,claude3-sonnet, claude3-haiku) (default "claude3-sonnet")
  -question string
        your question (default "What's weather in Beijing ?")
  -stream
        Use streaming
  -verbose
        setting to true will log messages being exchanged with LLM
```









### Credit 

https://github.com/xiwan/AWSTools/tree/main/Claude3/Golang
