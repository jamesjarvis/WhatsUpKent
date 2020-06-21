package api

import (
	"hash/fnv"
	"log"
	"time"

	badger "github.com/dgraph-io/badger"
)

//GetCache attempts to retrieve a cached version of the request
func (config *Config) GetCache(query string) (*string, error) {
	var valCopy []byte
	// var err error
	err := config.CacheDB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(query))
		if err != nil {
			return err
		}

		valCopy, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	temp := string(valCopy)
	log.Println("Successful cache-hit!")

	return &temp, nil
}

//SetCache stores the query and answer in the database
func (config *Config) SetCache(query string, answer string) error {
	err := config.CacheDB.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(query), []byte(answer)).WithTTL(time.Hour)
		err := txn.SetEntry(e)
		return err
	})
	log.Println("Saved to cache!")
	return err
}

//PerformCachedQuery is the main accessor with cache abilities
func (config *Config) PerformCachedQuery(query string) (*string, error) {
	//Try to retrieve from cache
	answer, err := config.GetCache(query)
	if err != nil {
		//If its not in the cache
		if err == badger.ErrKeyNotFound {
			// Lock.Lock() //The lock can be used if deemed necessary. Doesn't seem to be an issue for now
			res, queryErr := config.PerformQuery(query)
			// Lock.Unlock()
			if queryErr != nil {
				return nil, queryErr
			}

			//Save to cache in a goroutine
			go config.SetCache(query, *res)

			//Return result
			return res, queryErr
		}
		return nil, err
	}
	//If it is in the cache
	return answer, nil
}

//PerformQuery is the main db accessor, without caching abilities.
func (config *Config) PerformQuery(query string) (*string, error) {
	//Get client connection
	result, err := config.DBClient.ReadOnly(query)
	if err != nil {
		return nil, err
	}
	res := string(result)

	//Return result
	return &res, err
}

//Hash generates a uint32 hash of a string
func Hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
