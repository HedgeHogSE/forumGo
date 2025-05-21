package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/HedgeHogSE/forum/backend/forum/internal/models"
	"github.com/gorilla/websocket"
)

/*type Message struct {
	Content string `json:"content"`
	Time    string `json:"time"`
}*/

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	messages      = make([]models.CommentWithUsername, 0)
	messagesMutex sync.Mutex
	clients       = make(map[*websocket.Conn]bool)
	clientsMutex  sync.Mutex
)

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	topicID := r.URL.Query().Get("topic")
	num, err := strconv.Atoi(topicID)
	if err != nil {
		log.Println("Ошибка при преобразовании topicID в int:", err)
		return
	}
	messages, err = models.GetCommentsByTopicID(num)
	time := time.Now()
	for _, m := range messages {
		if int(time.Sub(m.CreatedAt).Hours()/24) >= 14 {
			err := models.DeleteCommentByID(m.ID)
			if err != nil {
				log.Println("Ошибка при удалении комментария:", err)
			}
		}
	}

	if err != nil {
		log.Println("Ошибка при получении комментариев:", err)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка при обновлении соединения:", err)
		return
	}
	defer ws.Close()

	clientsMutex.Lock()
	clients[ws] = true
	clientsMutex.Unlock()

	messagesMutex.Lock()
	history, _ := json.Marshal(messages)
	messagesMutex.Unlock()
	ws.WriteMessage(websocket.TextMessage, history)

	log.Println("Новое соединение установлено")

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println(msg)
			log.Println("Ошибка чтения сообщения:", err)
			break
		}

		type IncomingMessage struct {
			Content  string `json:"content"`
			TopicId  int    `json:"topic_id"`
			AuthorId int    `json:"author_id"`
			Username string `json:"username"`
		}

		var newMessage IncomingMessage
		if err := json.Unmarshal(msg, &newMessage); err != nil {
			log.Println("Ошибка разбора сообщения:", err)
			continue
		}

		comment := &models.Comment{
			Content:  newMessage.Content,
			TopicId:  newMessage.TopicId,
			AuthorId: newMessage.AuthorId,
		}

		messagesMutex.Lock()
		_, err = models.AddComment(comment)
		if err != nil {
			log.Println(err)
			return
		}
		messagesMutex.Unlock()

		clientsMutex.Lock()
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("Ошибка отправки:", err)
				client.Close()
				delete(clients, client)
			}
		}
		clientsMutex.Unlock()
	}

	clientsMutex.Lock()
	delete(clients, ws)
	clientsMutex.Unlock()
}
