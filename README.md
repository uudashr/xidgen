# xidgen

Xidgen is a cli tool to generate and decode [xid](https://github.com/rs/xid).

## Installation

```shell
go install github.com/kechako/xidgen@latest
```

## Usage

```
$ xidgen -help
Usage of xidgen:
  -decode string
        Decode xid
  -format value
        Output format [hex, binary] (default hex)
  -n int
        Generate n xid(s) (default 1)
  -o string
        Output file
  -passthru
        Passthru mode
  -separator string
        Separator (default "\n")
  -v    Verbose output
  -validate string
        Validate xid
```

### Generate xid

```
$ xidgen
cuf36vlp6umi15dn3130
```

### Generate multiple xid

```
$ xidgen -n 3
cuf373tp6umi1h3kc79g
cuf373tp6umi1h3kc7a0
cuf373tp6umi1h3kc7ag
```

### Output format

```
$ xidgen -format binary|xxd
00000000: 679e 33b7 b937 ad21 0403 162b            g.3..7.!...+
```

### Output to file

```
$ xidgen -o xid.txt
$ xidgen -n 3 -o xids.txt
```

### Decode

```
$ xidgen -decode cuf373tp6umi1h3kc79g
Timestamp:   2025-02-01T21:45:35+07:00
Machine ID:  b937ad
Process ID:  8388
Counter:     7627219
```

### Validate

```
$ xidgen -validate cuf373tp6umi1h3kc79g
(exit code 0)
```

```
$ xidgen -validate asdf
Invalid ID
(exit code 1)
```

```
$ echo cuf373tp6umi1h3kc79g | xidgen -validate -
```