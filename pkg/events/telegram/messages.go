package telegram

const msgHelp = `Moi! I'm a bot, my name is Oranssi Pupu.
I can help you learn and practice the Suomi language.
Send me message /words and I will give you a random word and its translation in 3 seconds.
So you have time to guess the word🤓`

const msgHello = "Hello again! \n\n" + msgHelp

const (
	msgUnknownCommand  = "Unknown command. Type /help to get more information 🤨"
	msgNoDataInStorage = "Call the doctor my memory is empty! I don't know a single word 😱"
)
