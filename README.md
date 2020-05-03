## Mikbak

Mikbak is a backup tool for mikrotik routeros.

## Install

go build:

```bash
go get -u github.com/wsvn53/mikbak
```

or, download from releases page.

## Usage

First, you need enable ssh service from RouterOS: [IP] -> [Sevices] -> [Enable SSH].

Command help:

```shell
Usage:
  mikbak-darwin [OPTIONS]

Application Options:
  -s, --server=   Mikrotik RouterOS ip address, example: 127.0.0.1
  -u, --user=     Username to login.
  -p, --password= Password to login, support base64 encoded password with prefix 'B:'.
  -o, --output=   Target directory to save backup file, by default will save to current directory. (default: .)
      --prefix=   Add prefix to backup filename. (default: ros)

Help Options:
  -h, --help      Show this help message
```

Example:

```bash
./mikbak -s 192.168.1.1 -u admin -p 123456 -o /local/backups
```

For more security, you can using base64 encoded password:

```bash
./mikbak -s 192.168.1.1 -u admin -p B:MTIzNDU2Cg== -o /local/backups
```

Base64 encoded password must be appended with suffix "B:"!

Custom filename:

```bash
./mikbak -s 192.168.1.1 -u admin -p B:MTIzNDU2Cg== --prefix=myros
```

will save to [myros-2020xxxx.backup].

## License

MIT.