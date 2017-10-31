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

	var allUser map[string]*User.User = GetAllUserInfo()
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

func OperateParticipants(title, operation string, participants []string){
	current_user := User.GetCurUserName()
	if current_user == ""{
		fmt.Println("You haven't logged in")
		return 
	}

	check_participants(participants)

	var allUser map[string]*User.User = GetAllUserInfo()
	var all_meetings map[string]Meeting = GetAllMeetingInfo()
	count := 0

	for i := range all_meetings{
		if all_meetings[i].Title == title && all_meetings[i].Sponsor == current_user {

			//delete meeting
			if operation == "del" {
				slice := all_meetings[i].Participants
				slice = remove_element_in_slice(slice,participants)		
				var meet *Meeting = &Meeting{all_meetings[i].Title,all_meetings[i].Sponsor, slice, all_meetings[i].StartTime, all_meetings[i].EndTime, all_meetings[i].Id}
				all_meetings[i] = *meet

				if(len(all_meetings[i].Participants) == 0){
					for user := range allUser{
						if allUser[user].Username == all_meetings[i].Sponsor{
							var meeting_id []string
							meeting_id = append(meeting_id, meet.Id)
							allUser[user].SponsorMeeting = remove_element_in_slice(allUser[user].SponsorMeeting,meeting_id)
						}
					}
					delete(all_meetings, i)
				}
				//delete participants's meeting messages
				for user := range allUser{
					for name := range participants{
						if allUser[user].Username == participants[name]{
							var meeting_id []string
							meeting_id = append(meeting_id,  meet.Id)
							allUser[user].ParticipantMeeting = remove_element_in_slice(allUser[user].ParticipantMeeting,meeting_id)
							fmt.Println("delete succeed!\n")
							break
						}
					}
				}
				break
			}

			if operation == "add" {
				//check participants repeat or not
				for j := range all_meetings[i].Participants{
					for k := range participants{
						if all_meetings[i].Participants[j] == participants[k]{
							fmt.Println("The user " + participants[k]+"is already in the meeting\n")
							return 
						}
					}
				}

				slice := all_meetings[i].Participants
				slice = add_element_in_slice(slice,participants)
				var meet *Meeting = &Meeting{all_meetings[i].Title,all_meetings[i].Sponsor, slice, all_meetings[i].StartTime, all_meetings[i].EndTime, all_meetings[i].Id}
				all_meetings[i] = *meet

				//check participants time contract

				for _, pName := range participants{
					for _, meetingId := range allUser[pName].ParticipantMeeting {
						if TimeContact(meet.StartTime, meet.EndTime, all_meetings[meetingId].StartTime, all_meetings[meetingId].EndTime, pName, meetingId) {
							return
						}
					}
				}

				for user := range allUser{
					for name := range participants{
						if allUser[user].Username == participants[name]{
							allUser[user].ParticipantMeeting = append(allUser[user].ParticipantMeeting,all_meetings[i].Id)
							fmt.Println("add succeed!\n")
							break;
						}
					}
				}
			}
			count++
		}
		
	}

	if count == 0{
		fmt.Println("The meeting isn't exited\n")
		return 
	}

	fout, _ := os.Create("data/Meeting.json")
	defer fout.Close()
	bytes,_ := json.Marshal(all_meetings)
	fout.Write(bytes)

	fout_user, _ := os.Create("data/User.json")
	defer fout_user.Close()
	buser, _ := json.Marshal(allUser)
	fout_user.Write(buser)


}

func remove_element_in_slice(oldslice []string, elemsslice []string) []string{
	var slice []string
	flag := 1
	for i := range oldslice{
		flag = 1
		for j := range elemsslice{
			if oldslice[i] == elemsslice[j]{
				flag = 0
				break
			}
		}
		if flag == 1{
			slice = append(slice, oldslice[i])
		}
	}
	return slice
}

func add_element_in_slice(slice []string, elemsslice []string) []string{
	for i:=range elemsslice{
		slice = append(slice,elemsslice[i])
	}

	return slice
}

func check_participants(slice []string) {
	var allUser map[string]*User.User = GetAllUserInfo()
	for _, pname := range slice{
		if _, ok := allUser[pname];!ok{
			fmt.Println("No such user:"+pname+"!")
			return
		}
	}
}

