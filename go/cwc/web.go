/*
Copyright (C) 2019 Graeme Sutherland, Nodestone Limited


This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package cwc

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

const WebPort = ":12345"

var cfg *Config

func WebServer(ctx context.Context, config *Config) {

	cfg = config
	glog.Infof("Got config with wpm of %d\n", cfg.KeyerSpeed)
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Static("/assets", "./assets")
	router.LoadHTMLFiles("index.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	router.GET("/ws", func(c *gin.Context) {
		wshandler(c.Writer, c.Request)
	})
	router.RunTLS(WebPort, "./certs/server.pem", "./certs/server.key")
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wshandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: %+v", err)
		return
	}

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))
		m := string(msg)

		if m == "fromC:status:connected" {
			var res string
			res = fmt.Sprintf("toC:wpm:%d", cfg.KeyerSpeed)
			glog.Infof("res is: %s", res)
			conn.WriteMessage(t, []byte(res))
		}

		s := strings.Split(m, ":")
		direction, name, value := s[0], s[1], s[2]

		if direction == "toC" {
			return
		}

		if name == "wpm" {
			cfg.KeyerSpeed, err = strconv.Atoi(value)
		}

		if false {
			conn.WriteMessage(t, msg)
		}
	}
}
