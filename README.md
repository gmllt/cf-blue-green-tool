# NOT USABLE RIGHT NOW WORK IN PROGRESS

# cf blue-green-tool

**A simple Blue/Green deployment tool based on manifest**

Inspired by [bluemixgaragelondon/cf-blue-green-deploy](https://github.com/bluemixgaragelondon/cf-blue-green-deploy)

## Installation

Clone repository, build, install as cf plugin.

```bash
$ mkdir -p $GOPATH/src/github.com/gmllt
$ git clone https://github.com/gmllt/cf-blue-green-tool.git $GOPATH/src/github.com/gmllt/cf-blue-green-tool
$ cd $GOPATH/src/github.com/gmllt/cf-blue-green-tool
$ go build
$ cf install-plugin cf-blue-green-tool
```

## Usage

### Deploy green versions
```bash
$ cf bgt green
```

### Rollback
```bash
$ cf bgt rollback
```

### Approve green versions
```bash
$ cf bgt approve --delete-old-apps
```

### Deploy with blue green
```bash
$ cf bgt deploy --delete-old-apps
```
