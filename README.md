# Smartling

A client implementation of the [Smartling Translation API](https://docs.smartling.com/display/docs/Smartling+Translation+API)


## Using the Library

```
import "github.com/99designs/smartling"

client := smartling.NewClient(apiKey, projectId)

client.List(smartling.ListRequest{})

```

## CLI tool

The CLI tool is designed for two purposes.

### API commands

Provide a familiar unix-like command interface to the Smartling API.
```
COMMANDS:
   ls           list remote files
   stat         display the translation status of a remote file
   get          downloads a remote file
   put          uploads a local file
   rename       renames a remote file
   rm           removes a remote file
   lastmodified shows when a remote file was modified last
   locales      list the locales for the project
```

### `project` command

The `project` command helps manage a local project's translation files.

When working in a dev environment, it is not desirable to be overwriting the Smartling project files constantly. So these commands use temporary files to allow pushing and pulling translations as required.

```
COMMANDS:
   status   show the status of the project's local files
   pull     translate the local project files using Smartling as a translation memory
   push     Upload the local project files to Smartling, using the git branch as a prefix
   help, h  Shows a list of commands or help for one command
```

### Configuration file

The CLI tool uses a project level config file called `smartling.yml` for configuration.

Example config:
```
apikey: "11111111-2222-3333-4444-555555555555"
projectid: "666666666"

fileconfig:
  filetype: "xliff"
  parserconfig:
    placeholderformat: "%[^%]+%"
  pullfilepath: "{{.BaseName}}.{{.Locale}}.{{.Extension}}"

files:
  - test.xlf
```
