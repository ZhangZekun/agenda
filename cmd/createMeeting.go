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
	"agenda/entity/Meeting"
	"fmt"
	"github.com/spf13/cobra"
	//"time"
)

// createMeetingCmd represents the createMeeting command
var createMeetingCmd = &cobra.Command{
	Use:   "createMeeting",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("createMeeting called")
		title, _ := cmd.Flags().GetString("title")
		participants, _ := cmd.Flags().GetStringSlice("participants")
		startTime, _ := cmd.Flags().GetString("startTime")
		endTime, _ := cmd.Flags().GetString("endTime")
		timeS, _ := Meeting.StringToDate(startTime)
		timeE, _ := Meeting.StringToDate(endTime)
		Meeting.CreateAMeeting(&Meeting.Meeting{title, "", participants, timeS, timeE, ""})
	},
}

func init() {
	RootCmd.AddCommand(createMeetingCmd)
	createMeetingCmd.Flags().StringP("title", "t", "", "title")
	createMeetingCmd.Flags().StringSliceP("participants", "p", make([]string, 0), "title")
	createMeetingCmd.Flags().StringP("startTime", "s", "", "startTime")
	createMeetingCmd.Flags().StringP("endTime", "e", "", "User endTime")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createMeetingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createMeetingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
