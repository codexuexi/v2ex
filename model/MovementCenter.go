package model

import (
	"fmt"
	"github.com/123456/c_code/mc"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"time"
)

//Notice
//Notification
type MovementCenter struct {
	ID primitive.ObjectID `json:"_id" bson:"_id"`
	//发起人
	MID MIDTYPE `json:"mid" bson:"mid"`
	//接收人
	M2ID        MIDTYPE     `json:"m2id" bson:"m2id"`
	Hash        string      `json:"hash" bson:"hash"`
	Read        bool        `json:"read" bson:"read"`
	Type        int         `json:"type" bson:"type"`
	V           interface{} `json:"v" bson:"v"`
	ReleaseTime time.Time   `json:"release_time" bson:"release_time"`
	//用来做查询的临时字段
	WID primitive.ObjectID `json:"-" bson:"-"`
}

func (t MovementCenter) Table() string {
	return "movement_center"
}
func (t MovementCenter) IndexList() []mongo.IndexModel {

	return []mongo.IndexModel{
		{
			Keys: bsonx.Doc{
				{"mid", bsonx.Int32(-1)},
				{"m2id", bsonx.Int32(-1)},
				{"hash", bsonx.Int32(-1)},
			},
		},
	}
}

const (
	//发布文章
	MovementArticleSend = 1

	//发布文章评论
	MovementArticleCommentRoot = 2

	//文章点赞
	MovementArticleCommentGood = 3

	//发布提问
	MovementQuestionSend = 4

	//提问回答
	MovementQuestionCommentRoot = 5

	//回答点赞
	MovementQuestionCommentGood = 6

	//文章 或者 问题收藏
	MovementCollect = 7
)

func Movement(mid, m2id MIDTYPE) MovementCenter {
	return MovementCenter{
		ID:          primitive.NewObjectID(),
		MID:         mid,
		M2ID:        m2id,
		ReleaseTime: time.Now(),
	}
}

func (t MovementCenter) verify() bool {
	if t.MID == t.M2ID {
		return false
	}
	return true
}

type articleSend struct {
	DID DIDTYPE `json:"did" bson:"did"`
}

//发布文章
func (t MovementCenter) AddArticleSend(index DataIndex) {
	if !t.verify() {
		return
	}
	t.Type = MovementArticleSend
	t.V = articleSend{
		DID: index.DID,
	}
	mc.Table(t.Table()).Insert(t)
}

//获取最后活跃的会员

func (t MovementCenter) GetHuoYueMID(nums int64) (m []gin.H) {
	pid := primitive.ObjectID{}
	mids := []MIDTYPE{}
	for {
		where := bson.M{}
		if pid.Hex() != mc.Empty {
			where["_id"] = bson.M{"$lt": pid}
		}
		tmp := []MovementCenter{}
		mc.Table(t.Table()).Where(where).Order(bson.M{"_id": -1}).Limit(nums + 5).Find(&tmp)
		if len(tmp) == 0 {
			break
		}
		pid = tmp[len(tmp)-1].ID
		for _, v := range tmp {

			if _inmids(v.MID, mids) {
				continue
			}
			mids = append(mids, v.MID)
			if len(mids) >= int(nums) {
				break
			}
		}
		if len(mids) >= int(nums) {
			break
		}

	}
	members := []Member{}
	mc.Table(Member{}.Table()).Where(bson.M{"mid": bson.M{"$in": mids}}).Find(&members)
	for _, v := range members {
		m = append(m, gin.H{
			"t":      v.UserName,
			"u":      UrlMember(v),
			"avatar": Avatar(v.Avatar),
			"mid":    v.MID,
		})
	}
	return
}

func _inmids(array_val MIDTYPE, array []MIDTYPE) bool {
	for _, v := range array {
		if v == array_val {
			return true
		}
	}
	return false
}

type questionSend struct {
	DID DIDTYPE `json:"did" bson:"did"`
}

//发布提问
func (t MovementCenter) AddQuestionSend(index DataIndex) {
	if !t.verify() {
		return
	}
	t.Type = MovementQuestionSend
	t.V = questionSend{
		DID: index.DID,
	}
	mc.Table(t.Table()).Insert(t)
}

