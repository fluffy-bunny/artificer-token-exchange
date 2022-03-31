/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package secrets

import (
	core_utils "github.com/fluffy-bunny/grpcdotnetgo/pkg/utils"

	"encoding/base64"
	"fmt"

	"github.com/gorilla/securecookie"
	"github.com/spf13/cobra"
)

// genCmd represents the gen command
var secretCmd = &cobra.Command{
	Use:   "secrets",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		key := securecookie.GenerateRandomKey(length)
		encodedString := base64.StdEncoding.EncodeToString(key)

		type output struct {
			Secret string
			Length int
			Encode string // base64 encoded
		}
		fmt.Println(core_utils.PrettyJSON(&output{
			Secret: encodedString,
			Length: length,
			Encode: "base64",
		}))
		fmt.Println("secrets called")
	},
}
var length int

func Init(rootCmd *cobra.Command) {
	rootCmd.AddCommand(secretCmd)

	secretCmd.Flags().IntVarP(&length, "length", "l", 0, "--length=[32|64]")
	secretCmd.MarkFlagRequired("length")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
