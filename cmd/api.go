// Copyright © 2017 edwin <edwin.lzh@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"log"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/lvzhihao/goutils"
	"github.com/lvzhihao/uchat/models"
	"github.com/lvzhihao/uchat/uchat"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		defer Logger.Sync()
		app := goutils.NewEcho()
		//app.Logger.SetLevel(log.INFO)
		client := uchat.NewClient(viper.GetString("merchant_no"), viper.GetString("merchant_secret"))
		/*
			tool, err := NewTool(fmt.Sprintf("amqp://%s:%s@%s/%s", viper.GetString("rabbitmq_user"), viper.GetString("rabbitmq_passwd"), viper.GetString("rabbitmq_host"), viper.GetString("rabbitmq_vhost")))
			if err != nil {
				Logger.Error("RabbitMQ Connect Error", zap.Error(err))
			}
		*/
		db, err := gorm.Open("mysql", viper.GetString("mysql_dns"))
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
			return viper.GetString("table_prefix") + "_" + defaultTableName
		}
		app.POST("/applycode", func(ctx echo.Context) error {
			myId := ctx.FormValue("my_id")
			if myId == "" {
				return ctx.HTML(http.StatusOK, "Error my_id")
			}
			subId := ctx.QueryParam("sub_id")
			var myRobots []models.MyRobot
			err := db.Where("my_id = ?", myId).Find(&myRobots).Error
			if err != nil {
				return ctx.HTML(http.StatusOK, "Error"+err.Error())
			}
			var robotIds []string
			for _, myRobot := range myRobots {
				robotIds = append(robotIds, myRobot.RobotSerialNo)
			}
			if len(robotIds) == 0 {
				return ctx.HTML(http.StatusOK, "Error robotIds")
			}
			var robots []models.Robot
			err = db.Where("serial_no in (?)", robotIds).Where("chat_room_count < ?", 30).Find(&robots).Error
			if err != nil {
				return ctx.HTML(http.StatusOK, "Error"+err.Error())
			}
			if len(robots) == 0 {
				return ctx.HTML(http.StatusOK, "Error robots")
			}
			//todo check subid chatroom limit
			params := map[string]string{
				"vcRobotSerialNo":    robots[0].SerialNo, //todo
				"nType":              "10",
				"vcChatRoomSerialNo": "",
				"nCodeCount":         "1",
				"nAddMinute":         "30", //暂定30分钟过期
				"nIsNotify":          "0",
				"vcNotifyContent":    "",
			}
			datas, err := client.ApplyCodeList(params)
			if err != nil {
				return ctx.HTML(http.StatusOK, "Error"+err.Error())
			}
			codeData := datas[0]
			loc, _ := time.LoadLocation("Asia/Shanghai")
			applyCode := &models.RobotApplyCode{}
			applyCode.RobotSerialNo = codeData["vcRobotSerialNo"]
			applyCode.RobotNickName = codeData["vcNickName"]
			applyCode.Type = "10"
			applyCode.ChatRoomSerialNo = codeData["vcChatRoomSerialNo"]
			applyCode.ExpireTime, _ = time.ParseInLocation("2006/1/2 15:04:05", codeData["dtEndDate"], loc)
			applyCode.CodeSerialNo = codeData["vcSerialNo"]
			applyCode.CodeValue = codeData["vcCode"]
			applyCode.CodeImages = codeData["vcCodeImages"]
			applyCode.CodeTime, _ = time.ParseInLocation("2006/1/2 15:04:05", codeData["dtCreateDate"], loc)
			applyCode.MyId = myId
			applyCode.SubId = subId
			applyCode.Used = false
			err = db.Create(applyCode).Error
			if err != nil {
				return ctx.HTML(http.StatusOK, "Error"+err.Error())
			}
			return ctx.JSON(http.StatusOK, applyCode)
		})
		goutils.EchoStartWithGracefulShutdown(app, ":8079")

	},
}

func init() {
	RootCmd.AddCommand(apiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// apiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// apiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}