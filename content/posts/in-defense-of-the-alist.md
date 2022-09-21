# In Defense of the Alist
## lisp, en
## 2018-05-24

Lisp is a well designed language, but sometimes its good design comes off as
cruft at first sight. This is the case with [association lists] (alists from
now on). Alists are used in Lisp programs in places where more recent languages
(Python, Ruby, JS) would use hash tables instead.  Not having first class
syntax for maps seems like a misdesign at first. And although adding first
class syntax is a mere [~30 lines] of [code away], I posit that not
having first class is an example of good design — in particular regarding
developer 'ergonomics'. If one can write the elements of the map in the source
code, you shouldn't be using a map in the first place!

The efficiency of algorithms is normally compared using Big-O notation. Let’s
not forget that Big-O notation is a way to analyze performance in general
cases, not in concrete ones. That is to say O(1) is better than O(n) as a
_general_ rule. However this is not case for small values of n; in the case of
hash maps in Lisp the threshold is around [25 elements]. Therefore by not
providing literal syntax for maps, Lisp ensures that the ergonomics favour the
better choice for representation.

There is however another virtue of alists, one that is related to an well known
aphorism in software development: “Delay the decision up to the last possible
moment”. In the case of alists that is that the retrieval strategy used by the
operation is decided at the call-site, not set in stone when the data structure
is created. This is important as it lets one, for example, change the equality
predicate of keys. A more interesting way to exploit this characteristic would
be to turn the alist into an most recently used (MRU) cache:

```lisp
  (defpackage "ALIST-MRU"
    (:use "CL"))
  (in-package "ALIST-MRU")

  (defun mru-assoc (item alist &key (key 'car) (test 'eql))
    (labels ((%mru-assoc (item alist key test head)
               (cond ((endp alist) nil)
                     ((funcall test item (funcall key (car alist)))
                      ;; =>
                      (progn
                        (let ((head-cons (car head)))
                          (rplaca head (car alist))
                          (rplaca alist head-cons))
                        (car alist)))
                     (t (%mru-assoc item (cdr alist) key test head)))))
      (%mru-assoc item alist key test alist)))
```

The gist of the above code is to save a reference to the first cons cell and
upon finding a successful match swapping the car of the first cons cell with
that of the successful match and updating the neighbor of the point to the cons
cell that the successful match used to point to.

[association lists]:http://www.lispworks.com/documentation/HyperSpec/Body/26_glo_a.htm#association_list
[~30 lines]: http://frank.kank.net/essays/hash.html
[code away]: https://github.com/vseloved/rutils/blob/db3c3f4ae897025b5f0cd81042ca147da60ca0c5/core/readtable.lisp#L28-L51
[25 elements]: http://funcall.blogspot.pe/2016/01/alist-vs-hash-table.html
