
test:
	go build -ldflags "-X main.FileName=go_test.mod" && ./gomod-check
	rm ./gomod-check