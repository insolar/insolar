#!/usr/bin/env bash

make vendor
git diff -w -G'(^[^\*# /])|(^#\w)|(^\s+[^\*#/])' --exit-code vendor/
