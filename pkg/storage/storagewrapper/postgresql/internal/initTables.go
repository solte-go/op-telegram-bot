package internal

func CreateTables() string {
	return `CREATE TABLE IF NOT EXISTS public.users (
    	id SERIAL PRIMARY KEY,
    	user_name VARCHAR NOT NULL CONSTRAINT "Users_pk" UNIQUE,
    	topic VARCHAR NOT NULL DEFAULT 'all', 
    	user_language VARCHAR NOT NULL DEFAULT 'ru',
    	seq_offset INTEGER NOT NULL DEFAULT 0,
    	seq_interval INTEGER NOT NULL DEFAULT 3
	);

	ALTER TABLE users
    	OWNER TO postgres;
	
	CREATE TABLE IF NOT EXISTS public.excluded_words (
	  	id SERIAL PRIMARY KEY,
	  	user_name VARCHAR NOT NULL,
	  	word VARCHAR NOT NULL CONSTRAINT "excluded_words_PK" UNIQUE,
	  	create_at timestamp DEFAULT CURRENT_TIMESTAMP
	);
	
	ALTER TABLE excluded_words
		OWNER TO postgres;

	CREATE TABLE IF NOT EXISTS public.words (
	    id SERIAL PRIMARY KEY,
	    letter VARCHAR NOT NULL,
	    topic VARCHAR NOT NULL,
	    suomi VARCHAR NOT NULL CONSTRAINT "words_PK" UNIQUE,
	    english VARCHAR NOT NULL,
	    russian VARCHAR NOT NULL,
	    create_at timestamp DEFAULT CURRENT_TIMESTAMP
	);

	ALTER TABLE words
    	OWNER TO postgres;
`
}

//CREATE TABLE IF NOT EXISTS public.links (
//id SERIAL PRIMARY KEY,
//user_id INTEGER NOT NULL,
//link TEXT NOT NULL CONSTRAINT "Links_pk" UNIQUE,
//create_at timestamp DEFAULT CURRENT_TIMESTAMP
//);
//
//ALTER TABLE links
//OWNER TO postgres;
