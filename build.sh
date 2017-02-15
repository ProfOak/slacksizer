#!/bin/bash

# cross compiling in go is nice (:

echo building linux
env GOOS=linux GOARCH=386 go build -o slacksizer
tar cvf linux_x86.tar slacksizer
rm slacksizer

echo building mac osx
env GOOS=darwin GOARCH=386 go build -o slacksizer
tar cvf osx_x86.tar slacksizer
rm slacksizer

echo building windows
env GOOS=windows GOARCH=386 go build -o slacksizer.exe
zip windows_x86.zip slacksizer.exe
rm slacksizer.exe

