package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

func getDatabaseArguments() (bool, []string) {
	readerDatabaseCustom := bufio.NewReader(os.Stdin)
	fmt.Print("Do you want to use your own database? [Y or N]")
	databaseCustom, _ := readerDatabaseCustom.ReadString('\n')
	if strings.TrimRight(strings.ToUpper(databaseCustom), "\n") == "N" {
		return false, nil
	}
	var databaseInfo []string
	readerDatabaseHost := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your database host: ")
	databaseHost, _ := readerDatabaseHost.ReadString('\n')
	databaseInfo = append(databaseInfo, databaseHost)

	readerDatabasePort := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your database port: ")
	databasePort, _ := readerDatabasePort.ReadString('\n')
	databaseInfo = append(databaseInfo, databasePort)

	readerDatabaseUser := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your databaseUser: ")
	databaseUser, _ := readerDatabaseUser.ReadString('\n')
	databaseInfo = append(databaseInfo, databaseUser)

	readerDatabasePassword := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your database password: ")
	databasePassword, _ := readerDatabasePassword.ReadString('\n')
	databaseInfo = append(databaseInfo, databasePassword)

	fmt.Println("\n\nhost:", databaseInfo[0])
	fmt.Println("port: ", databaseInfo[1])
	fmt.Println("username: ", databaseInfo[2])
	fmt.Println("password: ", databaseInfo[3])

	readerDatabaseConfirmation := bufio.NewReader(os.Stdin)
	fmt.Print("\nThe database information is correct? [Y or N]")
	databaseConfirmation, _ := readerDatabaseConfirmation.ReadString('\n')
	fmt.Print(databaseConfirmation)
	if strings.TrimRight(strings.ToUpper(databaseConfirmation), "\n") == "Y" {
		return true, databaseInfo
	}
	return getDatabaseArguments()

}
func getIngressArguments() error {
	return nil
}

func getVersionArguments() error {
	return nil
}

func renderYaml() error {
	return nil
}

func installCharles(c *cli.Context) error {
	useCustomDatabase, databaseInfo := getDatabaseArguments()
	fmt.Println(useCustomDatabase, databaseInfo)

	getIngressArguments()
	getVersionArguments()
	renderYaml()

	return nil
}

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "install",
				Aliases: []string{"i"},
				Usage:   "Start charles install",
				Action:  installCharles,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
