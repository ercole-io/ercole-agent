# ercole-agent
[![Build Status](https://travis-ci.org/amreo/ercole-agent.svg?branch=master)](https://travis-ci.org/amreo/ercole-agent)

This is the agent component for the Ercole project. Documentation available [here](https://ercole.netlify.com).

The agent is supposed to run on the same server of the Oracle instance you want to monitor.

Supported environments:

- Red Hat Enterprise Linux 5,6,7
- Windows Server 2012 and greater

## Requirements

- Go version 1.11 or greater.
- Go version 1.3.3 for the rhel5 branch

## How to build on Linux for Linux target

    make

## How to build on Linux for Windows target

    make windows

## How to run

Adjust the config.json configuration file with your server address, username
and password, then launch the binary:

    ./ercole-agent (Linux)
    ercole-agent.exe (Windows)
