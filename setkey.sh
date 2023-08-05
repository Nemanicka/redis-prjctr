
for i in $(seq $1 $2);
do
    docker exec -ti redis-master redis-cli set "thisIsAKey#$i" "ThisIsAwesomeValue#$i" # total ~30B
done
