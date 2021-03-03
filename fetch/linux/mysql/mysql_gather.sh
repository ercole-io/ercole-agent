#!/bin/bash
while [ "$1" != "" ]; do
    case $1 in
        -h | --host )           shift
                                host=$1
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

if [ -z "$host" ]; then
        exit 1
fi
if [ -z "$user" ]; then
        exit 1
fi
if [ -z "$password" ]; then
        exit 1
fi
if [ -z "$action" ]; then
        exit 1
fi

# set defaults
if [ -z "$database" ]; then
        database="mysql"
fi
if [ -z "$path" ]; then
        path="sql/mysql/"
fi

#actions
if [ $action ==  "instance" ]; then
	query="${path}instance.sql"
fi
if [ $action ==  "databases" ]; then
	query="${path}databases.sql"
fi
if [ $action ==  "table_schemas" ]; then
	query="${path}table_schemas.sql"
fi
if [ $action ==  "segment_advisors" ]; then
	query="${path}segment_advisors.sql"
fi

#echo "query=$query"
if [ -z "$loginPath" ]; then
        mysql -h $host --user=$user --password=$password --skip-column-names --database=$database -B <$query| sed "s/'/\'/;s/\t/\";\"/g;s/^/\"/;s/$/\"/;s/\n//g"
else
        mysql --login-path=$loginPath --skip-column-names --database=$database -B <$query| sed "s/'/\'/;s/\t/\";\"/g;s/^/\"/;s/$/\"/;s/\n//g"
fi
