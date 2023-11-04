package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"reflect"
	"time"
)

func main() {
	db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// create bucket and put
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("MyBucket"))
		if err != nil {
			return fmt.Errorf("create bucket err %v", err)
		}

		if err = b.Put([]byte("k1"), []byte("v1")); err != nil {
			return fmt.Errorf("put err %v", err)
		}
		if err = b.Put([]byte("k2"), []byte("v2")); err != nil {
			return fmt.Errorf("put err %v", err)
		}
		return nil
	})

	// get
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))
		v := b.Get([]byte("k1"))
		fmt.Println(reflect.DeepEqual(v, nil)) // false
		fmt.Println(v == nil)                  // false
		fmt.Println(string(v))

		v = b.Get([]byte("k3"))
		fmt.Println(reflect.DeepEqual(v, nil)) // false
		fmt.Println(v == nil)                  // true
		fmt.Println(string(v))
		return nil
	})

}
