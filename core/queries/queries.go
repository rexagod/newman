package queries

const (
	ADD                        = "add"
	ADDDELETEDMESSAGE          = "addDeletedMessage"
	CREATEDELETEDMESSAGESTABLE = "createDeletedMessages"
	CREATENOTESTABLE           = "createNotes"
	REMOVE                     = "remove"
	SHOW                       = "show"
	SHOWDELETEDMESSAGES        = "showDeletedMessages"
)

var Q = map[string]string{
	"add":                   "insert into notes (note) values ($1)",
	"addDeletedMessage":     "insert into deletedMessages (author, timestamp, message) values ($1, $2, $3)",
	"createDeletedMessages": "create table deletedMessages(id int not null identity(1,1),author varchar(50) not null,timestamp datetime not null,message varchar(255) not null)",
	"createNotes":           "create table notes(id int not null identity(1,1),note varchar(255) not null)",
	"remove":                "delete from notes where id = $1",
	"show":                  "select * from notes",
	"showDeletedMessages":   "select * from deletedMessages",
}
