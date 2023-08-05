### Results

## volitile-lru

# Steps to reproduce

* Create 30 keys with expiration
* Get keys 1-10 one time
* Get keys 21-30 one time
* Write keys until 10 keys are evicted
* observe that evicted keys are exactly keys 11, 13-20, and 3. The last one is weird but perhaps this is what "approximated LRU" means

## allkeys-lru

# Steps to reproduce

* Create 30 keys
* Set expiration for first 20 keys
* Get keys 11-20 three times
* Get keys 21-30 one time
* Write keys until 10 keys are evicted
* observe that evicted keys are among keys 1-100+ (meaning even recently inserted are deleted, since they are least recently used)


## volitile-lfu

# Steps to reproduce

* Set lfu-log-factor 0
* Create 30 keys with expiration time
* Get keys 11-20 101 times
* Get keys 21-30 101 time
* Get keys 1-10 1 time
* Write keys until 10 keys are evicted
* observe that 8 keys are evicted in range 1-10 and 2 keys are from range 21-30


## allkeys-lfu

# Steps to reproduce

* Set lfu-log-factor 0
* Create 30 keys with expiration time
* Get keys 1-10 1 times
* Get keys 11-20 1 time
* Get keys 21-30 1 time
* Write keys until 10 keys are evicted
* observe all keys are distributed evenly (kinda) across all keys


## volitile-random

# Steps to reproduce
* Create 30 keys with expiration time
* Get keys 11-20 101 times
* Get keys 21-30 101 time
* Get keys 1-10 1 time
* Write keys until 10 keys are evicted
* observe that deleted keys are evenly distributed across keys 1-30 (set to expire)

## volitile-random

# Steps to reproduce
* Create 30 keys with expiration time
* Get keys 11-20 101 times
* Get keys 21-30 101 time
* Get keys 1-10 1 time
* Write keys until 10 keys are evicted
* observe that deleted keys are evenly distributed across all keys

## volitile-ttl

# Steps to reproduce
* Create 10 keys with expiration time 100000, 10 keys with expiration time 500 and 10 keys with expiration time 200000 
* Get keys 1-10 1 times
* Get keys 11-20 1 time
* Get keys 21-30 1 time
* Write keys until 10 keys are evicted
* observe that deleted keys are keys 11-18, 20, and 2

## noeviction

# Steps to reproduce
* Write keys until out of memory
* Observe memory error
