## Simple Gist Creation Tool

### About

This tool reads from stdin until EOF, and creates a gist on github with the content.

### Usage

```bash
$ cat foo.txt | gist
https://gist.github.com/abcdefg
```

### Configuration

Create a `~/.gistrc` file with contents:

```json
{
    "key": "<Your API Key Goes Here>",
    "user": "github username, defaults to the current shell user",
    "apiUrl": "github api url, defaults to https://api.github.com"
}
```

For github enterprise, it may be necessary to specify baseUrl `https://[host]/api/v3`
