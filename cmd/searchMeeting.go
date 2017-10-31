// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"agenda/entity/Meeting"
)

// searchMeetingCmd represents the searchMeeting command
var searchMeetingCmd = &cobra.Command{
	Use:   "searchMeeting",
	Short: "search the meeting you create or participate in",
	Long: `you should input two arguments of the command.For example:
	./agenda searchMeeting -s=2013-10-10/10:00 -e=2017-10-10/10:00`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("searchMeeting called")
		startTime,_ := cmd.Flags().GetString("startTime")
		endTime,_ := cmd.Flags().GetString("endTime")
		Meeting.search_meeting(startTime, endTime)
	},
}

func init() {
	RootCmd.AddCommand(searchMeetingCmd)
	searchMeetingCmd.Flags().StringP("statTime", "s", "", "start time of meeting")
	searchMeetingCmd.Flags().StringP("endTime", "e", "", "end time of meeting")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchMeetingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchMeetingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
