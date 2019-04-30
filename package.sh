#!/bin/sh

makepkg --printsrcinfo > .SRCINFO

rm -rf pkg src nextver*

makepkg -si
