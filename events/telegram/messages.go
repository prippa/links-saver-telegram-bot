package telegram

const msgHelp = `I can save your links. Just send me a link and I will save it for you.

In order to get your saved links, just type /rand command.
After that this page will be deleted from the list.
`

const msgStart = "Hello and welcome! \n\n" + msgHelp

const (
	msgUnknownCmd    = "🤔 Unknown command. Type /help to get help"
	msgNoSavedPages  = "You have no saved pages 🤷‍♂️"
	msgSaved         = "Saved! ✅"
	msgAlreadyExists = "This page already exists in your list 🤗"
)
