# Smartling

A client implementation of the [Smartling Translation API](https://docs.smartling.com/display/docs/Smartling+Translation+API) in Go.

It consists of a library for use in other projects, and a CLI tool.

## Using the Library

You can find documentation at http://godoc.org/github.com/99designs/smartling

```go
import "github.com/99designs/smartling"
client := smartling.NewClient(apiKey, projectId)
client.List(smartling.ListRequest{})
```

## CLI tool

The `smartling` CLI tool provides a familiar unix-like command interface to the Smartling API, as well as providing a `project` command to manage a project's local files.

Install it with `go get github.com/99designs/smartling/cli/smartling`


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
   project      manage local project files
```

### The `project` command

When working in a dev environment, it is not desirable to be clobbering the Smartling project files. So these commands use temporary files and prefixes to allow "pushing" and "pulling" translations as required.

```
COMMANDS:
   status   show the status of the project's local files
   pull     translate local project files using Smartling as a translation memory
   push     upload local project files with new strings, using the git branch or user name as a prefix
   help, h  Shows a list of commands or help for one command
```

Other cool features:
- downloaded translation files are cached for 4 hours in `~/.smartling/cache`
- concurrent operations
- filetypes get detected automatically


### Configuration file

The CLI tool uses a project level config file called `smartling.yml` for configuration.

Example config:
```yaml
apikey: "11111111-2222-3333-4444-555555555555"             # Required
projectid: "666666666"                                     # Required

files:                                                     # Files in the project
  - translations/file1.xlf
  - translations/file2.xlf

fileconfig:                                                # Extra config for translation files
  filetype: "xliff"                                        # Override the detected file type
  parserconfig:
    placeholder_format_custom: "%[^%]+%"
  pullfilepath: "{{ TrimSuffix .Path .Ext }}.{{.Locale}}.{{.Ext}}" # The naming scheme when pulling files
```

## TODO
 - docs
 - tests
 - push should only upload if the file has untranslated strings
 - make more things configurable
  - push prefix
  - cache maxage
  - cache location
 - some error handling could be a bit nicer
