# Introduction

Solutions for *some* exercices of the [*The Go Programming Language*][gopl] book:

  - [ch2/popcount.go][gh-mb-gopl-ch2/popcount.go],
  [ch2/popcount_test.go][gh-mb-gopl-ch2/popcount_test.go]
  (benchmarks will be done later, in 11.6):
    - 2.3
    - 2.4
    - 2.5
  - [ch3/surface.go][gh-mb-gopl-ch3/surface.go]:
    - 3.1
    - 3.2 (a bit clumsy)
    - 3.4
  - [ch3/mandelbrot.go][gh-mb-gopl-ch3/mandelbrot.go]:
    - 3.5
    - 3.6
  - [ch3/z4-newton.go][gh-mb-gopl-ch3/z4-newton.go]:
    - 3.7
  - [ch6/intset.go][gh-mb-gopl-ch6/intset.go],
  [ch6/intset_test.go][gh-mb-gopl-ch6/intset_test.go]:
    - 6.1
    - 6.2
    - 6.3
    - 6.4
    - 6.5
  - [ch7/wc.go][gh-mb-gopl-ch7/wc.go], [ch7/wc_test.go][gh-mb-gopl-ch7/wc_test.go]:
    - 7.1
    - 7.2
  - [ch7/tree.go][gh-mb-gopl-ch7/tree.go], [ch7/tree_test.go][gh-mb-gopl-ch7/tree_test.go]:
    - 7.3
  - [ch8/clockwall.go][gh-mb-gopl-ch8/clockwall.go],
  [ch8/launch-clocks.sh][gh-mb-gopl-ch8/launch-clocks.sh]:
    - 8.1
  - [ch8/ftpd.go][gh-mb-gopl-ch8/ftpd.go]:
    - 8.2
  - [ch8/netcat3.go][gh-mb-gopl-ch8/netcat3.go],
   [ch8/reverb1.go][gh-mb-gopl-ch8/reverb1.go]:
    - 8.3
  - [ch8/reverb2.go][gh-mb-gopl-ch8/reverb2.go]:
    - 8.4
  - [ch8/parallel-mandelbrot.go][gh-mb-gopl-ch8/parallel-mandelbrot.go],
  [ch8/run-parallel-mandelbrot.sh][gh-mb-gopl-ch8/run-parallel-mandelbrot.sh]:
    - 8.5
  - [ch8/depth-limited-crawler.go][gh-mb-gopl-ch8/depth-limited-crawler.go]:
    - 8.6
  - [ch8/mirror.go][gh-mb-gopl-ch8/mirror.go]:
    - 8.7
  - [ch8/reverb2-with-timeout.go][gh-mb-gopl-ch8/reverb2-with-timeout.go]:
    - 8.8
  - [ch8/du.go][gh-mb-gopl-ch8/du.go]
    - 8.9
  - [ch8/depth-limited-crawler-with-cancel.go][gh-mb-gopl-ch8/depth-limited-crawler-with-cancel.go]:
    - 8.10
  - [ch8/fetch.go][gh-mb-gopl-ch8/fetch.go]:
    - 8.11
  - [ch8/chat.go][gh-mb-gopl-ch8/chat.go]:
    - 8.12
    - 8.13
    - 8.14
    - 8.15
  - [ch9/bank1.go][gh-mb-gopl-ch9/bank1.go]:
    - 9.1
  - [ch9/popcount.go][gh-mb-gopl-ch9/popcount.go]:
    - 9.2
  - [ch9/memo.go][gh-mb-gopl-ch9/memo.go] **UNTESTED / WIP**:
    - 9.3
  - [ch9/maxgo.go][gh-mb-gopl-ch9/maxgo.go],
  [ch9/run-maxgo.sh][gh-mb-gopl-ch9/run-maxgo.sh]:
    - 9.4
  - [ch9/pingpong.go][gh-mb-gopl-ch9/pingpong.go],
  [ch9/run-pingpong.sh][gh-mb-gopl-ch9/run-pingpong.sh]:
    - 9.5
  - [ch9/run-parallel-mandelbrot.sh][gh-mb-gopl-ch9/run-parallel-mandelbrot.sh]:
    - 9.6
  - [ch10/jpeg.go][gh-mb-gopl-ch10/jpeg.go]:
    - 10.1
  - [ch12/display.go][gh-mb-gopl-ch12/display.go]:
    - 12.1
    - 12.2
  - [ch12/sexp.go][gh-mb-gopl-ch12/sexp.go]:
    - 12.3
    - 12.4
    - 12.6
  - [ch12/json.go][gh-mb-gopl-ch12/json.go],
  [ch12/json_test.go][gh-mb-gopl-ch12/json_test.go]:
    - 12.5

