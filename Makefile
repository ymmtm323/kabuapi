all: build install

build:
	CGO_ENABLED=0 go build -tags osusergo,netgo -o cmd/kabuapi/kabuapi cmd/kabuapi/main.go

install:
	install -m 0755 cmd/kabuapi/kabuapi /usr/local/bin/kabuapi
	install -m 0644 misc/kabuapi.service /etc/systemd/system/kabuapi.service
	systemctl daemon-reload
	systemctl enable kabuapi
	systemctl restart kabuapi

uninstall:
	systemctl stop kabuapi
	systemctl disable kabuapi
	rm -f /etc/systemd/system/kabuapi.service
	systemctl daemon-reload
	rm -f /usr/local/bin/kabuapi

clean:
	rm -f cmd/kabuapi/kabuapi