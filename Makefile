binary = \
	bin/linux-amd64-edit-slack \
	bin/windows-amd64-edit-slack.exe \
	bin/darwin-amd64-edit-slack

all: $(binary)

bin/linux-amd64-edit-slack: main.go editslack/*
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $@

bin/windows-amd64-edit-slack.exe: main.go editslack/*
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o $@

bin/darwin-amd64-edit-slack: main.go editslack/*
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o $@

clean:
	rm $(binary)
