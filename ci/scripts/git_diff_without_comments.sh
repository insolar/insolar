#!/bin/bash

# Ignore all one-line comments.
git diff -w -G'(^[^\*# /])|(^#\w)|(^\s+[^\*#/])' --exit-code
