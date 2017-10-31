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

func create_meeting(meeting *Meeting) {
	currentName := User.get_cur_user_name()
	if currentName == "" {
		fmt.Println("[agenda][error]you haven't logged in\n")
		return
	}

	flog, err := os.OpenFile("data/input_output.log", os.O_APPEND|os.O_WRONLY, 0600)
	defer flog.Close()
	check_err(err)
	logger := log.New(flog, "", log.LstdFlags)

	var allMeetings map[string]Meeting = get_all_meeting_info()
	if check_title_repeate(allMeetings, meeting.Title){
		return
	}

	meeting.Id = strconv.Itoa(len(allMeetings)) //initial id is 0
	meeting.Sponsor = currentName
	meeting.Participants = append(meeting.Participants)
	allMeetings[meeting.Id] = *meeting

	var map[string]*User.User = GetAllUserInfo()
	//check all participanter exist, and time contract
	for _, pName := range meeting.Participants {
		//check if the user exist
		if _, ok := allUser[pName]; !ok {
			fmt.Println("[agenda][error]no such user " + pName + "\n")
			return
		}
		//check if the user's old meeting is contract with the new one
		for _, meetingId := range allUser[pName].ParticipantMeeting {
			if TimeContact(meeting.StartTime, meeting.EndTime, allMeetings[meetingId].StartTime, allMeetings[meetingId].EndTime, pName, title) {
				return
			}
		}
		allUser[pName].ParticipantMeeting = append(allUser[pName].ParticipantMeeting, meeting.Id)
	}
	allUser[currentName].SponsorMeeting = append(allUser[currentName].SponsorMeeting, meeting.Id)

	logger.Printf("[agenda][info]%s create a meeting %+v\n", currentName, *meeting)
	
	write_meeting_json(allMeetings)
	write_user_json(allUser)	
}

func operate_participants(title, operation string, participants []string){
	current_user := User.GetCurUserName()
	if current_user == ""{
		fmt.Println("[agenda][error]you haven't logged in\n")
		return 
	}
	check_participants_exist(participants)
	var allUser map[string]*User.User = get_all_user_info()
	var all_meetings map[string]Meeting = get_all_meeting_info()
	count := 0
	for i := range all_meetings{
		if all_meetings[i].Title == title && all_meetings[i].Sponsor == current_user {
			if operation == "del" {
				allUser, all_meetings = del_participants(allUser, all_meetings, participants)
				break
			}
			if operation == "add" {
				allUser, all_meetings = add_participants_main(allUser, all_meetings, participants)
				break
			}
			count++
		}
	}
	if count == 0{
		fmt.Println("The meeting isn't exited\n")
		return 
	}
	write_meeting_json(allMeetings)
	write_user_json(allUser)
}

//时间处理函数
func date_to_string(date time.Time) string {
	return date.Format("2006-01-02/15:04")
}

func string_to_date(date string) (time.Time, error) {
	the_time, err := time.Parse("2006-01-02/15:04", date)
	return the_time, err
}

func small_date(date1, date2 time.Time) bool {
	return date1.Before(date2) || date1.Equal(date2)
}

func large_date(date1, date2 time.Time) bool {
	return date1.After(date2) || date1.Equal(date2)
}

func time_contact(newDateS, newDateE, oriDateS, oriDateE time.Time, userName string, title string) bool {
	//newstart time is after newend time, wrong!
	if large_date(newDateS, newDateE) {
		fmt.Println("[agenda][error]start time can't be greater than end time")
		return true
	}
	if small_date(newDateE, oriDateS) || LargeDate(newDateS, oriDateE) {
		return false
	}
	fmt.Println("[agenda][error]time contract with " + userName + "'s meeting with title:" + title)
	return true
}


func get_all_meeting_info() map[string]Meeting {
	byteIn, err := ioutil.ReadFile("data/Meeting.json")
	check_err(err)
	var allMeetingInfo map[string]Meeting
	json.Unmarshal(byteIn, &allMeetingInfo)
	if allMeetingInfo == nil {
		allMeetingInfo = make(map[string]Meeting)
	}
	return allMeetingInfo
}

func remove_sponermeetings(allUser map[string]*User.User, meeting Meeting, id string)map[string]*User.User {
	for user := range allUser{
		if allUser[user].Username == meeting.Sponsor{
			var meeting_id []string
			meeting_id = append(meeting_id, id)
			allUser[user].SponsorMeeting = remove_element_in_slice(allUser[user].SponsorMeeting,meeting_id)
		}
	}
	return allUser				
}

