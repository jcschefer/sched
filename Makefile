install:
	go build sched.go
	mv sched /usr/bin/

uninstall:
	rm /usr/bin/sched
