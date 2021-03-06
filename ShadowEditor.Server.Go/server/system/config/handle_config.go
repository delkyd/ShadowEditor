package config

import (
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/tengge1/shadoweditor/helper"
	"github.com/tengge1/shadoweditor/server"
	"github.com/tengge1/shadoweditor/server/system/model"
)

func init() {
	config := Config{}
	server.Mux.UsingContext().Handle(http.MethodGet, "/api/Config/Get", config.Get)
}

// Config 配置控制器
type Config struct {
}

// Get 获取配置信息
func (Config) Get(w http.ResponseWriter, r *http.Request) {
	db, err := server.Mongo()
	if err != nil {
		helper.Write(w, err.Error())
		return
	}

	config := model.Config{}
	find, err := db.FindOne(server.ConfigCollectionName, bson.M{}, &config)
	if err != nil {
		helper.WriteJSON(w, server.Result{
			Code: 300,
			Msg:  err.Error(),
		})
		return
	}

	if !find {
		doc1 := bson.M{
			"ID":                  primitive.NewObjectID().Hex(),
			"Initialized":         false,
			"DefaultRegisterRole": "",
		}
		db.InsertOne(server.ConfigCollectionName, doc1)
		db.FindOne(server.ConfigCollectionName, bson.M{}, &config)
	}

	result := Result{
		ID:                   config.ID,
		EnableAuthority:      server.Config.Authority.Enabled,
		Initialized:          config.Initialized,
		DefaultRegisterRole:  config.DefaultRegisterRole,
		IsLogin:              false,
		Username:             "",
		Name:                 "",
		RoleID:               "",
		RoleName:             "",
		DeptID:               "",
		DeptName:             "",
		OperatingAuthorities: []string{},
		EnableRemoteEdit:     server.Config.Remote.Enabled,
		WebSocketServerPort:  server.Config.Remote.WebSocketPort,
	}

	user, err := server.GetCurrentUser(r)
	if err != nil {
		helper.Write(w, err.Error())
		return
	}

	if user != nil {
		result.IsLogin = true
		result.Username = user.Username
		result.Name = user.Name
		result.RoleID = user.RoleID
		result.RoleName = user.RoleName
		result.DeptID = user.DeptID
		result.DeptName = user.DeptName
	}

	helper.WriteJSON(w, server.Result{
		Code: 200,
		Msg:  "Get Successfully!",
		Data: result,
	})
}

// Result config to front end
type Result struct {
	ID                   string
	EnableAuthority      bool
	Initialized          bool
	DefaultRegisterRole  string
	IsLogin              bool
	Username             string
	Name                 string
	RoleID               string
	RoleName             string
	DeptID               string
	DeptName             string
	OperatingAuthorities []string
	EnableRemoteEdit     bool
	WebSocketServerPort  int
}
