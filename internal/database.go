package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
	"strconv"

	"github.com/go-redis/redis"
)

const redisFriendListName = "friends"
const dateFormat = "2006-01-02"

// TODO Remove
const data = `
{
	"name": "Anton",
	"dateContacted": "2019-06-19",
	"contactFrequency": 1
}
`

// Because Golang doesn't work nicely with easy time formats for some reason
type MyTime time.Time

func (m *MyTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse(dateFormat, s)

	if err != nil {
		return err
	}

	*m = MyTime(t)

	return nil
}

func (m MyTime) MarshalJSON() ([]byte, error) {
	time := m.Format(dateFormat)
	time = strconv.Quote(time)
	return []byte(time), nil
}

func (m MyTime) Before(t time.Time) bool {
	mt := time.Time(m)
	return mt.Before(t)
}

func (m MyTime) AddDays(days int) time.Time {
	t := time.Time(m)
	return t.AddDate(0, 0, days)
}

func (m MyTime) Format(format string) string {
	t := time.Time(m)
	return t.Format(format)
}

type Friend struct {
	id               int64
	Name             string `json:"name"`
	DateContacted    MyTime `json:"dateContacted"`
	ContactFrequency int    `json:"contactFrequency"`
}

func (f Friend) GetData() []byte {
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
	AddFriend(friend Friend)
	GetAllFriends() []Friend
	SetFriendAsContacted(friend Friend)
}

type RedisDatabase struct {
	client redis.Client
}

func GetDatabase() Database {
	// This will return the correct database based on settings
	// At the moment only Redis is allowed, so we'll just return a Redis database
	return NewRedisDatabase()
}

func ExampleRedis() {
	redis := NewRedisDatabase()
	friends := redis.GetAllFriends()
	fmt.Println(friends)
}

func NewRedisDatabase() RedisDatabase {
	url := os.Getenv("REDIS_URL")
	client := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: "",
		DB:       0,
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
	var friend Friend
	for i, value := range result {
		friend = NewFriendFromString(value)
		friend.id = int64(i)
		friends = append(friends, friend)
	}

	return friends
}

func (r RedisDatabase) SetFriendAsContacted(friend Friend) {
	friend.DateContacted = MyTime(time.Now())
	mess, err := json.Marshal(friend)
	if err != nil {
		panic(err)
	}
	fmt.Println("Putting", friend.Name, "as contacted:", string(mess))
	r.client.LSet(redisFriendListName, friend.id, mess)
}
