package Meeting

import (
	"agenda/entity/User"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

type Meeting struct {
	Title        string
	Sponsor      string
	Participants []string
	StartTime    time.Time
	EndTime      time.Time
	Id           string
}

//时间处理函数
func DateToString(date time.Time) string {
	return date.Format("2006-01-02/15:04")
}
func StringToDate(date string) (time.Time, error) {
	the_time, err := time.Parse("2006-01-02/15:04", date)
	return the_time, err
}
func SmallDate(date1, date2 time.Time) bool {
	return date1.Before(date2) || date1.Equal(date2)
}
func LargeDate(date1, date2 time.Time) bool {
	return date1.After(date2) || date1.Equal(date2)
}
func TimeContact(newDateS, newDateE, oriDateS, oriDateE time.Time, userName string, meetingId string) bool {
	//newstart time is after newend time, wrong!
	if LargeDate(newDateS, newDateE) {
		fmt.Println("start time can't be greater than end time")
		return true
	}
	if SmallDate(newDateE, oriDateS) || LargeDate(newDateS, oriDateE) {
		return false
	}
	fmt.Println("time contract with " + userName + "' meeting with id:" + meetingId)
	return true
}

func CreateAMeeting(meeting *Meeting) {
	currentName := User.GetCurUserName()
	if currentName == "" {
		fmt.Println("You haven't logged in")
		return
	}
	var allMeetings map[string]Meeting = GetAllMeetingInfo()
	meeting.Id = strconv.Itoa(len(allMeetings)) //initial id is 0
	meeting.Sponsor = currentName
	meeting.Participants = append(meeting.Participants, currentName)
	allMeetings[meeting.Id] = *meeting

	var map[string]*User.User = GetAllUserInfo()
	//check all participanter exist, and time contract
	for _, pName := range meeting.Participants {
		//check if the user exist
		if _, ok := allUser[pName]; !ok {
			fmt.Println("No such user:" + pName + "!")
			return
		}
		//check if the user's old meeting is contract with the new one
		for _, meetingId := range allUser[pName].ParticipantMeeting {
			if TimeContact(meeting.StartTime, meeting.EndTime, allMeetings[meetingId].StartTime, allMeetings[meetingId].EndTime, pName, meetingId) {
				return
			}
		}
		allUser[pName].ParticipantMeeting = append(allUser[pName].ParticipantMeeting, meeting.Id)
	}
	allUser[currentName].SponsorMeeting = append(allUser[currentName].SponsorMeeting, meeting.Id)

	fout, _ := os.Create("data/Meeting.json")
	defer fout.Close()
	b, _ := json.Marshal(allMeetings)
	//	fmt.Println(b)
	fout.Write(b)
	foutuser, _ := os.Create("data/User.json")
	defer foutuser.Close()
	buser, _ := json.Marshal(allUser)
	foutuser.Write(buser)
}

//func GetAllMeetingIDOfOneUser(name string,  map[string]Meeting)

//load all meeting infomation
func GetAllMeetingInfo() map[string]Meeting {
	byteIn, err := ioutil.ReadFile("data/Meeting.json")
	check(err)
	var allMeetingInfo map[string]Meeting
	json.Unmarshal(byteIn, &allMeetingInfo)
	if allMeetingInfo == nil {
		allMeetingInfo = make(map[string]Meeting)
	}
	return allMeetingInfo
}

//load all user infomation to User.AllUserInfo
func GetAllUserInfo() map[string]*User.User {
	byteIn, err := ioutil.ReadFile("data/User.json")
	check(err)
	var allUserInfo map[string]*User.User
	json.Unmarshal(byteIn, &allUserInfo)
	return allUserInfo
}

func check(r error) {
	if r != nil {
		log.Fatal(r)
	}
}
