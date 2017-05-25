#!/bin/sh

mkdir publish

echo "#"
echo "# Publishing x64 artifacts..."

cd bin/x64
zip sfx-x64.zip *
mv sfx-x64.zip ../../publish/

echo "#"
echo "# Publishing i386 artifacts..."

cd ../i386
zip sfx-i386.zip *
mv sfx-i386.zip ../../publish/

echo "#"
echo "# Publishing arm artifacts..."

cd ../arm
zip sfx-arm.zip *
mv sfx-arm.zip ../../publish/

cd ../..