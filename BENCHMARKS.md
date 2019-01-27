Prior to transposition table.

```bash
.[~/go/src/github.com/rictorlome/pupil]    (master)
$ go test -bench=.
goos: darwin
goarch: amd64
pkg: github.com/rictorlome/pupil
BenchmarkAlphaBetaInitial2-4    	    3000	    405515 ns/op
BenchmarkAlphaBetaInitial3-4    	     200	   7661683 ns/op
BenchmarkAlphaBetaInitial4-4    	      10	 112742338 ns/op
BenchmarkAlphaBetaInitial5-4    	       2	 906786002 ns/op
BenchmarkAlphaBetaKiwipete2-4   	    2000	    817643 ns/op
BenchmarkAlphaBetaKiwipete3-4   	     200	   9855421 ns/op
BenchmarkAlphaBetaKiwipete4-4   	      20	  93296598 ns/op
BenchmarkAlphaBetaKiwipete5-4   	       1	1077769417 ns/op
PASS
ok  	github.com/rictorlome/pupil	15.649s
```

After adding transposition table lookups at static_evaluate level

```bash
goos: darwin
goarch: amd64
pkg: github.com/rictorlome/pupil
BenchmarkAlphaBetaInitial2-4    	   10000	    108712 ns/op
BenchmarkAlphaBetaInitial3-4    	    1000	   2213521 ns/op
BenchmarkAlphaBetaInitial4-4    	      50	  38574630 ns/op
BenchmarkAlphaBetaInitial5-4    	       5	 305574264 ns/op
BenchmarkAlphaBetaKiwipete2-4   	   10000	    195509 ns/op
BenchmarkAlphaBetaKiwipete3-4   	     500	   2624494 ns/op
BenchmarkAlphaBetaKiwipete4-4   	      50	  29009556 ns/op
BenchmarkAlphaBetaKiwipete5-4   	       1	1182648464 ns/op
PASS
ok  	github.com/rictorlome/pupil	17.526s
```

After searching best move first

```bash
goos: darwin
goarch: amd64
pkg: github.com/rictorlome/pupil
BenchmarkAlphaBetaKiwipete2-4   	    5000	    350846 ns/op
BenchmarkAlphaBetaKiwipete3-4   	    1000	   1181788 ns/op
BenchmarkAlphaBetaKiwipete4-4   	     100	  18995218 ns/op
BenchmarkAlphaBetaKiwipete5-4   	      30	  45606258 ns/op
BenchmarkAlphaBetaInitial2-4    	   50000	     37440 ns/op
BenchmarkAlphaBetaInitial3-4    	   10000	    147593 ns/op
BenchmarkAlphaBetaInitial4-4    	    2000	    978088 ns/op
BenchmarkAlphaBetaInitial5-4    	     300	   4503095 ns/op
PASS
ok  	github.com/rictorlome/pupil	14.987s
```
