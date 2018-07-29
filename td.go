package td

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"strconv"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli"
)

var storage tdData

type tdData struct {
	Lists       map[string][]tdTodo
	CurrentList string
}

type tdTodo struct {
	Text string
	Done bool
}

func main() {

	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Author = "Sam Carlile"
	app.Version = "1.0.0"
	app.Name = "td"
	app.Usage = "Get your stuff done"

	app.Before = func(c *cli.Context) error {
		storage = loadData()
		return nil
	}

	app.After = func(c *cli.Context) error {
		return saveData(storage)
	}

	app.Action = func(c *cli.Context) error {
		if c.NArg() == 0 {
			displayList(storage.CurrentList)
		} else {
			if storage.Lists[c.Args().First()] != nil {
				displayList(c.Args().First())
			}
		}
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:      "list",
			Usage:     "Select the working list or display the lists",
			ArgsUsage: "[name]",
			Aliases:   []string{"l"},
			Action: func(c *cli.Context) error {
				listName := c.Args().First()

				if len(listName) == 0 {
					displayLists()
				} else {
					storage.CurrentList = listName

					if storage.Lists[listName] != nil {
						color.Yellow("Switched to list: %s", listName)
					} else {
						storage.Lists[listName] = []tdTodo{}
						color.Green("Created new list: %s", listName)
					}
				}
				return nil
			},
			BashComplete: func(c *cli.Context) {

				if c.NArg() > 0 { // If you've already types `td l listname`, this prevents further listname suggestions
					return
				}

				for l := range storage.Lists {
					color.New(color.FgYellow).Println(l)
				}
			},
		},
		{
			Name:      "add",
			Usage:     "Add a task to the working list",
			ArgsUsage: "take out the trash",
			Aliases:   []string{"a"},
			Action: func(c *cli.Context) error {
				todo := ""
				for i := 0; i < c.NArg(); i++ {
					todo += c.Args().Get(i)
					if i != c.NArg()-1 {
						todo += " "
					}
				}

				if len(todo) == 0 {
					return nil
				}

				storage.Lists[storage.CurrentList] = append(
					storage.Lists[storage.CurrentList],
					tdTodo{Text: todo, Done: false},
				)

				displayUpdateAdded(storage.CurrentList, len(storage.Lists[storage.CurrentList])-1)
				return nil
			},
		},
		{
			Name:    "delete",
			Aliases: []string{"d"},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "list, l",
					Usage: "delete a list by name",
				},
			},
			Action: func(c *cli.Context) error {

				if storage.Lists[c.String("list")] != nil {
					prompt := promptui.Prompt{
						IsConfirm: true,
						Label:     "Delete list " + c.String("list"),
					}
					result, err := prompt.Run()

					if err != nil {
						return err
					}

					if result == "N" {
						return nil
					}

					delete(storage.Lists, c.String("list"))
					color.New(
						color.FgRed,
						color.Bold,
					).Println("Deleted list: ", c.String("list"))
					return nil
				}

				rawID := c.Args().First()
				if len(rawID) == 0 {
					return nil
				}

				id, err := strconv.Atoi(rawID)
				if err != nil {
					return errors.New("Provide a todo number to delete it.")
				}

				if id > len(storage.Lists[storage.CurrentList])-1 || id < 0 {
					return nil
				}

				displayUpdateDeleted(storage.CurrentList, id)
				storage.Lists[storage.CurrentList] = append(storage.Lists[storage.CurrentList][:id], storage.Lists[storage.CurrentList][id+1:]...)
				return nil
			},
		},
		{
			Name:    "complete",
			Aliases: []string{"c"},
			Action: func(c *cli.Context) error {
				rawID := c.Args().First()
				if len(rawID) == 0 {
					return nil
				}

				id, err := strconv.Atoi(rawID)
				if err != nil {
					return errors.New("provide a todo number to delete it")
				}

				if id > len(storage.Lists[storage.CurrentList])-1 || id < 0 {
					return nil
				}

				storage.Lists[storage.CurrentList][id].Done = !storage.Lists[storage.CurrentList][id].Done
				displayUpdateCompleted(storage.CurrentList, id)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// Utility Functions

func displayLists() {

	color.New(
		color.FgWhite,
		color.Bold,
		color.Underline,
	).Println("Lists")

	for listName := range storage.Lists {
		if storage.CurrentList == listName {
			color.New(
				color.FgYellow,
				color.Bold,
			).Println("-", listName, "(active)")
		} else {
			color.New(
				color.FgYellow,
			).Println("-", listName)
		}

	}
}

func displayListName(listName string) {
	color.New(
		color.FgBlue,
		color.Underline,
		color.Bold,
	).Println(listName)
}

func displayList(name string) {
	displayListName(name)

	for index, todo := range storage.Lists[name] {
		if todo.Done {
			displayTodoCompleted(todo.Text, index)
		} else {
			displayTodo(todo.Text, index)
		}
	}
}

func displayTodo(text string, id int) {
	fmt.Println(
		color.New(color.FgWhite).Sprintf("%d", id),
		"-",
		color.New(color.Bold).Sprintf(text),
	)
}

func displayTodoCompleted(text string, id int) {
	fmt.Println(
		color.New(color.FgWhite).Sprintf("%d - %s", id, text),
		color.New(color.FgMagenta).Sprintf("(done)"),
	)
}

func displaySingleDeletedTodo(text string) {
	r := color.New(color.FgRed, color.Bold).SprintfFunc()
	fmt.Println(
		"❌ ",
		r(text),
	)
}

func displaySingleAddedTodo(text string) {
	r := color.New(color.FgMagenta, color.Bold).SprintfFunc()
	fmt.Println(
		"✨ ",
		r(text),
	)
}

func displaySingleCompletedTodo(text string) {
	r := color.New(color.FgGreen, color.Bold).PrintlnFunc()
	r("✔ " + text)
}

func displayUpdateDeleted(listName string, id int) {
	displayListName(listName)

	for index, todo := range storage.Lists[listName] {
		if index == id {
			displaySingleDeletedTodo(todo.Text)
		} else {
			newIndex := index
			if index >= id {
				newIndex--
			}
			displayTodo(todo.Text, newIndex)
		}
	}
}

func displayUpdateAdded(listName string, id int) {
	displayListName(listName)

	for index, todo := range storage.Lists[listName] {
		if index == id {
			displaySingleAddedTodo(todo.Text)
		} else {
			displayTodo(todo.Text, index)
		}
	}
}

func displayUpdateCompleted(listName string, id int) {
	displayListName(listName)

	for index, todo := range storage.Lists[listName] {
		if index == id {
			displaySingleCompletedTodo(todo.Text)
		} else {
			displayTodo(todo.Text, index)
		}
	}
}

/// Config stuff

// type tdConfig struct {
// 	listPath string `json:"listPath"`
// }

func loadData() tdData {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	filename := path.Join(usr.HomeDir, ".td_data")

	if _, err := os.Stat(filename); os.IsNotExist(err) {

		d := tdData{
			Lists: map[string][]tdTodo{
				"todos": []tdTodo{},
			},
			CurrentList: "todos",
		}
		return d
	}

	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading todo data file: %s", err)
	}

	var data tdData
	if err := json.Unmarshal(raw, &data); err != nil {
		log.Fatalf("Error parsing todo data file: %s", err)
	}

	return data

}

func saveData(data tdData) error {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path.Join(usr.HomeDir, ".td_data"), raw, 0644)
}
