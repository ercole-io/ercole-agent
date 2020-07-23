#!/bin/sh
ps -eocommand | grep -v grep | grep ora_pmon_ | cut -d'_' -f 3-