set GOOS=freebsd
go build -o bin/freebsd/check_unifi
set GOOS=linux
go build -o bin/linux/check_unifi
