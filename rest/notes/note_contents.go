package notes

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"onysakura.fun/Server/commons/sqlite"
)

type NoteContents struct {
	Id      int64
	Title   string
	Level   int
	Type    NoteContentsType
	Parent  int64
	Content []byte
	Status  NoteStatus
}

type NoteContentsType int

const (
	NoteContentsTypeDir     NoteContentsType = 0
	NoteContentsTypeContent NoteContentsType = 1
	NoteContentsTypeImage   NoteContentsType = 2
)

type NoteStatus int

const (
	NoteStatusDisabled NoteStatus = 0
	NoteStatusPrivate  NoteStatus = 1
	NoteStatusPublic   NoteStatus = 2
)

func CleanContents() {
	var err error
	var stmt *sql.Stmt
	defer func() {
		_ = stmt.Close()
		if err != nil {
			log.Warn("clear note contents fail: ", err)
		}
	}()
	stmt, err = sqlite.DB.Prepare("DELETE FROM note_contents")
	if err != nil {
		return
	}
	_, err = stmt.Exec()
	if err != nil {
		return
	}
	stmt, err = sqlite.DB.Prepare("DELETE FROM sqlite_sequence WHERE name = 'note_contents'")
	if err != nil {
		return
	}
	_, err = stmt.Exec()
}

func GetList() interface{} {
	var err error
	var stmt *sql.Stmt
	defer func() {
		_ = stmt.Close()
		if err != nil {
			log.Warn("get note list fail: ", err)
		}
	}()
	stmt, err = sqlite.DB.Prepare("select id, title, level, type, parent, status from note_contents")
	if err != nil {
		return nil
	}
	rows, err := stmt.Query()
	defer rows.Close()
	if err != nil {
		return nil
	}
	var list []NoteContents
	for rows.Next() {
		var noteContents NoteContents
		err = rows.Scan(&noteContents.Id, &noteContents.Title, &noteContents.Level, &noteContents.Type, &noteContents.Parent, &noteContents.Status)
		if err != nil {
			return nil
		}
		list = append(list, noteContents)
	}
	return list
}

func SaveNoteContents(noteContents NoteContents) int64 {
	var err error
	var stmt *sql.Stmt
	defer func() {
		_ = stmt.Close()
		if err != nil {
			log.Warn("saveToDatabase note contents fail: ", err)
		}
	}()
	stmt, err = sqlite.DB.Prepare("insert into note_contents(title, level, type, parent, content, status) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return 0
	}
	result, err := stmt.Exec(noteContents.Title, noteContents.Level, noteContents.Type, noteContents.Parent, noteContents.Content, noteContents.Status)
	if err != nil {
		return 0
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0
	}
	return id
}

func GetNoteContent(id string) (*[]byte, *NoteStatus) {
	var err error
	var stmt *sql.Stmt
	defer func() {
		_ = stmt.Close()
		if err != nil {
			log.Warn("get note fail: ", err)
		}
	}()
	var rows *sql.Rows
	stmt, err = sqlite.DB.Prepare("select content, status from note_contents where id = ?")
	if err != nil {
		return nil, nil
	}
	rows, err = stmt.Query(id)
	defer rows.Close()
	if err != nil {
		return nil, nil
	}
	for rows.Next() {
		var content []byte
		var status NoteStatus
		err = rows.Scan(&content, &status)
		if err != nil {
			return nil, nil
		}
		return &content, &status
	}
	return nil, nil
}
func GetNoteImage(parent string, title string) *[]byte {
	var err error
	var stmt *sql.Stmt
	defer func() {
		_ = stmt.Close()
		if err != nil {
			log.Warn("get note fail: ", err)
		}
	}()
	var rows *sql.Rows
	stmt, err = sqlite.DB.Prepare("select content from note_contents where parent = (select id from note_contents where parent = ? and title = 'images') and title = ?")
	if err != nil {
		return nil
	}
	rows, err = stmt.Query(parent, title)
	defer rows.Close()
	if err != nil {
		return nil
	}
	for rows.Next() {
		var content []byte
		err = rows.Scan(&content)
		if err != nil {
			return nil
		}
		return &content
	}
	return nil
}
