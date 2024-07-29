# Introduction

Solutions for *some* exercices of the [*The Go Programming Language*][gopl] book:

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

Eventually:
	7.4 / 7.5 : p194

	p204: subtle bits regarding interfaces containing a nil pointer;
	would be nice to clarify all those things with proper "memory diagrams".

	8.1 / 8.2 : p241
-->

[gopl]: https://www.gopl.io/

[gh-mb-gopl-ch6/intset.go]: https://github.com/mbivert/gopl/blob/master/ch6/intset.go
[gh-mb-gopl-ch6/intset_test.go]: https://github.com/mbivert/gopl/blob/master/ch6/intset_test.go

[gh-mb-gopl-ch7/wc.go]: https://github.com/mbivert/gopl/blob/master/ch7/wc.go
[gh-mb-gopl-ch7/wc_test.go]: https://github.com/mbivert/gopl/blob/master/ch7/wc_test.go

[gh-mb-gopl-ch7/tree.go]: https://github.com/mbivert/gopl/blob/master/ch7/tree.go
[gh-mb-gopl-ch7/tree_test.go]: https://github.com/mbivert/gopl/blob/master/ch7/tree_test.go
