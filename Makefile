tidy:
	@echo Checking dependencies...
	@go mod tidy -v

win: tidy
	@echo Building for windows...
	@GOOS=windows GOARCH=386 go build -o xml2csv-parser.exe ./

mac: tidy
	@echo Building for mac...
	@GOOS=darwin GOARCH=amd64 go build -o xml2csv-parser ./

linux: tidy
	@echo Building for linux...
	@GOOS=linux GOARCH=amd64 go build -o xml2csv-parser ./