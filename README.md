# td - *A nice simple CLI to-do list manager*

[![asciicast](https://asciinema.org/a/194005.png)](https://asciinema.org/a/194005)

---

## Getting Started
The first time you run `td`, it displays and creates a default list named *todo*. This list becomes the active list by default.

To add items to the active list, run `td a take out the trash` or `td add take out the trash`.

To delete a todo, run `td d <id>`.

To complete a todo, run `td c <id>`

That's pretty much it. You can create new lists or change the active list with `td l my-new-list` or `td list my-new-list`.

Run `td help` for the help page.

---

## Future Features
 - [ ] add option to move completed items to bottom of list when displaying list
 - [ ] customize display settings (colors, style, etc)
 - [ ] storage adapters (right now the default storage system is just a JSON file in your home directory)

--- 

## FAQ

### Where does it store the data?
Right now all the data is stored in JSON at `$HOME/.td_data`