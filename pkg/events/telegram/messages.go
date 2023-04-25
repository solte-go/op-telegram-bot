package telegram

const msgHelp = `
Oranssi Pupu Here!
/helpEn - Help in English
/helpRu - Помощь на Русском
`

const msgHello = "Moi! \n" + msgHelp

const (
	msgUnknownCommand   = "Unknown command. Type /help to get more information 🤨"
	msgNoDataInStorage  = "Call the doctor my memory is empty! I don't know a single word 😱"
	msgMissingArgument  = "Missing argument 🤨. To setup the language, type /setLang <ru/en>. Type /help to get more information"
	msgUnsupportedLang  = "Unsupported language argument 🤨. To setup the language, type /setLang <ru/en>. Type /help to get more information"
	msgUnsupportedTopic = "No topic with provided name 🤨. To see all available topics, type /topics. Type /help to get more information"
	msgSettingApplied   = "Settings applied ✅"
)

const msgHelpRu = `Привет!
Я бот, меня зовут Oranssi Pupu. Я могу помочь тебе учить и практиковать финский язык.

/words - Отправь мне это сообщение, и я пришлю тебе случайное слово и его перевод через 3 секунды. Так что у тебя будет время угадать слово🤓

/setLang <ru/en> - Для смены языка, на котором будут присылаться ответы.

/setTopic <название темы> - Это сообщение, если ты хочешь попрактиковать конкретную тему. Я буду присылать случайные слова из этой темы.

/topics - Чтобы увидеть все доступные темы.

Пример: /setTopic animals
Пример: /setLang en

Help in English: /helpEn
`

const msgHelpEn = `Moi! 
I'm a bot, my name is Oranssi Pupu. I can help you learn and practice Suomi language.

/words - Send me this message and I'll send you a random word and it's translation in 3 seconds. So you will have time to guess the word🤓
		 
/setLang <ru/en> - To change the language in which answers will be sent.

/setTopic <topic name> - Send me this message if you want to practice a specific topic. I'll start sending you random words from this topic.

/topics - To see all available topics send me "/topics".

Example: /setTopic animals
Example: /setLang en

Помощь на Русском: /helpRu
`
