package User

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

//User is used to encapsulate all functions
//that deal with User-relativedlogic
type User struct {
	Username           string
	Password           string
	Email              string
	SponsorMeeting     []string
	ParticipantMeeting []string
}

//register the user with name, password, email
func register_user(user *User) {
	AllUserInfo := get_all_user_Info()
	flog, err := os.OpenFile("data/input_output.log", os.O_APPEND|os.O_WRONLY, 0600)
	defer flog.Close()
	check_err(err)
	logger := log.New(flog, "", log.LstdFlags)

	if _, ok := AllUserInfo[user.Username]; !ok {
		AllUserInfo[user.Username] = *user
		os.Stdout.WriteString("[agenda][info]"+ user.Username + "registed succeed!\n")
		logger.Printf("[agenda][info] %s registed succeed\n", user.Username)
	} else {
		os.Stdout.WriteString("[agenda][warning]The userName " + user.Username + "have been registered\n")
	}

	fout, _ := os.Create("data/User.json")
	defer fout.Close()
	b, _ := json.Marshal(AllUserInfo)
	fout.Write(b)
}

//search all user
func search_all_user() {
	AllUserInfo := get_all_user_info()
	flog, err := os.OpenFile("data/input_output.log", os.O_APPEND|os.O_WRONLY, 0600)
	defer flog.Close()
	check_err(err)
	logger := log.New(flog, "", log.LstdFlags)
	logger.Printf("[agenda][info]" + get_cur_user_name() + " search all users")

	if get_cur_user_name() == "" {
		os.Stdout.WriteString("You haven't logged in, can't search for all users!\n")
		return
	} else {
		for _, val := range AllUserInfo {
			fmt.Println("[agenda][info]name: %s, email: %s", val.Username, val.Email)
		}
	}
}

//log in with name, password
func logIn(user *User) {
	AllUserInfo := get_all_user_info()

	flog, err := os.OpenFile("data/input_output.log", os.O_APPEND|os.O_WRONLY, 0600)
	defer flog.Close()
	check_err(err)
	logger := log.New(flog, "", log.LstdFlags)

	if get_cur_user_name() != "" {
		os.Stdout.WriteString("[agenda][warning]You have log in already\n")
		return
	}
	if _, ok := AllUserInfo[user.Username]; !ok {
		os.Stdout.WriteString("[agenda][error]Username or password is incorrect!\n")
	} else {
		correctPass := AllUserInfo[user.Username].Password
		if correctPass == user.Password {
			fout, _ := os.Create("data/current.txt")
			defer fout.Close()
			fout.WriteString(user.Username)
			os.Stdout.WriteString("[agenda][info]"+user.Username+" haved logged in\n")
			logger.Print("[agenda][info]"+user.Username+" haved logged in\n")
		} else {
			os.Stdout.WriteString("[agenda][error]Username or password is incorrect!\n")
		}
	}
}

//log out with name, password
func logout() {
	//	AllUserInfo := GetAllUserInfo()
	curuser = get_cur_user_name()
	flog, err := os.OpenFile("data/input_output.log", os.O_APPEND|os.O_WRONLY, 0600)
	defer flog.Close()
	check_err(err)
	logger := log.New(flog, "", log.LstdFlags)

	if get_cur_user_name() == "" {
		os.Stdout.WriteString("[agenda][error]You haven't logged in!\n")
	} else {
		fout, _ := os.Create("data/current.txt")
		defer fout.Close()
		fout.WriteString("")
		os.Stdout.WriteString("[agenda][info]" + curuser +" logged out!\n")
		logger.Print("[agenda][info]" + curuser +" logged out!\n")
	}
}

//load all user infomation to User.AllUserInfo
func get_all_user_info() map[string]User {

	byteIn, err := ioutil.ReadFile("data/User.json")
	check_err(err)
	var allUserInfo map[string]User
	json.Unmarshal(byteIn, &allUserInfo)
	return allUserInfo
}

func get_cur_user_name() string {
	fin, err0 := os.Open("data/current.txt")
	check_err(err0)
	defer fin.Close()
	reader := bufio.NewReader(fin)
	str, _ := reader.ReadString('\n')
	return str
}

func check_err(r error) {
	if r != nil {
		log.Fatal(r)
	}
}


