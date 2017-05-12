@echo off

set OPT=%1
set VERSION=%2

if ["%OPT%"]==["clean"] set CLEAN=1
if ["%OPT%"]==["clean-src"] set CLEAN=1

if ["%OPT%"]==["src"] set SRC=1
if ["%OPT%"]==["clean-src"] set SRC=1

if ["%OPT%"]==["quick"] set QUICK=1

if [%VERSION%]==[] set VERSION=1.0.0

if defined CLEAN (
    echo #
    echo # Cleaning %TEMP%\co.touchify.testsfx

    rmdir /S /Q "%TEMP%\co.touchify.testsfx"
)

if not defined QUICK (
    echo #
    echo # Building SFX and bundler

    if defined SRC (
        bash ./build.sh
    ) else (
        go build -o test/sfx.exe ^
            -tags verbose ^
            ./base

        go build -o test/bundler.exe ^
            ./bundler
    )
)

echo #
echo # Bundling SFX...

if defined SRC (
    set BUNDLER=bin\x64\bundler.exe
    copy bin\x64\sfxv.exe test\sfx.exe
) else (
    set BUNDLER=test\bundler.exe
)

%BUNDLER% -v ^
    -exe test/sfx.exe ^
    -dir project ^
    -compress 9 ^
    -id co.touchify.testsfx ^
    -version %VERSION% ^
    -args "--sfx"

echo #
echo # Running SFX...

test\sfx.exe --test

