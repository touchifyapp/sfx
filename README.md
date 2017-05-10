# sfx (Self Extracting Archive) [![Build Status](https://travis-ci.org/touchifyapp/sfx.png)](https://travis-ci.org/GeertJohan/go.rice)

`sfx` is [Go](http://golang.org/) package that allows creating __Self Extracting Archive (sfx)__ for Windows.

## Usage

This package provides a bundler which append an archive into an sfx bootstrapper.

```cmd
$ copy sfx.exe dest/myprogram.exe
$ bundler.exe -exe dest/myprogram.exe -dir myarchivedir -compress 9 -id myprogramid
```

## SFX bootstrappers

This package provides two SFX bootstrappers:
 * `sfx.exe`: Silent SFX boostrapper
 * `sfxv.exe`: Verbose SFX bootstrapper

## Bundler options

Run `bundler -h` to print usage help :

| Option        | Type      | Description                                                           |
|---------------|-----------|-----------------------------------------------------------------------|
| -args         | string    | arguments to pass to executable                                       |
| -compress     | int       | The program to run in the project directory.                          |
| -dest         | string    | The absolute destination path to extract project in (default: temp).  |
| -dir          | string    | The directory to bundle into sfx. (default "project")                 |
| -exe          | string    | The program to bundle the project in. (default "sfx.exe")             |
| -id           | string    | The unique ID for this package. (default "co.touchify.sfx")           |
| -run          | string    | The program to run in the project directory (default: auto-detect).   |
| -v            | boolean   | Enable program output.                                                |
| -version      | string    | The program version to check for updates. (default `1.0.0`)           |

## License

[MIT](https://github.com/touchifyapp/sfx/blob/master/LICENSE)

## Changelog

 * `1.0.0`: Initial release
 * `1.0.1`:
     * Append _Microsoft Windows File Properties/Version Info_ to `sfx.exe` and `sfxv.exe` using [goversioninfo](https://github.com/josephspurrier/goversioninfo).
     * This allows _File Properties/Version Info_ modification via tools like [rcedit](https://github.com/electron/rcedit).