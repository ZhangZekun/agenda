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
	Phone			   string
	SponsorMeeting     []string
	ParticipantMeeting []string
}

//register the user with name, password, email
func Register_user(user *User) {
	if(user.Username == "" || user.Password == "" || user.Email == "" || user.Phone == ""){
		fmt.Println("you need to input user message to create an account, the arguments include usernam, password, email, telephone.For example:\n./agenda register -u=zhangzemian -p=12345678 -e=1106066690@qq.com -t=15018377821\n")
		return
	}
	AllUserInfo := Get_all_user_info()
	flog, err := os.OpenFile("data/input_output.log", os.O_APPEND|os.O_WRONLY, 0600)
	defer flog.Close()
	check_err(err)
	logger := log.New(flog, "", log.LstdFlags)

	if _, ok := AllUserInfo[user.Username]; !ok {
		AllUserInfo[user.Username] = *user
		os.Stdout.WriteString("[agenda][info] "+ user.Username + " registed succeed!\n")
		logger.Printf("[agenda][info] "+user.Username+" registed succeed\n")
	} else {
		os.Stdout.WriteString("[agenda][warning]The userName " + user.Username + " have been registered\n")
		return
	}

	fout, _ := os.Create("data/User.json")
	defer fout.Close()
	b, _ := json.Marshal(AllUserInfo)
	fout.Write(b)
}

//search all user
func Search_all_user() {
	AllUserInfo := Get_all_user_info()
	flog, err := os.OpenFile("data/input_output.log", os.O_APPEND|os.O_WRONLY, 0600)
	defer flog.Close()
	check_err(err)
	logger := log.New(flog, "", log.LstdFlags)
	logger.Printf("[agenda][info]" + Get_cur_user_name() + " search all users")

	if Get_cur_user_name() == "" {
		os.Stdout.WriteString("You haven't logged in, can't search for all users!\n")
		return
	} else {
		for _, val := range AllUserInfo {
			fmt.Println("[agenda][info]name:" + val.Username + " , email:" + val.Email+"\n")
		}
	}
}

//log in with name, password
func Login(user *User) {
	if(user.Username == "" || user.Password==""){
		fmt.Println("you need to input the username and password,for example:\n./agenda login -u=zhangzemian -p=12345678\n")
		return
	}
	AllUserInfo := Get_all_user_info()
	flog, err := os.OpenFile("data/input_output.log", os.O_APPEND|os.O_WRONLY, 0600)
	defer flog.Close()
	check_err(err)
	logger := log.New(flog, "", log.LstdFlags)

	if Get_cur_user_name() != "" {
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
func Logout() {
	//	AllUserInfo := GetAllUserInfo()
	curuser := Get_cur_user_name()
	flog, err := os.OpenFile("data/input_output.log", os.O_APPEND|os.O_WRONLY, 0600)
	defer flog.Close()
	check_err(err)
	logger := log.New(flog, "", log.LstdFlags)

	if Get_cur_user_name() == "" {
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
func Get_all_user_info() map[string]User {

	byteIn, err := ioutil.ReadFile("data/User.json")
	check_err(err)
	var allUserInfo map[string]User
	json.Unmarshal(byteIn, &allUserInfo)
	return allUserInfo
}

func Get_cur_user_name() string {
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