**<u>Quick book review:</u>** The books feels great; in particular:

  - interesting exercises, close to real-world, varied (e.g.
  Exercise 8.2: implementing a concurrent FTP server);
  - professional/pragmatic code/solutions (e.g. "no need to parse
  your XML as a tree if what you want to do with it can be achieved by
  parsing it as a list of tokens");
  - real-world code used as example (e.g. Go standard library);
  - complete nicely Go's documentation (e.g. detailling the behavior of the
  ``http.HandlerFunc(x)`` *type conversion* (not function call)).

Not for absolute beginners, as stated in the introduction, buf very
good for CS students, or people who got programming skills as they go,
but who lack on the engineering side.

<!--

3.3:
	probably, compute the derivative to determine peeks & valleys;
	not sure how to get the proportionality correct from there to grab
	correct blue<->red gradients (isn't the derivative too local?)
	Maybe there's a clever approach, e.g. from the deformation of the
	polygons?

3.8
	numerical work with benchmark, p82

3.9
	easy

Eventually:
	7.4 / 7.5 : p194

	p204: subtle bits regarding interfaces containing a nil pointer;
	would be nice to clarify all those things with proper "memory diagrams".

10.2 (p307) easy

10.3 p312

10.4 p319

11.1, 11.2 p326, writing tests
11.3, 11.4 p327 again, just about writing tests
11.5 p336 more tests

11.6 11.7 benchmarks, for resp. 2.4/2.5 and 6.1 to 6.5 (IntSet)
	(todo)

12.3 -> 12.7 p360
12.8 12.9 12.10 p366
12.11 12.12 12.13 p369/370
	reflections; at least some of them

13.1 13.2 p380 unsafe

13.3 13.4 p385: theoretically still unsafe

-->

[gopl]: https://www.gopl.io/

[gh-mb-gopl-ch2/popcount.go]: https://github.com/mbivert/gopl/blob/master/ch2/popcount.go
[gh-mb-gopl-ch2/popcount_test.go]: https://github.com/mbivert/gopl/blob/master/ch2/popcount_test.go

[gh-mb-gopl-ch3/surface.go]: https://github.com/mbivert/gopl/blob/master/ch3/surface.go

[gh-mb-gopl-ch3/mandelbrot.go]: https://github.com/mbivert/gopl/blob/master/ch3/mandelbrot.go

[gh-mb-gopl-ch3/z4-newton.go]: https://github.com/mbivert/gopl/blob/master/ch3/z4-newton.go

[gh-mb-gopl-ch6/intset.go]: https://github.com/mbivert/gopl/blob/master/ch6/intset.go
[gh-mb-gopl-ch6/intset_test.go]: https://github.com/mbivert/gopl/blob/master/ch6/intset_test.go

[gh-mb-gopl-ch7/wc.go]: https://github.com/mbivert/gopl/blob/master/ch7/wc.go
[gh-mb-gopl-ch7/wc_test.go]: https://github.com/mbivert/gopl/blob/master/ch7/wc_test.go

[gh-mb-gopl-ch7/tree.go]: https://github.com/mbivert/gopl/blob/master/ch7/tree.go
[gh-mb-gopl-ch7/tree_test.go]: https://github.com/mbivert/gopl/blob/master/ch7/tree_test.go

[gh-mb-gopl-ch8/clockwall.go]: https://github.com/mbivert/gopl/blob/master/ch8/clockwall.go
[gh-mb-gopl-ch8/launch-clocks.sh]: https://github.com/mbivert/gopl/blob/master/ch8/launch-clocks.sh

[gh-mb-gopl-ch8/ftpd.go]: https://github.com/mbivert/gopl/blob/master/ch8/ftpd.go

[gh-mb-gopl-ch8/netcat3.go]: https://github.com/mbivert/gopl/blob/master/ch8/netcat3.go
[gh-mb-gopl-ch8/reverb1.go]: https://github.com/mbivert/gopl/blob/master/ch8/reverb1.go

[gh-mb-gopl-ch8/reverb2.go]: https://github.com/mbivert/gopl/blob/master/ch8/reverb2.go

[gh-mb-gopl-ch8/parallel-mandelbrot.go]: https://github.com/mbivert/gopl/blob/master/ch8/parallel-mandelbrot.go
[gh-mb-gopl-ch8/run-parallel-mandelbrot.sh]: https://github.com/mbivert/gopl/blob/master/ch8/run-parallel-mandelbrot.sh

[gh-mb-gopl-ch8/depth-limited-crawler.go]: https://github.com/mbivert/gopl/blob/master/ch8/depth-limited-crawler.go

[gh-mb-gopl-ch8/mirror.go]: https://github.com/mbivert/gopl/blob/master/ch8/mirror.go

[gh-mb-gopl-ch8/reverb2-with-timeout.go]: https://github.com/mbivert/gopl/blob/master/ch8/reverb2-with-timeout.go

[gh-mb-gopl-ch8/du.go]: https://github.com/mbivert/gopl/blob/master/ch8/du.go

[gh-mb-gopl-ch8/depth-limited-crawler-with-cancel.go]: https://github.com/mbivert/gopl/blob/master/ch8/depth-limited-crawler-with-cancel.go

[gh-mb-gopl-ch8/fetch.go]: https://github.com/mbivert/gopl/blob/master/ch8/fetch.go

[gh-mb-gopl-ch8/chat.go]: https://github.com/mbivert/gopl/blob/master/ch8/chat.go

[gh-mb-gopl-ch9/bank1.go]: https://github.com/mbivert/gopl/blob/master/ch9/bank1.go

[gh-mb-gopl-ch9/popcount.go]: https://github.com/mbivert/gopl/blob/master/ch9/popcount.go

[gh-mb-gopl-ch9/memo.go]: https://github.com/mbivert/gopl/blob/master/ch9/memo.go

[gh-mb-gopl-ch9/maxgo.go]: https://github.com/mbivert/gopl/blob/master/ch9/maxgo.go
[gh-mb-gopl-ch9/run-maxgo.sh]: https://github.com/mbivert/gopl/blob/master/ch9/run-maxgo.sh

[gh-mb-gopl-ch9/pingpong.go]:  https://github.com/mbivert/gopl/blob/master/ch9/pingpong.go
[gh-mb-gopl-ch9/run-pingpong.sh]: https://github.com/mbivert/gopl/blob/master/ch9/run-pingpong.sh

[gh-mb-gopl-ch9/run-parallel-mandelbrot.sh]: https://github.com/mbivert/gopl/blob/master/ch9/run-parallel-mandelbrot.sh

[gh-mb-gopl-ch10/jpeg.go]: https://github.com/mbivert/gopl/blob/master/ch10/jpeg.go

[gh-mb-gopl-ch12/display.go]: https://github.com/mbivert/gopl/blob/master/ch12/display.go

[gh-mb-gopl-ch12/sexp.go]: https://github.com/mbivert/gopl/blob/master/ch12/sexp.go

[gh-mb-gopl-ch12/json.go]: https://github.com/mbivert/gopl/blob/master/ch12/json.go
[gh-mb-gopl-ch12/json_test.go]: https://github.com/mbivert/gopl/blob/master/ch12/json_test.go
