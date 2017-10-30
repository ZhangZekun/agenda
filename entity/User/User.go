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
func RegisterAnUser(user *User) {
	AllUserInfo := GetAllUserInfo()
	flog, err := os.OpenFile("data/input_output.log", os.O_APPEND|os.O_WRONLY, 0600)
	defer flog.Close()
	check(err)
	logger := log.New(flog, "", log.LstdFlags)
	logger.Printf("agenda register -u %s -p %s -e %s", user.Username, user.Password, user.Email)

	if _, ok := AllUserInfo[user.Username]; !ok {
		AllUserInfo[user.Username] = *user
		os.Stdout.WriteString("register succeed!\n")
		logger.Print("register succeed!\n")
	} else {
		os.Stdout.WriteString("The userName have been registered\n")
		logger.Print("The userName have been registered\n")
	}

	fout, _ := os.Create("data/User.json")
	defer fout.Close()
	b, _ := json.Marshal(AllUserInfo)
	fout.Write(b)
}

//search all user
func SearchAllUser() {
	AllUserInfo := GetAllUserInfo()
	flog, err := os.OpenFile("data/input_output.log", os.O_APPEND|os.O_WRONLY, 0600)
	defer flog.Close()
	check(err)
	logger := log.New(flog, "", log.LstdFlags)
	logger.Printf("agenda searchUser")

	if GetCurUserName() == "" {
		os.Stdout.WriteString("You haven't logged in, can't search for all users!\n")
		return
	} else {
		for _, val := range AllUserInfo {
			fmt.Println(val.Username, val.Password, val.Email)
		}
	}
}

//log in with name, password
func LogIn(user *User) {
	AllUserInfo := GetAllUserInfo()

	flog, err := os.OpenFile("data/input_output.log", os.O_APPEND|os.O_WRONLY, 0600)
	defer flog.Close()
	check(err)
	logger := log.New(flog, "", log.LstdFlags)
	logger.Printf("agenda login -u %s -p %s", user.Username, user.Password)

	if GetCurUserName() != "" {
		os.Stdout.WriteString("You have log in already!\n")
		logger.Print("You have log in already!\n")
		return
	}
	if _, ok := AllUserInfo[user.Username]; !ok {
		os.Stdout.WriteString("Username or password is incorrect!\n")
		logger.Print("Username or password is incorrect!\n")
	} else {
		correctPass := AllUserInfo[user.Username].Password
		if correctPass == user.Password {
			fout, _ := os.Create("data/current.txt")
			defer fout.Close()
			fout.WriteString(user.Username)
			os.Stdout.WriteString("you haved logged in\n")
			logger.Print("you haved logged in\n")
		} else {
			os.Stdout.WriteString("Username or password is incorrect!\n")
			logger.Print("Username or password is incorrect!\n")
		}
	}
}

//log out with name, password
func LogOut() {
	//	AllUserInfo := GetAllUserInfo()
	flog, err := os.OpenFile("data/input_output.log", os.O_APPEND|os.O_WRONLY, 0600)
	defer flog.Close()
	check(err)
	logger := log.New(flog, "", log.LstdFlags)
	logger.Printf("agenda logout")

	if GetCurUserName() == "" {
		os.Stdout.WriteString("You haven't logged in!\n")
		logger.Print("You haven't logged in!\n")
	} else {
		fout, _ := os.Create("data/current.txt")
		defer fout.Close()
		fout.WriteString("")
		os.Stdout.WriteString("you have logged out!\n")
		logger.Print("you have logged out!\n")
	}
}

//load all user infomation to User.AllUserInfo
func GetAllUserInfo() map[string]User {

	byteIn, err := ioutil.ReadFile("data/User.json")
	check(err)
	var allUserInfo map[string]User
	json.Unmarshal(byteIn, &allUserInfo)
	return allUserInfo
}

func GetCurUserName() string {
	fin, err0 := os.Open("data/current.txt")
	check(err0)
	defer fin.Close()
	reader := bufio.NewReader(fin)
	str, _ := reader.ReadString('\n')
	return str
}

func check(r error) {
	if r != nil {
		log.Fatal(r)
	}
}
