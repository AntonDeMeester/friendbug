package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/joho/godotenv"

	// Own packages
	"friendbug/internal"
)

//TODO TO BE TRANSFERED TO SETTINGS
const maxNumberOfContactsPerDay = 3
const dateFormat = "2006-01-02"
const targetNumber = "+32496952214"

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	client := internal.GetDatabase()
	friends := client.GetAllFriends()

	toContactFriends := []internal.Friend{}

	for _, friend := range friends {
		if friend.DateContacted.Before(time.Now().AddDate(0, 0, -1)) {
			toContactFriends = append(toContactFriends, friend)
		}
	}

	contactedFriends := sendReminderForFriends(toContactFriends)
	sendMessage(contactedFriends)
	for _, friend := range contactedFriends {
		client.SetFriendAsContacted(friend)
	}
}

func sendReminderForFriends(friends []internal.Friend) (toContactFriends []internal.Friend) {
	toSendRemindersFor := []internal.Friend{}
	// If we have more friends than the max number to cintact, we need to choose from them
	if len(friends) > maxNumberOfContactsPerDay {
		selectedFriend := internal.Friend{}
		for i := 0; i < maxNumberOfContactsPerDay; i++ {
			selectedFriend, friends = selectRandomFriendWeighted(friends)
			toSendRemindersFor = append(toSendRemindersFor, selectedFriend)
		}
	} else {
		// If we have less than the max amount, we just add them all
		for _, friend := range(friends) {
			toSendRemindersFor = append(toSendRemindersFor, friend)
		}
	}

	return toSendRemindersFor

}

func selectRandomFriendWeighted(friends []internal.Friend) (selectedFriend internal.Friend, remainingFriends []internal.Friend) {
	if len(friends) == 0 {
		return internal.Friend{}, friends
	}

	totalWeight := 0.0
	for _, friend := range friends {
		totalWeight += 360 / float64(friend.ContactFrequency)
	}

	selected := rand.Float64() * totalWeight
	for i, friend := range friends {
		selected -= 360 / float64(friend.ContactFrequency)
		if selected < 0 {
			// Save the friend
			selectedFriend = friend
			// Copy the last element to the place of the previous friend
			friends[i] = friends[len(friends)-1]
			// Remove the last element of the friends list
			friends = friends[:len(friends)-1]

			return selectedFriend, friends
		}
	}
	// If we somehow don't find one, we return the last one
	return friends[len(friends)], friends[:len(friends)-1]
}

func sendMessage(friends []internal.Friend) {
	if len(friends) == 0 {
		return
	}

	length := len(friends)
	friendsString := ""
	for i, friend := range friends {
		friendsString += friend.Name
		if i < length-2 {
			friendsString += ", "
		} else if i == length-2 {
			friendsString += " and "
		}
	}

	sendMessage := "Hey Anton, you should probably contact " + friendsString + " today!"
	fmt.Println(sendMessage)
	//internal.SendMessageTwilio(sendMessage, targetNumber)
}
