# Only GNU Make is supported

meseircd: *.go
	go build -o $@
