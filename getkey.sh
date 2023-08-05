for i in $(seq $1 $2);
do
    docker exec -ti redis-master redis-cli get "thisIsAKey#$i" 
done
