package uchatlib

/*
  关群调用
*/
func SetChatRoomOver(chatRoomSerialNo, comment string, client *UchatClient) error {
	ctx := make(map[string]string, 0)
	ctx["vcChatRoomSerialNo"] = chatRoomSerialNo
	ctx["vcComment"] = comment
	return client.ChatRoomOver(ctx)
}

/*
  群内时时消息开启设置
*/
func SetChatRoomOpenGetMessage(chatRoomSerialNo string, client *UchatClient) error {
	ctx := make(map[string]string, 0)
	ctx["vcChatRoomSerialNo"] = chatRoomSerialNo
	return client.ChatRoomOpenGetMessages(ctx)
}

/*
  群内时时消息关闭设置
*/
func SetChatRoomCloseGetMessage(chatRoomSerialNo string, client *UchatClient) error {
	ctx := make(map[string]string, 0)
	ctx["vcChatRoomSerialNo"] = chatRoomSerialNo
	return client.ChatRoomCloseGetMessages(ctx)
}

func ApplyRobotAddUser(robotSerialNo, userWexinId string, client *UchatClient) error {
	ctx := make(map[string]string, 0)
	ctx["vcRobotSerialNo"] = robotSerialNo
	ctx["vcUserWeixinId"] = userWexinId
	return client.RobotAddUser(ctx)
}
