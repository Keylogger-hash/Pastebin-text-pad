package db

import (
	"log"

	"github.com/boltdb/bolt"
)

func ConvertByte(s string) []byte {
	// b,err := hex.DecodeString(s)
	// if err != nil {
	// 	fmt.Println("Can't convert")
	// 	return nil
	// }
	return []byte(s)
}
func ConvertString(b []byte) string {

	return string(b)
}

type DatabaseManager struct {
	db *bolt.DB
}

func AsyncGetDB(db *bolt.DB, bucket string, key []byte) chan []byte {
	ch := make(chan []byte)
	go func() {
		db.View(func(tx *bolt.Tx) error {
			bucket := ConvertByte(bucket)
			b := tx.Bucket(bucket)
			ans := b.Get(key)
			ch <- ans
			return nil
		})
	}()
	return ch
}
func AsyncUpdateDB(db *bolt.DB, bucket string, key []byte, value []byte) chan error {
	ch := make(chan error)
	go func() {
		err := db.Update(func(tx *bolt.Tx) error {
			tx.CreateBucketIfNotExists([]byte(bucket))
			bucket := ConvertByte(bucket)
			b := tx.Bucket(bucket)

			err := b.Put(key, value)
			return err
		})
		ch <- err
	}()
	return ch
}
func InitDB(connectDatabase string) *bolt.DB {
	db, err := bolt.Open(connectDatabase, 0777, nil)
	if err != nil {
		log.Fatal("Can't open database")
	}
	return db
}
