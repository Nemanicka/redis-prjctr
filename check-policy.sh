docker exec -ti redis-master redis-cli flushall
docker exec -ti redis-master redis-cli config resetstat
docker exec -ti redis-master redis-cli config set maxmemory-policy $1
#docker exec -ti redis-master redis-cli config set lfu-log-factor 0
./setkey.sh 1 30
./expirekey.sh 1 10 100000 
./expirekey.sh 11 20 500 
./expirekey.sh 21 30 200000 

sleep 1

#for i in $(seq 1 101);
#do
    ./getkey.sh 11 20
#done

#for i in $(seq 1 101);
#do
    ./getkey.sh 21 30
#done

./getkey.sh 1 10


c=100

echo $ev

while true
do
    ((c=c+1))
    ./setkey.sh $c $c
    ev=`./check_evict_num.sh`
    echo $ev
    if [[ $ev == '10' ]]; then
        echo "Evicted 10 keys"
        break
    fi
done

./getkey.sh 1 30

