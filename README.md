# third-eye

Command to encrypt / decrypt file name.

## Usage

```console
$ go build
```

```console
$ ./third-eye -key <key> <path>
$ ./third-eye -d -key <key> <path>
```

or

```console
$ echo $KEY | ./third-eye <path>
$ echo $KEY | ./third-eye -d <path>
```
