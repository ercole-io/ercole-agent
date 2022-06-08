#!/bin/sh

# Copyright (c) 2022 Sorint.lab S.p.A.
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.


PORT=$1
USER=$2
PASSWORD=$3
FILE=$4
DBNAME=$5
SCHEMA_NAME=$6


if [ -z "$PORT" ]; then
  >&2 echo "Missing PORT parameter"
  exit 1
fi
if [ -z "$USER" ]; then
  >&2 echo "Missing USER parameter"
  exit 1
fi
if [ -z "$PASSWORD" ]; then
  >&2 echo "Missing PASSWORD parameter"
  exit 1
fi
if [ -z "$FILE" ]; then
  >&2 echo "Missing FILE parameter"
  exit 1
fi
if [ -z "$DBNAME" ]; then
  >&2 echo "Missing DBNAME parameter"
  exit 1
fi
if [ -z "$SCHEMA_NAME" ]; then
  >&2 echo "Missing SCHEMA_NAME parameter"
  exit 1
fi

export PGPASSWORD=$PASSWORD
psql -t -A -d $DBNAME -p $PORT -U $USER -f $FILE -v schema_name=$SCHEMA_NAME
