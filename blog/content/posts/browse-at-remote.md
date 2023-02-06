# browse-at-remote snippet
## en, emacs, protip
## 2020-10-12

```elisp
(use-package browse-at-remote
  :after (magit)
  :ensure t
  :config
  (transient-append-suffix 'magit-file-dispatch "m"
    '("o" "Browse file" browse-at-remote)
    (transient-replace-suffix 'magit-dispatch "o"
      '("o" "Browse file" browse-at-remote))))
```
