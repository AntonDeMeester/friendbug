package friendbug

import (
	"fmt"
	"math/rand"
	"time"
	"os"

	// Own packages
	"friendbug/internal"
)

//TODO TO BE TRANSFERED TO SETTINGS
const maxNumberOfContactsPerDay = 2
const dateFormat = "2006-01-02"

func ContactFriends() {
	client := internal.GetDatabase()
	friends := client.GetAllFriends()

	toContactFriends := []internal.Friend{}

	for _, friend := range friends {
		// Take 23 hour to not interfere with exact time differences.
		if friend.DateContacted.Before(time.Now().Add(time.Hour * -23).Add(time.Duration(-24 * friend.ContactFrequency) * time.Hour)) {
			toContactFriends = append(toContactFriends, friend)
		} 
	}

	contactedFriends := sendReminderForFriends(toContactFriends)
	sendMessage(contactedFriends)
	for _, friend := range contactedFriends {
		client.SetFriendAsContacted(friend)
		fmt.Println(friend)
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
	targetNumber := os.Getenv("TARGET_NUMBER")
	internal.SendMessageTwilio(sendMessage, targetNumber)
}
