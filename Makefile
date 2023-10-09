objects = src/main.go

blog-editor = src/blog-editor.go

run: $(objects)
	@go build $(objects)
	@./main test 8000

run-edit: $(blog-editor)
	@go run $(blog-editor)
