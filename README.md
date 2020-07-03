# ercole-agent
[![Build Status](https://travis-ci.org/ercole-io/ercole-agent.svg?branch=master)](https://travis-ci.org/ercole-io/ercole-agent)
[![codecov](https://codecov.io/gh/ercole-io/ercole-agent/branch/master/graph/badge.svg)](https://codecov.io/gh/ercole-io/ercole-agent)
[![Go Report Card](https://goreportcard.com/badge/github.com/ercole-io/ercole-agent)](https://goreportcard.com/report/github.com/ercole-io/ercole-agent)

This is the agent component for the Ercole project.
Documentation available [here](https://ercole.io).

Supported environments:

- Red Hat Enterprise Linux 5, 6, 7, 8
- Windows Server 2012 and greater

## Requirements
- Go version 1.13 or greater.

## How to build on Linux for Linux target
    make

## How to build on Linux for Windows target
    make windows

## How to run

Adjust the config.json configuration file with your server address, username
and password, then launch the binary:

    ./ercole-agent (Linux)
    ercole-agent.exe (Windows)
