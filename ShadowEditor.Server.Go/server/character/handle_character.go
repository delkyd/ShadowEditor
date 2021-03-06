package character

import (
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/tengge1/shadoweditor/helper"
	"github.com/tengge1/shadoweditor/server"
	"github.com/tengge1/shadoweditor/server/category"
)

func init() {
	handler := Character{}
	server.Mux.UsingContext().Handle(http.MethodGet, "/api/Character/List", handler.List)
	server.Mux.UsingContext().Handle(http.MethodGet, "/api/Character/Get", handler.Get)
	server.Mux.UsingContext().Handle(http.MethodPost, "/api/Character/Edit", handler.Edit)
	server.Mux.UsingContext().Handle(http.MethodPost, "/api/Character/Save", handler.Save)
	server.Mux.UsingContext().Handle(http.MethodPost, "/api/Character/Delete", handler.Delete)
}

// Character 人物控制器
type Character struct {
}

// List 获取列表
func (Character) List(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	db, err := server.Mongo()
	if err != nil {
		helper.WriteJSON(w, server.Result{
			Code: 300,
			Msg:  err.Error(),
		})
		return
	}

	// 获取所有类别
	filter := bson.M{
		"Type": "Character",
	}
	categories := []category.Model{}
	db.FindMany(server.CategoryCollectionName, filter, &categories)

	docs := bson.A{}

	opts := options.FindOptions{
		Sort: bson.M{
			"_id": -1,
		},
	}

	if server.Config.Authority.Enabled {
		user, _ := server.GetCurrentUser(r)

		if user != nil {
			filter1 := bson.M{
				"UserID": user.ID,
			}

			if user.Name == "Administrator" {
				filter2 := bson.M{
					"UserID": bson.M{
						"$exists": 0,
					},
				}
				filter1 = bson.M{
					"$or": bson.A{
						filter1,
						filter2,
					},
				}
			}
			db.FindMany(server.CharacterCollectionName, filter1, &docs, &opts)
		}
	} else {
		db.FindAll(server.CharacterCollectionName, &docs, &opts)
	}

	list := []Model{}
	for _, i := range docs {
		doc := i.(primitive.D).Map()
		categoryID := ""
		categoryName := ""

		if doc["Category"] != nil {
			for _, category := range categories {
				if category.ID == doc["Category"].(string) {
					categoryID = category.ID
					categoryName = category.Name
					break
				}
			}
		}

		thumbnail, _ := doc["Thumbnail"].(string)

		info := Model{
			ID:           doc["_id"].(primitive.ObjectID).Hex(),
			Name:         doc["Name"].(string),
			CategoryID:   categoryID,
			CategoryName: categoryName,
			TotalPinYin:  helper.PinYinToString(doc["TotalPinYin"]),
			FirstPinYin:  helper.PinYinToString(doc["FirstPinYin"]),
			CreateTime:   doc["CreateTime"].(primitive.DateTime).Time(),
			UpdateTime:   doc["UpdateTime"].(primitive.DateTime).Time(),
			Thumbnail:    thumbnail,
		}

		list = append(list, info)
	}

	helper.WriteJSON(w, server.Result{
		Code: 200,
		Msg:  "Get Successfully!",
		Data: list,
	})
}

// Get 获取
func (Character) Get(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(r.FormValue("ID")))
	if err != nil {
		helper.WriteJSON(w, server.Result{
			Code: 300,
			Msg:  "ID is not allowed.",
		})
	}

	db, err := server.Mongo()
	if err != nil {
		helper.WriteJSON(w, server.Result{
			Code: 300,
			Msg:  err.Error(),
		})
		return
	}

	filter := bson.M{
		"ID": id,
	}
	doc := bson.M{}
	find, _ := db.FindOne(server.CharacterCollectionName, filter, &doc)

	if !find {
		helper.WriteJSON(w, server.Result{
			Code: 300,
			Msg:  "The character is not existed.",
		})
		return
	}

	thumbnail, _ := doc["Thumbnail"].(string)
	obj := Model{
		ID:          doc["_id"].(primitive.ObjectID).Hex(),
		Name:        doc["Name"].(string),
		TotalPinYin: helper.PinYinToString(doc["TotalPinYin"]),
		FirstPinYin: helper.PinYinToString(doc["FirstPinYin"]),
		CreateTime:  doc["CreateTime"].(primitive.DateTime).Time(),
		UpdateTime:  doc["UpdateTime"].(primitive.DateTime).Time(),
		Data:        doc["Data"].(string),
		Thumbnail:   thumbnail,
	}

	helper.WriteJSON(w, server.Result{
		Code: 200,
		Msg:  "Get Successfully!",
		Data: obj,
	})
}

