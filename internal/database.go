package internal

import (
	"os"
	"github.com/go-redis/redis"
	"time"
	"fmt"
	"encoding/json"
)

const redisFriendListName = "friends"
const dateFormat = "200T6-01-02"

// TODO Remove
const data = `
{
	"name": Anton",
	"dateContacted": "2019-06-19",
	"contactFrequency": "1"
}
`


// Because Golang doesn't work nicely with easy time formats for some reason
type MyTime time.Time

func (m *MyTime) UnmarshalJSON(p []byte) error {
	t, err := time.Parse(dateFormat, string(p))

	if err != nil {
		return err
	}

	*m = MyTime(t)

	return nil
}

type Friend struct {
	Name 				string 		`json:"name"`
	DateContacted 		MyTime	 	`json:"dateContacted"`
	ContactFrequency 	int 		`json:"contactFrequency"`
}

func (f Friend) GetData() []byte  {
	value, err := json.Marshal(f)
	if err != nil {
		panic(err)
	}
	return value
}

func NewFriendFromString(value string) Friend {
	friend := Friend{}
	err := json.Unmarshal([]byte(value), &friend)
	if err != nil {
		panic(err)
	}
	return friend
}

type Database interface {
	AddFriend()
	GetFriends()
	SetFriendAsContacted()
}

type RedisDatabase struct {
	client redis.Client
}

func ExampleRedis() {
	redis := NewRedisDatabase()
	friend := NewFriendFromString(data)
	redis.AddFriend(friend)
	friends := redis.GetAllFriends()
	fmt.Println(friends)
}

func NewRedisDatabase() RedisDatabase {
	url := os.Getenv("REDIS_URL")
	client := redis.NewClient(&redis.Options{
		Addr:			url,
		Password:		"",
		DB:				0,
	})
	db := RedisDatabase{
		client: *client,
	}
	return db
}

func (r RedisDatabase) AddFriend(friend Friend) {
	r.client.RPush(redisFriendListName, friend.GetData())
}

func (r RedisDatabase) GetAllFriends() []Friend {
	result, err := r.client.LRange(redisFriendListName, 0, -1).Result()
	if err != nil {
		panic(err)
	}

	var friends []Friend
	for _, value := range result {
		friends = append(friends, NewFriendFromString(value))
	}

	return friends
	
}
