package auth

import (
	"bpm/core/queue"

	"github.com/streadway/amqp"
)

type NewAuthCreated struct {
	AuthID     int64  `json:"auth_id"`
	AuthType   int    `json:"auth_type"`
	Identifier string `json:"identifier"`
	Credential string `json:"credential"`
	Name       string `json:"name"`
	Email      string `json:"email"`
}

type NewProfileCreated struct {
	AuthID int64 `json:"auth_id"`
	UserID int64 `json:"user_id"`
}

func Subscribe(conn *queue.Conn) {
	conn.StartConsumer("UpdateAuthUserID", "NewProfileCreated", UpdateAuthUserID)
}

func UpdateAuthUserID(d amqp.Delivery) bool {
	// if d.Body == nil {
	// 	return false
	// }
	// var NewProfileCreated NewProfileCreated
	// err := json.Unmarshal(d.Body, &NewProfileCreated)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return false
	// }
	// var authInfo UserAuth
	// authInfo.UserID = NewProfileCreated.UserID
	// authInfo.ID = NewProfileCreated.AuthID
	// db := database.InitMySQL()
	// repo := NewAuthRepository(db)
	// err = repo.UpdateUserID(authInfo)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return false
	// }
	return true
}
