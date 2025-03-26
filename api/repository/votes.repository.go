package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/kallel-anobom/event_voting_go/api/model"
	"github.com/kallel-anobom/event_voting_go/api/services/database"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type VotesRepository interface {
	AddVote(vote model.Vote) error
	GetAllVotes() ([]model.Vote, error)
}

type votesRepository struct {
	mongoService *database.MongoService
	mongoClient  *mongo.Client
	dbName       string
}

type MongoRow struct {
	ParticipantID int       `bson:"participant_id"`
	Date          time.Time `bson:"date"`
}

func NewVotesRepository(ms *database.MongoService) VotesRepository {
	return &votesRepository{
		mongoService: ms,
		dbName:       "event_voting_db",
	}
}

func (vt *votesRepository) AddVote(vote model.Vote) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("Attempting to insert vote: %+v", vote) // Log do voto sendo inserido

	collection := vt.mongoClient.Database(vt.dbName).Collection("votes")
	vote.Date = time.Now()

	result, err := collection.InsertOne(ctx, vote)
	if err != nil {
		log.Printf("Error inserting vote: %v", err) // Log detalhado do erro
		return fmt.Errorf("failed to insert vote: %v", err)
	}

	log.Printf("Vote inserted successfully! InsertedID: %v", result.InsertedID)
	return nil
	/*
		Essa função só precisa inserir o voto no mongo
	*/

	/*
		vote.Date = time.Now()
		vt.mongoService.Put(vote)
	*/

	// ctx := context.Background()

	// data, err := json.Marshal(votes)
	// if err != nil {
	// 	return errors.New("failed to marshal votes")
	// }

	// key := "votes:" + votes.ID
	// err = vt.mongoService.Client.Set(ctx, key, data, 168*time.Hour).Err()
	// if err != nil {
	// 	return errors.New("failed to save votes in Redis")
	// }
}

func (r *votesRepository) GetAllVotes() ([]model.Vote, error) {
	if r.mongoClient == nil {
		return nil, errors.New("mongo client is not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := r.mongoClient.Database(r.dbName).Collection("votes")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.New("failed to find votes in MongoDB")
	}
	defer cursor.Close(ctx)

	var votes []model.Vote
	if err = cursor.All(ctx, &votes); err != nil {
		return nil, errors.New("failed to decode votes from MongoDB")
	}

	/*
		Essa função busca todos os votos no MONGO e retorna eles para o USECASE

		o USE CASE será o responsável por calcular e processar a média
	*/
	return votes, nil
}
