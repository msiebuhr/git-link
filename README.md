# git-link

A git plugin that generates web-url links for various bits of git-content; namely files and commits:

```sh
% git link
https://github.com/msiebuhr/git-link
% git link HEAD
https://github.com/msiebuhr/git-link/commit/90246ec398933e676d8d940abd85e38a2a565ddf
% git link README.md
https://github.com/msiebuhr/git-link/blob/master/README.md
```

Only extra option is `--open` (or `-open` for Gophers) that opens the link in a browser. This is done by calling `xdg-open`, which I belive is a Linux-only thing.