func SearchMeeting(search_start, search_end string) {
	current_user := User.GetCurUserName()
	if current_user == ""{
		fmt.Println("You haven't logged in")
		return 
	}
	var all_meetings map[string]Meeting = GetAllMeetingInfo()
	start,_ := StringToDate(search_start)
	end,_ := StringToDate(search_end)

	if !SmallDate(start, end){
		fmt.Println("the StartTime is larger than EndTime\n")
		return
	}

	count := 0
	for meeting := range all_meetings{
		meet := all_meetings[meeting]
		if check_user_in_meeting(&meet, current_user){
			if !(SmallDate(meet.EndTime, start) || SmallDate(end, meet.StartTime)){
				fmt.Println(all_meetings[meeting].StartTime, all_meetings[meeting].EndTime, all_meetings[meeting].Title, all_meetings[meeting].Sponsor, all_meetings[meeting].Participants)
				count++
			}
		}
	}

	if count == 0{
		fmt.Println("You haven't create or participant in any meeting\n")
		return
	}

}

func check_user_in_meeting(meeting *Meeting, username string) bool {
	if meeting.Sponsor == username{
		return true
	}
	for participant := range meeting.Participants{
		if meeting.Participants[participant] == username{
			return true
		}
	}
	return false
}

func CancelMeeting(title string) {
	current_user := User.GetCurUserName()
	if current_user == ""{
		fmt.Println("You haven't logged in")
		return 
	}

	var allUser map[string]*User.User = GetAllUserInfo() 
	var all_meetings map[string]Meeting = GetAllMeetingInfo()
	count := 0

	for meeting := range all_meetings{
		if all_meetings[meeting].Sponsor == current_user && all_meetings[meeting].Title == title{
			for user := range allUser{
				if allUser[user].Username == all_meetings[meeting].Sponsor{
					var meeting_id []string
					meeting_id = append(meeting_id, all_meetings[meeting].Id)
					allUser[user].SponsorMeeting = remove_element_in_slice(allUser[user].SponsorMeeting,meeting_id)
				}
				for name := range all_meetings[meeting].Participants{
					if allUser[user].Username == all_meetings[meeting].Participants[name]{
						var meeting_id []string
						meeting_id = append(meeting_id,  all_meetings[meeting].Id)
						allUser[user].ParticipantMeeting = remove_element_in_slice(allUser[user].ParticipantMeeting,meeting_id)
					}
				}
			}
			delete(all_meetings, meeting)
			fmt.Println("the " + title + " meeting deleted success\n")
			count++
		}
	}

	if count == 0{
		fmt.Println("You haven't create or participant in the " + title + " meeting\n")
	}

	fout, _ := os.Create("data/Meeting.json")
	defer fout.Close()
	bytes,_ := json.Marshal(all_meetings)
	fout.Write(bytes)

	fout_user, _ := os.Create("data/User.json")
	defer fout_user.Close()
	buser, _ := json.Marshal(allUser)
	fout_user.Write(buser)
}

func ExitMeeting(title string) {
	current_user := User.GetCurUserName()
	if current_user == ""{
		fmt.Println("You haven't logged in")
		return 
	}

	var participants []string
	participants = append(participants, current_user)
	var allUser map[string]*User.User = GetAllUserInfo() 
	var all_meetings map[string]Meeting = GetAllMeetingInfo()
	count := 0

	for i := range all_meetings{
		if all_meetings[i].Title == title{
			slice := all_meetings[i].Participants
			slice = remove_element_in_slice(slice,participants)		
			var meet *Meeting = &Meeting{all_meetings[i].Title,all_meetings[i].Sponsor, slice, all_meetings[i].StartTime, all_meetings[i].EndTime, all_meetings[i].Id}
			all_meetings[i] = *meet

			if(len(all_meetings[i].Participants) == 0){
				for user := range allUser{
					if allUser[user].Username == all_meetings[i].Sponsor{
						var meeting_id []string
						meeting_id = append(meeting_id, meet.Id)
						allUser[user].SponsorMeeting = remove_element_in_slice(allUser[user].SponsorMeeting,meeting_id)
					}
				}
				delete(all_meetings, i)
			}
			//delete participants's meeting messages
			for user := range allUser{
				for name := range participants{
					if allUser[user].Username == participants[name]{
						var meeting_id []string
						meeting_id = append(meeting_id,  meet.Id)
						allUser[user].ParticipantMeeting = remove_element_in_slice(allUser[user].ParticipantMeeting,meeting_id)
						fmt.Println("exit succeed!\n")
						break
					}
				}
			}
			count++
		}
	}

	if count == 0{
		fmt.Println("You haven't participant in the " + title + " meeting\n")
		return
	}

	fout, _ := os.Create("data/Meeting.json")
	defer fout.Close()
	bytes,_ := json.Marshal(all_meetings)
	fout.Write(bytes)

	fout_user, _ := os.Create("data/User.json")
	defer fout_user.Close()
	buser, _ := json.Marshal(allUser)
	fout_user.Write(buser)
}