func remove_participantmeeting(allUser map[string]*User.User, participants []string, id string)map[string]*User.User {
	for user := range allUser{
		for name := range participants{
			if allUser[user].Username == participants[name]{
				var meeting_id []string
				meeting_id = append(meeting_id,  id)
				allUser[user].ParticipantMeeting = remove_element_in_slice(allUser[user].ParticipantMeeting,meeting_id)
				break
			}
		}
	}
	return allUser			
}

func remove_participant(meeting Meeting, participants []string) Meeting {
	slice := meeting.Participants
	slice = remove_element_in_slice(slice,participants)		
	var meet *Meeting = &Meeting{meeting.Title,meeting.Sponsor, slice, meeting.StartTime, meeting.EndTime, meeting.Id}
	return *meet
}

func slice_to_string(slice []string)string {
	slice_str := ""
	for i := range slice{
		slice_str = slice_str + " " + slice[i]
	}
	return slice_str
}

func del_participants(allUser map[string]*User.User, all_meetings map[string]Meeting, participants []string)(map[string]*User.User, map[string]Meeting) {
	all_meetings[i] = remove_participant(all_meetings[i], participants)
	allUser = remove_participantmeeting(allUser, participants, all_meetings[i].Id)
	if(len(all_meetings[i].Participants) == 0){
		allUser = remove_sponermeetings(allUser, all_meetings[i], all_meetings[i].Id)
		delete(all_meetings, i)
	}
	str_mes := slice_to_string(participants)
	fmt.Println("[agenda][info]delete the participants(" + str_mes + " ) in " + title + " meeting succeed!\n")
	logger.Print("[agenda][info]delete the participants(" + str_mes + " ) in " + title + " meeting succeed!\n")
	return allUser, all_meetings
}

func check_participants_repeat(participants []string, m_participants []string) bool{
	for j := range m_participants{
		for k := range participants{
			if m_participants == participants[k]{
				fmt.Println("[agenda][error]The user " + participants[k]+"is already in the meeting\n")
				return true
			}
		}
	}
	return false
}

func add_participant(meeting Meeting, participants []string) Meeting {
	slice := meeting.Participants
	slice = add_element_in_slice(slice,participants)		
	var meet *Meeting = &Meeting{meeting.Title,meeting.Sponsor, slice, meeting.StartTime, meeting.EndTime, meeting.Id}
	return *meet
}

func check_time_contract(allUser map[string]*User.User, meet Meeting, all_meetings map[string]Meeting, participants []string)bool {
	for _, pName := range participants{
		for _, meetingId := range allUser[pName].ParticipantMeeting {
			if TimeContact(meet.StartTime, meet.EndTime, all_meetings[meetingId].StartTime, all_meetings[meetingId].EndTime, pName, meetingId) {
				return true
			}
		}
	}
	return false
}

func add_participantmeeting(allUser map[string]*User.User, participants []string) map[string]*User.User {
	for user := range allUser{
		for name := range participants{
			if allUser[user].Username == participants[name]{
				allUser[user].ParticipantMeeting = append(allUser[user].ParticipantMeeting,all_meetings[i].Id)
				break;
			}
		}
	}
	return allUser
}

func add_participants_main(allUser map[string]*User.User, all_meetings map[string]Meeting, participants []string)(map[string]*User.User, map[string]Meeting) {
	if check_participants_repeat(participants, all_meetings[i].Participants){
		return
	}
	all_meetings[i] = add_participant(allUser, all_meetings, participants)
	if check_time_contract(allUser, allMeetings[i], allMeetings, participants){
		return
	}
	allUser = add_participantmeeting(allUser, participants)
	str_mes := slice_to_string(participants)
	fmt.Println("[agenda][info]delete the participants(" + str_mes + " ) in " + title + " meeting succeed!\n")
	logger.Print("[agenda][info]delete the participants(" + str_mes + " ) in " + title + " meeting succeed!\n")
	return allUser,all_meetings
}

