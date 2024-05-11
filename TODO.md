# dosapp

## v1.0

- [ ] Fill out help text
- [ ] Add `--version` flag
- [ ] Write README
- [ ] Turn repository public

## Packages

- [ ] turbo-pascal
- [ ] DOStodon

## Go Rewrite

- [ ] Create cobra app with the same shape as the bash command
- [ ] Use `exec.Command` to call existing bash procedures
- [ ] Load dotenv file in go
- [ ] Pull logic from bash procedures to go - call `task` and `gomplate` with
      `exec.Command`
- [ ] Create `template` command that replaces gomplate
- [ ] Create `task` command that embeds go-task
- [ ] Create `unpack` and `pack` commands that use `sevenzip`, etc
- [ ] Create `download` command that replaces `curl`
- [ ] Publish packages in the form of github releases
- [ ] Create `fetch` command that fetches a release from github
