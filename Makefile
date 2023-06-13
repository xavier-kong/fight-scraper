build:
	GOOS=linux GOARCH=amd64 go build -o bin/main main.go && zip bin/main.zip bin/main  

deploy:
	aws lambda update-function-code --function-name fightScraper --zip-file fileb://bin/main.zip

