mv profiles/*_test.go .
tmpfile=$(mktemp /tmp/cpuprof.XXXXXX)
go test -v -cpuprofile $tmpfile
mv *_test.go profiles/
~/go/bin/pprof -http=":" $tmpfile
rm $tmpfile
