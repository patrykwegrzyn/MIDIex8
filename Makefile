build:
	go build -o bin/midiex8 main.go

run:
	go run main.go

install:
	go run main.go

install:
	echo "Installing"
	cp bin/midiex8 /usr/bin
	cp midiex8.service /etc/systemd/system
	systemctl daemon-reload
	systemctl start midiex8.service
