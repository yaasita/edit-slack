binary = \
	bin/linux-386-edit-slack \
	bin/linux-amd64-edit-slack \
	bin/freebsd-386-edit-slack \
	bin/freebsd-amd64-edit-slack \
	bin/windows-386-edit-slack.exe \
	bin/windows-amd64-edit-slack.exe \
	bin/darwin-386-edit-slack \
	bin/darwin-amd64-edit-slack

all: $(binary)

bin/linux-386-edit-slack: main.go slack/*
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo -o $@

bin/linux-amd64-edit-slack: main.go slack/*
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo -o $@

bin/freebsd-386-edit-slack: main.go slack/*
	GOOS=freebsd GOARCH=386 CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo -o $@

bin/freebsd-amd64-edit-slack: main.go slack/*
	GOOS=freebsd GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo -o $@

bin/windows-386-edit-slack.exe: main.go slack/*
	GOOS=windows GOARCH=386 CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo -o $@

bin/windows-amd64-edit-slack.exe: main.go slack/*
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo -o $@

bin/darwin-386-edit-slack: main.go slack/*
	GOOS=darwin GOARCH=386 CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo -o $@

bin/darwin-amd64-edit-slack: main.go slack/*
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo -o $@

clean:
	rm $(binary)
