mv profiles/*_test.go .
tmpfile=$(mktemp /tmp/cpuprof.XXXXXX)
go test -v -cpuprofile $tmpfile
mv *_test.go profiles/
go tool pprof -http=":" $tmpfile
rm $tmpfile
