@echo off

bash -version >nul 2>&1 && set USEBASH=1

if [%USEBASH%]==[1] (
    bash ./build.sh
    exit /b 0
) 

echo #
echo # Bash is not installed, buiding from Windows...
echo #
echo # Building SFX...

go generate ./base
go build -o bin/sfx.exe ^
    -ldflags "-s -w -H windowsgui" ^
    ./base

echo #
echo # Building SFX (verbose)...

go generate ./base
go build -o bin/sfxv.exe ^
    -ldflags "-s -w" ^
    -tags verbose ^
    ./base


echo #
echo # Building Bundler...

go build -o bin/bundler.exe ^
    ./bundler