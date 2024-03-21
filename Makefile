clear:
	@rm -rf output
	@mkdir output

build: clear
	@go build -o output/diary app/diary/cmd/main.go

run: build
	@cd output && ./diary
