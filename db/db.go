package db

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

func ConvertByte(s string) []byte {
	return []byte(s)
}
func ConvertString(b []byte) string{
	return string(b)
}

type DatabaseManager struct {
	db *bolt.DB
}
func UpdateDB(db *bolt.DB, bucket string, key string, value string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("Paste"))
		bucket := ConvertByte(bucket)
		b := tx.Bucket(bucket)
		k := ConvertByte(key)
		v := ConvertByte(value)
		err := b.Put(k, v)
		return err
	})
	return err
}
func GetDB(db *bolt.DB, bucket string, key string) ([]byte,error) {
	var ans []byte
	err := db.View(func(tx *bolt.Tx) error {
		bucket := ConvertByte(bucket)
		b := tx.Bucket(bucket)
		key := ConvertByte(key)
		ans = b.Get(key)
		if ans == nil {
			return fmt.Errorf("Key is not exist")
		}
		return nil
	})
	return ans,err
}
func InitDB(connectDatabase string) *bolt.DB{
	db, err := bolt.Open(connectDatabase, 0777, nil)
	if err != nil {
		log.Fatal("Can't open database")
	}
	return db
}
// func main() {
// 	db := InitDB("bolt.db")
// 	// for i := 0; i < 10; i++ {
// 	// 	UpdateDB(db, "Paste", fmt.Sprintf("%v", i), fmt.Sprintf("Hello world %v!", i))
// 	// }
// 	for i := 0; i < 10; i++ {
// 		ans, err := GetDB(db, "Paste", fmt.Sprintf("%v",i))
// 		if err != nil {
// 			fmt.Println("Skip")
// 		} else {
// 			fmt.Println(convertString(ans))
// 		}

// 	}
// }
