package main

import (
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

func main() {

	pool := newPool()
	conn := pool.Get()
	defer conn.Close()

	err := ping(conn)
	if err != nil {
		fmt.Println(err)
	}

	err = set(conn)
	if err != nil {
		fmt.Println(err)
	}

	err = get(conn)
	if err != nil {
		fmt.Println(err)
	}

	err = setStruct(conn)
	if err != nil {
		fmt.Println(err)
	}

	err = getStruct(conn)
	if err != nil {
		fmt.Println(err)
	}
}

func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

// ping tests connectivity for redis (PONG should be returned)
func ping(c redis.Conn) error {
	pong, err := c.Do("PING")
	if err != nil {
		return err
	}
	fmt.Println(pong, err)
	// Output: PONG <nil>

	return nil
}

// set executes the redis Set command
func set(c redis.Conn) error {
	_, err := c.Do("SET", "key", "value")
	if err != nil {
		return err
	}
	return nil
}

func get(c redis.Conn) error {
	reply, err := c.Do("GET", "key")
	val, err := redis.String(reply, err)
	if err != nil {
		return (err)
	}
	fmt.Println("key", val)

	val2, err := redis.String(c.Do("GET", "key2"))
	if err == redis.ErrNil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist

	return nil
}

// User is a simple user struct for this example
type User struct {
	Username  string `json:"username"`
	MobileID  int    `json:"mobile_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func setStruct(c redis.Conn) error {

	usr := User{
		Username:  "otto",
		MobileID:  1234567890,
		Email:     "ottoM@repoman.com",
		FirstName: "Otto",
		LastName:  "Maddox",
	}

	json, err := json.Marshal(usr)
	if err != nil {
		return err
	}

	_, err = c.Do("SET", "user:"+usr.Username, json)
	if err != nil {
		return err
	}

	return nil
}

func getStruct(c redis.Conn) error {

	s, err := redis.String(c.Do("GET", "user:otto"))
	if err == redis.ErrNil {
		fmt.Println("User does not exist")
	} else if err != nil {
		return err
	}

	usr := User{}
	err = json.Unmarshal([]byte(s), &usr)

	fmt.Printf("%+v\n", usr)

	return nil

}
