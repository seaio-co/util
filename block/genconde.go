package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"io"
	"log"
	"time"
	"unsafe"
)

var (
	_ = unsafe.Sizeof(0)
	_ = io.ReadFull
	_ = time.Now()
)

type Person struct {
	Name string
	Age uint8
}

func main() {

	db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("person"))
		if err != nil {
			return err
		}
		return nil
	})

	person := Person{
		Name: "testname",
		Age:  20,
	}

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("person"))

		buf, _ := person.Marshal(nil)
		fmt.Println(buf, string(buf))
		b.Put([]byte(person.Name), buf)

		fmt.Printf("Gencode encoded size: %v\n", len(buf))

		return nil
	})

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("person"))

		if err := b.ForEach(func(k, v []byte) error {
			fmt.Printf("%s is %s.\n", k, v)
			p := Person{}
			p.Unmarshal(v)
			fmt.Println(p)
			return nil
		}); err != nil {
			return err
		}

		return nil
	})

}

