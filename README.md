Packer is a simplified implementation of a package manager.  
It includes a command-line program for managing project libraries, as well as a server program for managing repositories.  

## pack.yml
```yml
registry: 'http://127.0.0.1:3000'
name: 'packer'
version: 0.1.0

dependencies:
  - packer-utils@0.1.0

# Executed by the quick.js runtime built into packer
# You can write any script based on the quick.js engine
scripts:
  # pack run install
  install: './pack_scripts/install.mjs'
  serve: './pack_scripts/serve.mjs'
  bundle: './pack_scripts/bundle.mjs'

hooks:
  # Callback after executing pack sync
  after_sync: 'pack_scripts/after_sync.mjs'
```
## pack_libs
The source code will be saved in `~/.pack/libs`.
Create a `pack_libs` dir in the root directory of the project and symlink it to `~/.pack/libs`.

## bin

```bash
# update all packages
pack sync

# add new package
pack add <package>

# delete
pack rm <package>

# publish
pack publish
```

## server
```bash
pack server
```