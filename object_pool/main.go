package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync"
)

var ObjectPoolThresholdExited = fmt.Errorf("ObjectPoolThresholdExited")

type ObjectPoolInstanceBuilder = func() interface{}

type ObjectPoolInstance interface{}

type ObjectPool struct {
	sync.Mutex
	available       []ObjectPoolInstance
	inUse           []ObjectPoolInstance
	threshold       int
	instanceBuilder ObjectPoolInstanceBuilder
}

func NewObjectPool(threshold int, instanceBuilder func() interface{}) ObjectPool {
	return ObjectPool{
		available:       []ObjectPoolInstance{},
		inUse:           []ObjectPoolInstance{},
		threshold:       threshold,
		instanceBuilder: instanceBuilder,
	}
}

func (op *ObjectPool) acquire() (ObjectPoolInstance, error) {
	if len(op.available) == 0 {
		if len(op.inUse) >= op.threshold {
			return nil, ObjectPoolThresholdExited
		}
		instance := op.instanceBuilder()
		fmt.Println(instance)
		op.available = append(op.available, instance)
	}

	op.Lock()
	defer op.Unlock()
	instance := op.available[len(op.available)-1]
	op.inUse = append(op.inUse, instance)
	op.available = op.available[0 : len(op.available)-1]

	return instance, nil
}

func (op *ObjectPool) free(instance ObjectPoolInstance) {
	for i, _instance := range op.inUse {
		if instance == _instance {
			op.Lock()
			op.inUse = append(op.inUse[0:i], op.inUse[i+1:]...)
			op.available = append(op.available, instance)
			op.Unlock()

			break
		}
	}
	fmt.Println("Available instances:", len(op.available), op.available)
	fmt.Println("In use instances:", len(op.inUse), op.inUse)
}

type PoolObjectImpl struct {
	ObjectPoolInstance
	Name string
}

func BuilderImpl() interface{} {
	obj := PoolObjectImpl{
		Name: "Pool",
	}
	return &obj
}

func BuilderRedis() interface{} {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return rdb
}

func main() {
	var ctx = context.Background()
	pool := NewObjectPool(4, BuilderRedis)

	var instances []ObjectPoolInstance
	for i := 0; i < 6; i++ {
		instance, err := pool.acquire()
		if instance != nil {
			instance.(*redis.Client).Set(ctx, "testKey", "test value", 1000)
		}
		instances = append(instances, instance)
		fmt.Println("Acquire instance: ", instance, err)
	}

	for _, instance := range instances {
		fmt.Println("Free", instance)
		if instance != nil {
			fmt.Println(instance.(*redis.Client).Get(ctx, "testKey"))
		}
		pool.free(instance)
	}
}
