evicted=`docker exec -ti redis-master redis-cli info stats | grep evicted_keys  |  grep -oE '[0-9]+'`
echo $evicted
#num=${evicted#*:}
#echo $num
