package telegram

const msgHelp = `I'm a bot, my name is Oranssi Pupu. I can save and keep your links.
	Also I can send you to reed a random link from previously saved links.

	In order to save, just send a link to me.

	In order to get a random link, send me a command /rnd.
	Be aware! That page will be deleted from my memory after sending it to you.`

const msgHello = "Hello again! \n\n" + msgHelp

const (
	msgUnknownCommand = "Unknown command. Type /help to get more information."
	msgNoSavedPages   = "You don't have any saved links. Send me a link to save it."
	msgPageSaved      = "Link saved."
	msgAlreadyExists  = "Link already saved."
)
