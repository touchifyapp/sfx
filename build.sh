#!/bin/sh

echo "#"
echo "# Building Bundler (linux x64)..."

GOOS=linux GOARCH=amd64 go build \
    -o bin/x64/bundler \
    ./bundler

echo "#"
echo "# Building Bundler (linux i386)..."

GOOS=linux GOARCH=386 go build \
    -o bin/i386/bundler \
    ./bundler

echo "#"
echo "# Building Bundler (linux arm)..."

GOOS=linux GOARCH=arm go build \
    -o bin/arm/bundler \
    ./bundler

echo "#"
echo "# Building Bundler (windows x64)..."

GOOS=windows GOARCH=amd64 go build \
    -o bin/x64/bundler.exe \
    ./bundler

echo "#"
echo "# Building Bundler (windows i386)..."

GOOS=windows GOARCH=386 go build \
    -o bin/i386/bundler.exe \
    ./bundler

echo "#"
echo "# Building SFX (windows x64)..."

GOOS=windows GOARCH=amd64 go generate ./base
GOOS=windows GOARCH=amd64 go build \
    -o bin/x64/sfx.exe \
    -ldflags "-H windowsgui" \
    ./base

echo "#"
echo "# Building SFX (windows i386)..."

GOOS=windows GOARCH=386 go generate ./base
GOOS=windows GOARCH=386 go build \
    -o bin/i386/sfx.exe \
    -ldflags "-H windowsgui" \
    ./base

echo "#"
echo "# Building SFX verbose (windows x64)..."

GOOS=windows GOARCH=amd64 go generate ./base
GOOS=windows GOARCH=amd64 go build \
    -o bin/x64/sfxv.exe \
    -tags verbose \
    ./base
    
echo "#"
echo "# Building SFX verbose (windows i386)..."

GOOS=windows GOARCH=386 go generate ./base
GOOS=windows GOARCH=386 go build \
    -o bin/i386/sfxv.exe \
    -tags verbose \
    ./base