func DeleteAllMeetings() {
	current_user := User.GetCurUserName()
	if current_user == ""{
		fmt.Println("You haven't logged in")
		return 
	}

	var allUser map[string]*User.User = GetAllUserInfo() 
	var all_meetings map[string]Meeting = GetAllMeetingInfo()
	count := 0

	for meeting := range all_meetings{
		if all_meetings[meeting].Sponsor == current_user{
			for user := range allUser{
				if allUser[user].Username == all_meetings[meeting].Sponsor{
					var meeting_id []string
					meeting_id = append(meeting_id, all_meetings[meeting].Id)
					allUser[user].SponsorMeeting = remove_element_in_slice(allUser[user].SponsorMeeting,meeting_id)
				}
				for name := range all_meetings[meeting].Participants{
					if allUser[user].Username == all_meetings[meeting].Participants[name]{
						var meeting_id []string
						meeting_id = append(meeting_id,  all_meetings[meeting].Id)
						allUser[user].ParticipantMeeting = remove_element_in_slice(allUser[user].ParticipantMeeting,meeting_id)
					}
				}
			}
			delete(all_meetings, meeting)
			fmt.Println("all meetings of " + current_user + "  deleted success\n")
			count++
		}
	}

	if count == 0{
		fmt.Println("You haven't create the any meeting\n")
	}

	fout, _ := os.Create("data/Meeting.json")
	defer fout.Close()
	bytes,_ := json.Marshal(all_meetings)
	fout.Write(bytes)

	fout_user, _ := os.Create("data/User.json")
	defer fout_user.Close()
	buser, _ := json.Marshal(allUser)
	fout_user.Write(buser)
}

func  exitAllMeetings() {
	current_user := User.GetCurUserName()
	if current_user == ""{
		fmt.Println("You haven't logged in")
		return 
	}

	var allUser map[string]*User.User = GetAllUserInfo() 
	var all_meetings map[string]Meeting = GetAllMeetingInfo()
	var participants []string
	participants = append(participants, current_user)

	for i := range all_meetings{
		slice := all_meetings[i].Participants
		slice = remove_element_in_slice(slice,participants)		
		var meet *Meeting = &Meeting{all_meetings[i].Title,all_meetings[i].Sponsor, slice, all_meetings[i].StartTime, all_meetings[i].EndTime, all_meetings[i].Id}
		all_meetings[i] = *meet

		if(len(all_meetings[i].Participants) == 0){
			for user := range allUser{
				if allUser[user].Username == all_meetings[i].Sponsor{
					var meeting_id []string
					meeting_id = append(meeting_id, meet.Id)
					allUser[user].SponsorMeeting = remove_element_in_slice(allUser[user].SponsorMeeting,meeting_id)
				}
			}
			delete(all_meetings, i)
		}
	}
	fout, _ := os.Create("data/Meeting.json")
	defer fout.Close()
	bytes,_ := json.Marshal(all_meetings)
	fout.Write(bytes)

	fout_user, _ := os.Create("data/User.json")
	defer fout_user.Close()
	buser, _ := json.Marshal(allUser)
	fout_user.Write(buser)

}

func DeleteUser() {
	currentName := User.GetCurUserName()
	if currentName == "" {
		fmt.Println("You haven't logged in")
		return
	}
	var allUser map[string]*User.User = GetAllUserInfo()
	var userslice []string;
	userslice = append(userslice, currentName)
	count := 0

	exitAllMeetings()
	DeleteAllMeetings()

	for user := range allUser{
		if(allUser[user].Username == currentName){
			delete(allUser, user)
			fmt.Println("delete user success\n")
			count++
		}
	}

	if count == 0{
		fmt.Println("delete user fail\n")
	}

	User.LogOut()
	fout_user, _ := os.Create("data/User.json")
	defer fout_user.Close()
	buser, _ := json.Marshal(allUser)
	fout_user.Write(buser)

}