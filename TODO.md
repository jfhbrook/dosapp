# dosapp

## v1.0

- [ ] 'overwrite' behavior for `dosapp init` and Taskfile/conf
- [ ] 'link' command that templates out `~/.local/bin/*` with `packages/{{app}}/bin/*.tmpl`
- [ ] 'publish' task in wordperfect package that publishes a release to github
- [ ] 'fetch' command in `./bin/dosapp` that fetches a release from github

## Packages

- [ ] turbo-pascal
- [ ] DOStodon

## Go Rewrite

- Create cobra app with the same shape as the bash command
- use `exec.Command` to call go-task and gomplate
- create `dosapp template` command that replaces gomplate
- create `dosapp task` command that embeds go-task
