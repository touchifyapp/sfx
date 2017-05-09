@echo off

go build -o test/sfx.exe -tags verbose ^
    base/log_verbose.go ^
    base/config.go ^
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

echo Bundling SFX...
test\bundler.exe -v ^
    -exe test/sfx.exe ^
    -dir project ^
    -compress 9 ^
    -id co.touchify.testsfx ^
    -version 2.0.0 ^
    -args "--sfx"

echo Running SFX...
test\sfx.exe

