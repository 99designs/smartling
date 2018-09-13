# Smartling

A CLI tool implementation of the [Smartling Translation API](https://developer.smartling.com/docs/list-of-smartling-apis) in Go.

## Using the CLI tool

The `smartling` CLI tool provides a familiar unix-like command interface to the Smartling API, as well as providing a `project` command to manage a project's local files.

Install it with `go get github.com/99designs/smartling/cli/smartling`
*or*
run it as a docker container `docker run --rm -v MyProject:/work 99designs/smartling ls`


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


### The `smartling project` command

The `smartling project` commands are designed for some common use-cases in a dev or CI environment.

```
COMMANDS:
   files  lists the local files
   status show the status of the project's remote files
   pull   translate local project files using Smartling as a translation memory
   push   upload local project files that contain untranslated strings
```

"Pushing" uploads files to a smartling project using a prefix. By default it uses the git branch name , but you can also specifiy the wanted prefix as an argument. A hash is also used in the prefix to prevent clobbering.

"Pulling" translates local project files using Smartling as a translation memory.

Other cool features:
- downloaded translation files are cached (default is 4 hours) in `~/.smartling/cache`
- operations mostly happen concurrently
- filetypes get detected automatically


### Configuration file

The CLI tool uses a project level config file called `smartling.yml` for configuration.

Example config:
```yaml
# Required config
ApiKey: "aaaaaabbbbbbbbcccccddddddd"                       # Smartling API Token Secret token
UserID: "a1b2c3d4e5f6"                                     # Smartling User Identifier
ProjectId: "666666666"                                     # Smartling Project Id
Files:                                                     # Files in the project
  - translations/*.xlf                                     # Globbing can be used,
  - foo/bar.xlf                                            # as well as individual files

# Optional config
CacheMaxAge: "4h"                                          # How long to cache translated files for
FileType: "xliff"                                          # Override the detected file type
ParserConfig:                                              # Add a custom configuration
  placeholder_format_custom: "%[^%]+%"
PullFilePath: "{{ TrimSuffix .Path .Ext }}.{{.Locale}}{{.Ext}}" # The naming scheme when pulling files
```
