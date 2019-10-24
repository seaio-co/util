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

func (d *Person) Size() (s uint64) {

	{
		l := uint64(len(d.Name))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	s += 1
	return
}
func (d *Person) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	{
		l := uint64(len(d.Name))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+0] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+0] = byte(t)
			i++

		}
		copy(buf[i+0:], d.Name)
		i += l
	}
	{

		buf[i+0+0] = byte(d.Age >> 0)

	}
	return buf[:i+1], nil
}

func (d *Person) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+0] & 0x7F)
			for buf[i+0]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+0]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		d.Name = string(buf[i+0 : i+0+l])
		i += l
	}
	{

		d.Age = 0 | (uint8(buf[i+0+0]) << 0)

	}
	return i + 1, nil
}
