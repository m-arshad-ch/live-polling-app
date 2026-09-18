package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"live-polling-app/models"
)

type PollController struct {
	db        *mongo.Database
	redis     *redis.Client
	jwtSecret string
}

func NewPollController(db *mongo.Database, redisClient *redis.Client, jwtSecret string) *PollController {
	return &PollController{db: db, redis: redisClient, jwtSecret: jwtSecret}
}

type createPollRequest struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
}

type voteRequest struct {
	OptionID string `json:"optionId"`
	VoterID  string `json:"voterId"`
}

type resultItem struct {
	OptionID string `json:"optionId"`
	Text     string `json:"text"`
	Votes    int64  `json:"votes"`
}

func (p *PollController) CreatePoll(c *gin.Context) {
	var req createPollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	question := strings.TrimSpace(req.Question)
	if len(question) < 3 || len(question) > 200 {
		c.JSON(400, gin.H{"error": "question must be between 3 and 200 characters"})
		return
	}

	if len(req.Options) < 2 || len(req.Options) > 6 {
		c.JSON(400, gin.H{"error": "poll must have 2 to 6 options"})
		return
	}

	options := make([]models.PollOption, 0, len(req.Options))
	seen := map[string]bool{}
	for i, text := range req.Options {
		text = strings.TrimSpace(text)
		if len(text) < 1 || len(text) > 100 {
			c.JSON(400, gin.H{"error": "each option must be 1 to 100 characters"})
			return
		}
		if seen[strings.ToLower(text)] {
			c.JSON(400, gin.H{"error": "options must be different"})
			return
		}
		seen[strings.ToLower(text)] = true
		options = append(options, models.PollOption{
			ID:   fmt.Sprintf("option-%d", i+1),
			Text: text,
		})
	}

	userID := c.MustGet("userID").(primitive.ObjectID)
	poll := models.Poll{
		Question:  question,
		Options:   options,
		CreatedBy: userID,
		CreatedAt: time.Now().Unix(),
		IsClosed:  false,
	}

	result, err := p.db.Collection("polls").InsertOne(context.Background(), poll)
	if err != nil {
		c.JSON(500, gin.H{"error": "could not create poll"})
		return
	}
	poll.ID = result.InsertedID.(primitive.ObjectID)

	c.JSON(201, poll)
}

func (p *PollController) GetPoll(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid poll id"})
		return
	}

	var poll models.Poll
	err = p.db.Collection("polls").FindOne(context.Background(), bson.M{"_id": id}).Decode(&poll)
	if err != nil {
		c.JSON(404, gin.H{"error": "poll not found"})
		return
	}

	c.JSON(200, poll)
}

func (p *PollController) MyPolls(c *gin.Context) {
	userID := c.MustGet("userID").(primitive.ObjectID)
	cursor, err := p.db.Collection("polls").Find(
		context.Background(),
		bson.M{"createdBy": userID},
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "could not load polls"})
		return
	}
	defer cursor.Close(context.Background())

	var polls []models.Poll
	if err := cursor.All(context.Background(), &polls); err != nil {
		c.JSON(500, gin.H{"error": "could not load polls"})
		return
	}
	if polls == nil {
		polls = []models.Poll{}
	}
	c.JSON(200, polls)
}

