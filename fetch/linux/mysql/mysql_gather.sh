#!/bin/bash
while [ "$1" != "" ]; do
    case $1 in
        -h | --host )           shift
                                host=$1
                                ;;
        -P | --port )           shift
                                port=$1
                                ;;
        -u | --user )           shift
                                user=$1
                                ;;
        -p | --password )       shift
                                password=$1
                                ;;
        -d | --database )       shift
                                database=$1
                                ;;
        --path )                shift
                                path=$1
                                ;;
        -a | --action )         shift
                                action=$1
                                ;;
        --loginPath )           shift
                                loginPath=$1
                                ;;
        * )
                                exit 1
    esac
    shift
done

if [ -z "$host" ] ||
   [ -z "$user" ] ||
   [ -z "$password" ] ||
   [ -z "$action" ]; then
        exit 1
fi

# set defaults
if [ -z "$database" ]; then
        database="mysql"
fi
if [ -z "$path" ]; then
        path="sql/mysql/"
fi
if [ -z "$port" ]; then
        port="3306"
fi

#actions
if [ "$action" ==  "instance" ] || 
   [ "$action" ==  "old_instance" ] ||
   [ "$action" ==  "databases" ] || 
   [ "$action" ==  "table_schemas" ] || 
   [ "$action" ==  "segment_advisors" ] ||
   [ "$action" ==  "slave_hosts" ] ||
   [ "$action" ==  "slave_status" ] ||
   [ "$action" == "high_availability" ] ||
   [ "$action" == "version" ]; then
	query="${path}${action}.sql"
else 
        exit 1
fi

if [ -z "$loginPath" ]; then
        mysql -h "$host" --port="$port" --user="$user" --password="$password" --skip-column-names --database="$database" -B <"$query"| sed "s/'/\'/;s/\t/\";\"/g;s/^/\"/;s/$/\"/;s/\n//g"
else
        mysql --login-path="$loginPath" --skip-column-names --database="$database" -B <"$query"| sed "s/'/\'/;s/\t/\";\"/g;s/^/\"/;s/$/\"/;s/\n//g"
fi