// Edit 编辑
func (Character) Edit(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(r.FormValue("ID")))
	if err != nil {
		helper.WriteJSON(w, server.Result{
			Code: 300,
			Msg:  "ID is not allowed.",
		})
	}

	name := strings.TrimSpace(r.FormValue("Name"))
	if name == "" {
		helper.WriteJSON(w, server.Result{
			Code: 300,
			Msg:  "Name is not allowed to be empty.",
		})
		return
	}

	thumbnail := strings.TrimSpace(r.FormValue("Thumbnail"))
	category := strings.TrimSpace(r.FormValue("Category"))

	// update mongo
	db, err := server.Mongo()
	if err != nil {
		helper.WriteJSON(w, server.Result{
			Code: 300,
			Msg:  err.Error(),
		})
		return
	}

	pinyin := helper.ConvertToPinYin(name)

	filter := bson.M{
		"ID": id,
	}
	set := bson.M{
		"Name":        name,
		"TotalPinYin": pinyin.TotalPinYin,
		"FirstPinYin": pinyin.FirstPinYin,
		"Thumbnail":   thumbnail,
	}
	update := bson.M{
		"$set": set,
	}
	if category == "" {
		update["$unset"] = bson.M{
			"Category": 1,
		}
	} else {
		set["Category"] = category
	}

	db.UpdateOne(server.CharacterCollectionName, filter, update)

	helper.WriteJSON(w, server.Result{
		Code: 200,
		Msg:  "Saved successfully!",
	})
}

// Save 保存
func (Character) Save(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(r.FormValue("ID")))
	if err != nil {
		id = primitive.NewObjectID()
	}

	name := strings.TrimSpace(r.FormValue("Name"))
	if name == "" {
		helper.WriteJSON(w, server.Result{
			Code: 300,
			Msg:  "Name is not allowed to be empty.",
		})
		return
	}

	data := strings.TrimSpace(r.FormValue("Data"))

	db, err := server.Mongo()
	if err != nil {
		helper.WriteJSON(w, server.Result{
			Code: 300,
			Msg:  err.Error(),
		})
		return
	}

	filter := bson.M{
		"ID": id,
	}
	doc := bson.M{}
	find, _ := db.FindOne(server.CharacterCollectionName, filter, &doc)

	now := time.Now()

	if !find {
		pinyin := helper.ConvertToPinYin(name)
		doc = bson.M{
			"ID":           id,
			"Name":         name,
			"CategoryID":   0,
			"CategoryName": "",
			"TotalPinYin":  pinyin.TotalPinYin,
			"FirstPinYin":  pinyin.FirstPinYin,
			"Version":      0,
			"CreateTime":   now,
			"UpdateTime":   now,
			"Data":         data,
			"Thumbnail":    "",
		}
		if server.Config.Authority.Enabled {
			user, err := server.GetCurrentUser(r)

			if err != nil && user != nil {
				doc["UserID"] = user.ID
			}
		}

		db.InsertOne(server.CharacterCollectionName, doc)
	} else {
		update := bson.M{
			"UpdateTime": now,
			"Data":       data,
		}
		db.UpdateOne(server.CharacterCollectionName, filter, update)
	}

	helper.WriteJSON(w, server.Result{
		Code: 200,
		Msg:  "Saved successfully!",
	})
}

// Delete 删除
func (Character) Delete(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id, err := primitive.ObjectIDFromHex(r.FormValue("ID"))
	if err != nil {
		helper.WriteJSON(w, server.Result{
			Code: 300,
			Msg:  "ID is not allowed.",
		})
		return
	}

	db, err := server.Mongo()
	if err != nil {
		helper.WriteJSON(w, server.Result{
			Code: 300,
			Msg:  err.Error(),
		})
		return
	}

	filter := bson.M{
		"ID": id,
	}

	doc := bson.M{}
	find, _ := db.FindOne(server.CharacterCollectionName, filter, &doc)

	if !find {
		helper.WriteJSON(w, server.Result{
			Code: 300,
			Msg:  "The asset is not existed!",
		})
		return
	}

	db.DeleteOne(server.CharacterCollectionName, filter)

	helper.WriteJSON(w, server.Result{
		Code: 200,
		Msg:  "Delete successfully!",
	})
}