func cancel_meeting(title string) {
	current_user := User.get_cur_user_name()
	if current_user == ""{
		fmt.Println("[agenda][error]you haven't logged in\n")
		return 
	}
	var allUser map[string]*User.User = get_all_user_info() 
	var all_meetings map[string]Meeting = get_all_meeting_info()
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
			fmt.Println("[agenda][info]the " + title + " meeting cancel successfully\n")
			logger.Print("[agenda][info]%s cancel the %s meeting\n", current_user, title)
			count++
		}
	}

	if count == 0{
		fmt.Println("[agenda][error]You haven't create or participant in the " + title + " meeting\n")
	}

	write_meeting_json(all_meetings)
	write_user_json(allUser)
}

func delete_all_meetings() {
	current_user := User.get_cur_user_name()
	if current_user == ""{
		fmt.Println("[agenda][error]you haven't logged in\n")
		return 
	}

	var allUser map[string]*User.User = get_all_user_info() 
	var all_meetings map[string]Meeting = get_all_meeting_info()
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
			fmt.Println("[agenda][info]all meetings of " + current_user + "  deleted success\n")
			logger.Print("[agenda][info]%s delete all his meeting\n", current_user)
			count++
		}
	}

	if count == 0{
		fmt.Println("You haven't create the any meeting\n")
	}

	write_meeting_json(allMeetings)
	write_user_json(allUser)
}

func exit_meeting(title string) {
	current_user := User.get_cur_user_name()
	if current_user == ""{
		fmt.Println("[agenda][error]you haven't logged in\n")
		return 
	}

	var participants []string
	participants = append(participants, current_user)
	var allUser map[string]*User.User = get_all_user_info() 
	var all_meetings map[string]Meeting = get_all_meeting_info()
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
						fmt.Println("[agenda][info]exit %s meeting succeed!\n", title)
						logger.Print("[agenda][info]%s exit %s meeting succeed!\n", current_user, title)
						break
					}
				}
			}
			count++
		}
	}

	if count == 0{
		fmt.Println("[agenda][error]You haven't participant in the " + title + " meeting\n")
		return
	}

	write_meeting_json(all_meetings)
	write_user_json(allUser)
}



func  exit_all_meetings() {
	current_user := User.get_cur_user_name()
	if current_user == ""{
		fmt.Println("You haven't logged in")
		return 
	}

	var allUser map[string]*User.User = get_all_user_info() 
	var all_meetings map[string]Meeting = get_all_meeting_info()
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
	write_meeting_json(all_meetings)
	write_user_json(allUser)
}


//load all user infomation to User.AllUserInfo
func get_all_user_info() map[string]*User.User {
	byteIn, err := ioutil.ReadFile("data/User.json")
	check(err)
	var allUserInfo map[string]*User.User
	json.Unmarshal(byteIn, &allUserInfo)
	return allUserInfo
}

func check_err(r error) {
	if r != nil {
		log.Fatal(r)
	}
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

func check_participants_exist(slice []string) {
	var allUser map[string]*User.User = get_all_user_info()
	for _, pname := range slice{
		if _, ok := allUser[pname];!ok{
			fmt.Println("[agenda][error]No such user:"+pname+"!")
			return
		}
	}
}

func search_meeting(search_start, search_end string) {
	current_user := User.get_cur_user_name()
	if current_user == ""{
		fmt.Println("You haven't logged in")
		return 
	}
	var all_meetings map[string]Meeting = get_all_meeting_info()
	start,_ := string_to_date(search_start)
	end,_ := string_to_date(search_end)

	if !small_date(start, end){
		fmt.Println("[agenda][error]the StartTime is larger than EndTime\n")
		return
	}

	count := 0
	for meeting := range all_meetings{
		meet := all_meetings[meeting]
		if check_user_in_meeting(&meet, current_user){
			if !(small_date(meet.EndTime, start) || small_date(end, meet.StartTime)){
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




func delete_user() {
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
	logger.Print("[agenda][info]%s deleted", currentName)
	write_user_json(allUser)
}

func write_meeting_json(allMeetings map[string]Meeting) {
	fout, _ := os.Create("data/Meeting.json")
	defer fout.Close()
	b, _ := json.Marshal(allMeetings)
	fout.Write(b)
}

func write_user_json(allUser map[string]*User.User) {
	foutuser, _ := os.Create("data/User.json")
	defer foutuser.Close()
	buser, _ := json.Marshal(allUser)
	foutuser.Write(buser)
}

func check_title_repeate(allMeetings map[string]Meeting, title string) bool {
	for index := range allMeetings{
		if allMeetings[index].Title == title{
			fmt.Println("[agenda][error]the title has exited\n")
			return true
		}
	}
	return false
}