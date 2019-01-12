fen=$3
depth=$1
break=$2
colordiff <(go run *.go "${fen}" "${depth}" "${break}" | sed 's/ *//g' | sort) <(~/workspace/a.out ${depth} -${break} "${fen}" | sed 's/ *//g' | sort)
rm log.txt
