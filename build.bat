@echo off

echo #
echo # Building SFX...

go build -o bin/sfx.exe ^
    -ldflags "-H windowsgui" ^
    base/log_silent.go ^
    base/config.go ^
    base/mode.go ^
    base/pecontent.go ^
    base/uncompress.go ^
    base/run.go ^
    base/main.go

echo #
echo # Building SFX (verbose)...

go build -o bin/sfxv.exe ^
    -tags verbose ^
    base/log_verbose.go ^
    base/mode.go ^
    base/config.go ^
    base/pecontent.go ^
    base/uncompress.go ^
    base/run.go ^
    base/main.go


echo #
echo # Building Bundler...

go build -o bin/bundler.exe ^
    bundler/util.go ^
    bundler/args.go ^
    bundler/config.go ^
    bundler/compress.go ^
    bundler/main.go