func (p *PollController) Vote(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid poll id"})
		return
	}

	var req voteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}
	req.VoterID = strings.TrimSpace(req.VoterID)
	req.OptionID = strings.TrimSpace(req.OptionID)

	if len(req.VoterID) < 10 || len(req.VoterID) > 100 {
		c.JSON(400, gin.H{"error": "invalid voter id"})
		return
	}

	var poll models.Poll
	if err := p.db.Collection("polls").FindOne(context.Background(), bson.M{"_id": id}).Decode(&poll); err != nil {
		c.JSON(404, gin.H{"error": "poll not found"})
		return
	}
	if poll.IsClosed {
		c.JSON(400, gin.H{"error": "this poll is closed"})
		return
	}

	validOption := false
	for _, option := range poll.Options {
		if option.ID == req.OptionID {
			validOption = true
			break
		}
	}
	if !validOption {
		c.JSON(400, gin.H{"error": "invalid option"})
		return
	}

	vote := models.Vote{
		PollID:    id,
		OptionID:  req.OptionID,
		VoterID:   req.VoterID,
		CreatedAt: time.Now().Unix(),
	}

	_, err = p.db.Collection("votes").InsertOne(context.Background(), vote)
	if err != nil {
		// Duplicate compound index means this voter already voted.
		c.JSON(409, gin.H{"error": "you have already voted in this poll"})
		return
	}

	ctx := context.Background()
	key := "poll:" + id.Hex() + ":counts"
	count, err := p.redis.HIncrBy(ctx, key, req.OptionID, 1).Result()
	if err != nil {
		c.JSON(500, gin.H{"error": "vote saved but live update failed"})
		return
	}

	results, err := p.buildResults(ctx, poll)
	if err != nil {
		c.JSON(500, gin.H{"error": "vote saved but results could not be prepared"})
		return
	}

	payload, _ := json.Marshal(gin.H{
		"optionID": req.OptionID,
		"newCount": count,
		"results":  results,
	})

	channel := "poll:" + id.Hex() + ":live"
	if err := p.redis.Publish(ctx, channel, payload).Err(); err != nil {
		c.JSON(500, gin.H{"error": "vote saved but live event failed"})
		return
	}

	c.JSON(200, gin.H{"message": "vote recorded", "results": results})
}

func (p *PollController) Results(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid poll id"})
		return
	}

	var poll models.Poll
	if err := p.db.Collection("polls").FindOne(context.Background(), bson.M{"_id": id}).Decode(&poll); err != nil {
		c.JSON(404, gin.H{"error": "poll not found"})
		return
	}

	results, err := p.buildResults(context.Background(), poll)
	if err != nil {
		c.JSON(500, gin.H{"error": "could not load results"})
		return
	}

	c.JSON(200, gin.H{"results": results})
}

func (p *PollController) buildResults(ctx context.Context, poll models.Poll) ([]resultItem, error) {
	key := "poll:" + poll.ID.Hex() + ":counts"
	fields := make([]string, 0, len(poll.Options))
	for _, option := range poll.Options {
		fields = append(fields, option.ID)
	}

	values, err := p.redis.HMGet(ctx, key, fields...).Result()
	if err != nil {
		return nil, err
	}

	needSeed := true
	for _, value := range values {
		if value != nil {
			needSeed = false
			break
		}
	}

	if needSeed {
		pipeline := p.redis.Pipeline()
		for _, option := range poll.Options {
			pipeline.HSet(ctx, key, option.ID, 0)
		}

		cursor, err := p.db.Collection("votes").Find(ctx, bson.M{"pollId": poll.ID})
		if err != nil {
			return nil, err
		}
		var votes []models.Vote
		if err := cursor.All(ctx, &votes); err != nil {
			return nil, err
		}

		counts := map[string]int64{}
		for _, vote := range votes {
			counts[vote.OptionID]++
		}
		for optionID, count := range counts {
			pipeline.HSet(ctx, key, optionID, count)
		}
		if _, err := pipeline.Exec(ctx); err != nil {
			return nil, err
		}

		values = make([]interface{}, len(fields))
		for i, option := range poll.Options {
			values[i] = counts[option.ID]
		}
	}

	results := make([]resultItem, 0, len(poll.Options))
	for i, option := range poll.Options {
		votes := int64(0)
		switch value := values[i].(type) {
		case string:
			votes, _ = strconv.ParseInt(value, 10, 64)
		case int64:
			votes = value
		case nil:
			votes = 0
		}
		results = append(results, resultItem{
			OptionID: option.ID,
			Text:     option.Text,
			Votes:    votes,
		})
	}
	return results, nil
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (p *PollController) Live(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid poll id"})
		return
	}

	var poll models.Poll
	if err := p.db.Collection("polls").FindOne(context.Background(), bson.M{"_id": id}).Decode(&poll); err != nil {
		c.JSON(404, gin.H{"error": "poll not found"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	pubsub := p.redis.Subscribe(context.Background(), "poll:"+id.Hex()+":live")
	defer pubsub.Close()

	for {
		msg, err := pubsub.ReceiveMessage(context.Background())
		if err != nil {
			return
		}
		if err := conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
			return
		}
	}
}
