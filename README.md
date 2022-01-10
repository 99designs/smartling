# Smartling

A CLI tool implementation of the [Smartling Translation API](https://developer.smartling.com/docs/list-of-smartling-apis) in Go.

## Using the CLI tool

The `smartling` CLI tool provides a familiar unix-like command interface to the Smartling API, as well as providing a `project` command to manage a project's local files.

### Installing

Install smartling by downloading the binary from the releases page.

Or you can run it as a docker container e.g. `docker run --rm -v MyProject:/work 99designs/smartling ls`

### Usage

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

Other features:
- downloaded translation files are cached (default is 4 hours) in `~/.smartling/cache`
- operations mostly happen concurrently
- filetypes get detected automatically


### Configuration file

The CLI tool uses a project level config file called `smartling.yml` for configuration.

Example config:
```yaml
# Required config
api_key: "aaaaaabbbbbbbbcccccddddddd"                       # Smartling API Token Secret token
user_id: "a1b2c3d4e5f6"                                     # Smartling User Identifier
project_id: "666666666"                                     # Smartling Project Id
files:                                                     # Files in the project
  - translations/*.xlf                                     # Globbing can be used,
  - foo/bar.xlf                                            # as well as individual files

# Optional config
cache_max_age: "4h"                                          # How long to cache translated files for
file_type: "xliff"                                          # Override the detected file type
parser_config:                                              # Add a custom configuration
  placeholder_format_custom: "%[^%]+%"
pull_file_path: "{{ TrimSuffix .Path .Ext }}.{{.Locale}}{{.Ext}}" # The naming scheme when pulling files
```

### How to make a release

1. Check out the the commit you want to create a release for, and tag it with appropriate semver convention:

```
$ git tag vx.x.x
$ git push --tags
```

2. Create the binaries:

```
$ make clean
$ make release
```

3. Go to https://github.com/99designs/smartling/releases/new

4. Select the tag version you just created

5. Attach the binaries from `./bin/*`

### How to build multi-arch DockerHub image

If building a multi-arch image for the first time, you may need to switch your builder. Checkout Docker guide on [building multi-arch images](https://docs.docker.com/desktop/multi-arch/#build-multi-arch-images-with-buildx)

```
$ docker buildx create --use
```

Build and push to target repository

```
$ docker buildx build --push --tag 99designs/smartling:x.x.x --platform linux/amd64,linux/arm64 .
```
