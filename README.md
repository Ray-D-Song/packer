Packer is a simplified implementation of a package manager.  
It includes a command-line program for managing project libraries, as well as a server program for managing repositories.  

## pack.yml

## pack_libs

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
# packer will read config from ~/.packer/config.yml
pack server
```