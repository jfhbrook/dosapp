## dosapp

- [ ] create folder structure I want
- [X] create bash script that wraps go-task
- [X] install gomplate
- [X] "init" command that templates out ~/.config/dosapp
  - mkdir -p
  - copy Taskfile.yml
  - run init task
- [X] stub wordperfect package
  - [X] `./packages/{{app}}/Taskfile.yml`
  - [ ] `./packages/{{app}}/dosapp.env.tmpl`
    - location of wordperfect root drive
    - location of documents
  - [ ] `./packages/{{app}}/*.conf.tmpl`
    - use wordperfect root drive
    - use documents location
- [ ] "install" command
  - template out `~/.local/share/dosapp/apps/{{app}}/Taskfile.yml`
  - template out `~/.local/share/dosapp/apps/{{app}}/*.conf`
    - use new `Taskfile.yml` to run this task
  - run "install" task
    - `./packages/{{app}}/Taskfile.yml`
    - store downloads in `~/.cache/dosapp/downloads`
    - store disks in `~/.local/share/dosapp/disks`
- [ ] template out ~/.config/dosapp/apps/{{app}} with gomplate
  - and run that, not packages/{{app}}/Taskfile.yml etc
- [ ] 'link' command that templates out `~/.local/bin/*` with `packages/{{app}}/bin/*.tmpl`
- [ ] 'publish' command that just calls `~/.config/dosapp/Taskfile.yml`

That will give me:

- user volumes
  - `joshiverse/dos/disks`
  - `~/Documents`
- commands
  - `dosapp init`
  - `dosapp install {{app}}`
  - `dosapp start {{app}}`
  - `dosapp link {{app}}`
  - `dosapp publish`
- config files
  - `~/.config/dosapp/Taskfile.yml`
  - `~/.config/dosapp/dosapp.env`
  - `~/.config/dosapp/main.conf.tmpl`
  - `~/.config/dosapp/apps/{{app}}/Taskfile.yml`
  - `~/.config/dosapp/apps/{{app}}/*.conf.tmpl`
- cache files
  - `~/.cache/dosapp/packages/{{app}}`
    - (in `joshiverse/packages` for now)
  - `~/.cache/dosapp/downloads`
- state files
  - `~/.local/state/dosapp/main.conf`
  - `~/.local/state/dosapp/apps/{{app}}/*.conf`

This is a LOT!!

If I get that working with go-task and gomplate:

- Create cobra app with the same shape as the bash command
- use `exec.Command` to call go-task and gomplate
- create `dosapp template` command that replaces gomplate
- create `dosapp task` command that embeds go-task

For packages:

- create/push a public repo
- `dosapp publish` command that shells out tar and `gh release`
- `dosapp fetch` command that pulls/unpacks tars from gh releases

## app fun LOL

- install turbo pascal
