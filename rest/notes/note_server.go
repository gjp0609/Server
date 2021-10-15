package notes

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"onysakura.fun/Server/commons"
	"onysakura.fun/Server/commons/data"
	"onysakura.fun/Server/commons/logrus"
	"onysakura.fun/Server/rest/user"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"
	"time"
)

var log = logrus.GetLogger()
var ignoredFile = []string{".git", ".github", ".gitignore", ".idea", "list.json", "node_modules", "package-lock.json", "package.json", "prettier.config.js"}

func init() {
}

func GetNote(writer http.ResponseWriter, request *http.Request) {
	var err error
	defer func() {
		if err != nil {
			log.Warning("get note err, ", err)
		}
	}()
	id := request.URL.Query().Get("id")
	parent := request.URL.Query().Get("parent")
	title := request.URL.Query().Get("title")
	var noteContent *[]byte
	if len(id) > 0 {
		var status *NoteStatus
		noteContent, status = GetNoteContent(id)
		if noteContent != nil {
			if *status == NoteStatusPrivate {
				authorization := request.Header.Get("Authorization")
				_, err = user.Auth(authorization)
				if err != nil {
					writer.WriteHeader(401)
					return
				}
			}
		}
	} else {
		noteContent = GetNoteImage(parent, title)
	}
	if noteContent == nil {
		err = errors.New("note not found")
		return
	}
	http.ServeContent(writer, request, "asd.png", time.Time{}, bytes.NewReader(*noteContent))
	_, _ = writer.Write(*noteContent)
}

func List(writer http.ResponseWriter, request *http.Request) {
	var returnMsg = data.NewErrorMsg()
	var err error
	defer func() {
		if err != nil {
			log.Warning(err)
			returnMsg.Msg = returnMsg.Msg + ", Reason: " + err.Error()
		}
		_, _ = writer.Write(returnMsg.ToString())
	}()
	tree := GetList()
	returnMsg = data.Msg{Code: data.MsgOk, Msg: "ok", Data: tree}
}

func AddNote(writer http.ResponseWriter, request *http.Request) {
	var returnMsg = data.NewErrorMsg()
	var err error
	defer func() {
		if err != nil {
			log.Warning(err)
			returnMsg.Msg = returnMsg.Msg + ", Reason: " + err.Error()
		}
		_, _ = writer.Write(returnMsg.ToString())
	}()
	authorization := request.Header.Get("Authorization")
	log.Info("Authorization: ", authorization)
	var username *string
	username, err = user.Auth(authorization)
	if err != nil {
		returnMsg.Msg = "auth fail"
		return
	}
	log.Info(fmt.Sprintf("username: %s", *username))
	returnMsg = data.Msg{Code: data.MsgOk, Msg: "Hello " + *username}
}

func UpdateNotes(writer http.ResponseWriter, request *http.Request) {
	defer func() {
		writer.WriteHeader(200)
	}()
	command := exec.Command("git", "-C", commons.Configs.Notes.Path, "pull")
	output, err := command.Output()
	log.Info("update notes: ", string(output), err)
	saveToDatabase()
	log.Info("save to database done")
}

func ServeNotes(writer http.ResponseWriter, request *http.Request) {
	var err error
	defer func() {
		if err != nil {
			writer.WriteHeader(401)
		}
	}()
	filePath := commons.Configs.Notes.Path + request.URL.Path
	filePath = strings.Replace(filePath, "/notes/", "", 1)
	log.Debug("serve file path: ", filePath)
	if strings.Contains(filePath, ".md") {
		var file *os.File
		file, err = os.Open(filePath)
		defer func(file *os.File) { _ = file.Close() }(file)
		if err != nil {
			return
		}
		scanner := bufio.NewScanner(file)
		log.Info("sc")
		if scanner.Scan() {
			text := scanner.Text()
			if strings.Contains(text, "[comment]: <> (private)") {
				authorization := request.Header.Get("Authorization")
				_, err = user.Auth(authorization)
				if err != nil {
					return
				}
			}
		}
	}
	http.ServeFile(writer, request, filePath)
}

func saveToDatabase() {
	var err error
	defer func() {
		if err != nil {
			log.Warning(err)
		}
	}()
	CleanContents()
	rootPath := commons.Configs.Notes.Path
	getFileList(rootPath, 0, 0)
}

func getFileList(dir string, level int, parentId int64) {
	level++
	root, err := os.Open(dir)
	if err != nil {
		return
	}
	files, err := root.ReadDir(-1)
	if err != nil {
		return
	}
	for index := range files {
		name := files[index].Name()
		if !contains(ignoredFile, name) {
			absPath := path.Join(dir, name)
			if files[index].IsDir() {
				noteContents := NoteContents{
					Title:   name,
					Level:   level,
					Type:    NoteContentsTypeDir,
					Parent:  parentId,
					Content: []byte(""),
					Status:  NoteStatusPrivate,
				}
				id := SaveNoteContents(noteContents)
				if id != 0 {
					getFileList(absPath, level, id)
				}
			} else {
				file, err := os.ReadFile(absPath)
				if err != nil {
					return
				}
				var noteContentType NoteContentsType
				var noteStatus NoteStatus
				var title string
				if strings.Contains(absPath, ".md") {
					title = strings.Replace(name, ".md", "", -1)
					noteContentType = NoteContentsTypeContent
					buffer := bytes.NewBuffer(file)
					private := []byte("[comment]: <> (private)")
					head := buffer.Next(len(private))
					if bytes.Compare(private, head) == 0 {
						noteStatus = NoteStatusPrivate
					} else {
						noteStatus = NoteStatusPublic
					}
				} else {
					noteContentType = NoteContentsTypeImage
					title = name
					noteStatus = NoteStatusPublic
				}
				noteContents := NoteContents{
					Title:   title,
					Level:   level,
					Type:    noteContentType,
					Parent:  parentId,
					Content: file,
					Status:  noteStatus,
				}
				_ = SaveNoteContents(noteContents)
			}
		}
	}
}

func contains(s []string, searchterm string) bool {
	i := sort.SearchStrings(s, searchterm)
	return i < len(s) && s[i] == searchterm
}