//提交评论
type articleCommentRoot struct {
	RID primitive.ObjectID `json:"rid" bson:"rid"`
}

func (t MovementCenter) AddArticleCommentRoot(index CommentRoot) {
	if !t.verify() {
		return
	}
	t.Type = MovementArticleCommentRoot
	t.V = articleCommentRoot{
		RID: index.ID,
	}
	mc.Table(t.Table()).Insert(t)
}

//提交回答
type questionCommentRoot struct {
	RID primitive.ObjectID `json:"rid" bson:"rid"`
}

func (t MovementCenter) AddQuestionCommentRoot(index CommentQuestionRoot) {
	if !t.verify() {
		return
	}
	t.Type = MovementQuestionCommentRoot
	t.V = questionCommentRoot{
		RID: index.ID,
	}
	mc.Table(t.Table()).Insert(t)
}

//文章点赞
type articleCommentZan struct {
	RID primitive.ObjectID `json:"rid" bson:"rid"`
}

func (t MovementCenter) AddArticleCommentZan(index CommentRoot) {
	if !t.verify() {
		return
	}
	t.Type = MovementArticleCommentGood
	t.V = articleCommentZan{
		RID: index.ID,
	}
	where := bson.M{"mid": t.MID, "m2id": t.M2ID, "hash": index.ID.Hex()}
	//查询是否存在
	_t := MovementCenter{}
	mc.Table(t.Table()).Where(where).FindOne(&_t)
	if _t.ID.Hex() != mc.Empty {
		return
	}
	t.Hash = index.ID.Hex()
	mc.Table(t.Table()).Insert(t)
}

//回答点赞
type questionAnswerZan struct {
	RID primitive.ObjectID `json:"rid" bson:"rid"`
}

func (t MovementCenter) AddQuestionAnswerZan(index CommentQuestionRoot) {
	if !t.verify() {
		return
	}
	t.Type = MovementQuestionCommentGood
	t.V = questionAnswerZan{
		RID: index.ID,
	}
	where := bson.M{"mid": t.MID, "m2id": t.M2ID, "hash": index.ID.Hex()}
	//查询是否存在
	_t := MovementCenter{}
	mc.Table(t.Table()).Where(where).FindOne(&_t)
	if _t.ID.Hex() != mc.Empty {
		return
	}
	t.Hash = index.ID.Hex()
	mc.Table(t.Table()).Insert(t)
}

//收藏
type collect struct {
	DID DIDTYPE
}

func (t MovementCenter) AddCollect(index DataIndex) {
	if !t.verify() {
		return
	}
	t.Type = MovementCollect
	t.V = collect{
		DID: index.DID,
	}
	where := bson.M{"mid": t.MID, "m2id": t.M2ID, "hash": index.ID.Hex()}
	//查询是否存在
	_t := MovementCenter{}
	mc.Table(t.Table()).Where(where).FindOne(&_t)
	if _t.ID.Hex() != mc.Empty {
		return
	}
	t.Hash = index.ID.Hex()
	mc.Table(t.Table()).Insert(t)
}

//删除收藏信息
func (t MovementCenter) DelCollect(index DataIndex) {
	if !t.verify() {
		return
	}
	where := bson.M{"mid": t.MID, "m2id": t.M2ID, "hash": index.ID.Hex()}
	fmt.Println(where)
	mc.Table(t.Table()).Where(where).DelOne()
}

//删除文章点赞
func (t MovementCenter) DelArticleCommentZan(index CommentRoot) {
	if !t.verify() {
		return
	}
	where := bson.M{"mid": t.MID, "m2id": t.M2ID, "hash": index.ID.Hex()}
	fmt.Println(where)
	mc.Table(t.Table()).Where(where).DelOne()
}

//删除回答点赞
func (t MovementCenter) DelQuestionCommentZan(index CommentQuestionRoot) {
	if !t.verify() {
		return
	}
	where := bson.M{"mid": t.MID, "m2id": t.M2ID, "hash": index.ID.Hex()}
	fmt.Println(where)
	mc.Table(t.Table()).Where(where).DelOne()
}
