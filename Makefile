objects = src/*.go

run: $(objects)
	@go build $(objects)
	@./main test 8000
