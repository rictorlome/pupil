fen=$1
depth=$2
break=$3
colordiff <(go run *.go "${fen}" "${depth}" "${break}" | sed 's/ *//g' | sort) <(~/workspace/a.out ${depth} -${break} "${fen}" | sed 's/ *//g' | sort)
