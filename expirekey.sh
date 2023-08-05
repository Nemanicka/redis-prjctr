
for i in $(seq $1 $2);
do
    docker exec -ti redis-master redis-cli EXPIRE "thisIsAKey#$i" $3
done
