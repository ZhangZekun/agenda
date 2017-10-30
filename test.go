package main

import (
	//"bufio"
	//	"encoding/json"
	"fmt"
	//"log"
	//	"io/ioutil"
	//"os"
	"time"
)

type User struct {
	Username string
	Password string
	Email    string
}

//时间处理函数
func DateToString(date time.Time) string {
	return date.Format("2006-01-02 15:04")
}
func StringToDate(date string) (time.Time, error) {
	the_time, err := time.Parse("2006-01-02 15:04", date)
	return the_time, err
}

func main() {
	var str string
	str = "2017-10-30 11:48"
	t, _ := StringToDate(str)
	newStr := DateToString(t)
	fmt.Println(newStr)
	/*
		byteIn, _ := ioutil.ReadFile("myUser.json")
		var info map[string]User
		json.Unmarshal(byteIn, &info)
		fmt.Println(info)
		for key, val := range info {
			fmt.Println(key, val)
			fmt.Println(val.Password)
		}
		)*/

	/*
		fin, err := os.Open("myUser.json")
		defer fin.Close()
		byteIn := make([]byte, 3096)
		reader := bufio.NewReader(fin)
		byteRead, _ := reader.Read(byteIn)

		var info interface{}
		os.Stdout.Write(byteIn)

		err2 := json.Unmarshal(byteIn[:byteRead], &info)
		if err2 != nil {
			//	os.Stdout.Write([]byte("err in unmarshal"))
			log.Fatal(err2)
		}

		m, ok := info.(map[string]interface{})

		if ok == true {
			os.Stdout.Write([]byte("success to change"))
		}
		age := m["Name"].(string)
		os.Stdout.WriteString(age)
	*/
	//os.Stdout.Write(info)

	/*
		fout, _ := os.Create("myUser.json")
		defer fout.Close()

		f := map[string]interface{}{
			"a": User{"a", "b", "c"},
			"b": User{"z", "x", "c"},
		}

		b, _ := json.Marshal(f)
		fout.Write(b)
	*/
}
