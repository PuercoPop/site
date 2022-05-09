# Per-project org-agenda
# en,emacs,org-mode
# 2020-12-28

For my personal projects I tend to keep a HACKING.org file where I write down
project specific stuff. Adding each of those files to org-agenda-files would be
tiresome plus I would see the information relating to every project each time I
open the agenda. This would be a surefire way to get overwhelmed by the amount
of information. Instead I use the following snippet to showing the project
specific agenda when I'm inside a project that has a =HACKING.org= file.

```elisp
  (defun my/org-agenda-list ()
    (interactive)
    (cl-flet ((find-hacking-file ()
                              (let* ((tld (or (locate-dominating-file default-directory ".git/")
                                              (locate-dominating-file default-directory "HACKING.org")))
                                     (hacking-file (concat  default-directory "HACKING.org")))
                                (and (file-exists-p hacking-file)
                                     hacking-file))))
      (let ((local-agenda-file (find-hacking-file)))
        (if local-agenda-file
            (let* ((lexical-binding nil)
                   (org-agenda-files (list local-agenda-file)))
              (org-agenda-list))
          (org-agenda-list)))))
```
