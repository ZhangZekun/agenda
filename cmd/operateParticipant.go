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

// operateParticipantCmd represents the operateParticipant command
var operateParticipantCmd = &cobra.Command{
	Use:   "operateParticipant",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("operateParticipant called")
		title, _ := cmd.Flags().GetString("title")
		operation, _ := cmd.Flags().GetString("op")
		participants, _ := cmd.Flags().GetStringSlice("participants")
		Meeting.OperateParticipants(title,operation, participants)
	},
}

func init() {
	RootCmd.AddCommand(operateParticipantCmd)
	operateParticipantCmd.Flags().StringP("title", "t", "", "title")
	operateParticipantCmd.Flags().StringP("op", "o", "", "Operation To Participant")
	operateParticipantCmd.Flags().StringSliceP("participants", "p", make([]string, 0), "Names of Participants")



	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// operateParticipantCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// operateParticipantCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
