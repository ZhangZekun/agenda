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

func CreateAMeeting(meeting *Meeting) {
	currentName := User.GetCurUserName()
	if currentName == "" {
		fmt.Println("You haven't logged in")
		return
	}
	var allMeetings map[string]Meeting = GetAllMeetingInfo()
	fmt.Println(allMeetings)
	meeting.Id = strconv.Itoa(len(allMeetings)) //initial id is 0
	meeting.Sponsor = currentName
	meeting.Participants = append(meeting.Participants, currentName)
	allMeetings[meeting.Id] = *meeting

	/*
		allUser := User.GetAllUserInfo()

		for key, val := range meeting.Participants {

		}
	*/
	fmt.Println(allMeetings)
	fout, _ := os.Create("data/Meeting.json")
	defer fout.Close()
	b, _ := json.Marshal(allMeetings)
	fmt.Println(b)
	fout.Write(b)
}

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

func check(r error) {
	if r != nil {
		log.Fatal(r)
	}
}
