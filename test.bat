@echo off

echo #
echo # Building SFX and bundler

go build -o test/sfx.exe -tags verbose ^
    base/log_verbose.go ^
    base/config.go ^
    base/mode.go ^
    base/pecontent.go ^
    base/uncompress.go ^
    base/run.go ^
    base/main.go

go build -o test/bundler.exe ^
    bundler/util.go ^
    bundler/args.go ^
    bundler/config.go ^
    bundler/compress.go ^
    bundler/main.go

echo #
echo # Bundling SFX...

test\bundler.exe -v ^
    -exe test/sfx.exe ^
    -dir project ^
    -compress 9 ^
    -id co.touchify.testsfx ^
    -version 1.0.0 ^
    -args "--sfx"

echo #
echo # Running SFX...

test\sfx.exe

