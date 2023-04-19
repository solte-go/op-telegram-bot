package telegram

const msgHelp = `Moi! I'm a bot, my name is Oranssi Pupu.
I can help you learn and practice the Suomi language.
Send me message /words and I will give you a random word and its translation in 3 seconds.
So you have time to guess the word🤓
You can also send me /setLang <ru/en> to change the language.
If you want to learn a specific topic, send me /setTopic <topic name> and I will give you a random word from this topic.
To see all available topics, send me /topics.
`

const msgHello = "Hello again! \n\n" + msgHelp

const (
	msgUnknownCommand   = "Unknown command. Type /help to get more information 🤨"
	msgNoDataInStorage  = "Call the doctor my memory is empty! I don't know a single word 😱"
	msgMissingArgument  = "Missing argument 🤨. To setup the language, type /setLang <ru/en>. Type /help to get more information"
	msgUnsupportedLang  = "Unsupported language argument 🤨. To setup the language, type /setLang <ru/en>. Type /help to get more information"
	msgUnsupportedTopic = "No topic with provided name 🤨. To see all available topics, type /topics. Type /help to get more information"
	msgSettingApplied   = "Settings applied ✅"
)
