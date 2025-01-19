package repository

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/iotassss/puzzdra-monster-rating-v2/internal/entity"
)

type MonsterRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewMonsterRepository() *MonsterRepository {
	client := newDynamoDBClient()

	tableNameMonsters := os.Getenv("DYNAMODB_TABLE_NAME")

	return &MonsterRepository{
		client:    client,
		tableName: tableNameMonsters,
	}
}

func newDynamoDBClient() *dynamodb.Client {
	var cfg aws.Config
	var err error

	// 環境変数でローカル実行かどうかを判定
	if os.Getenv("AWS_SAM_LOCAL") == "true" {
		// ローカルDynamoDB Local用の設定
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
				if service == dynamodb.ServiceID {
					return aws.Endpoint{
						URL: "http://host.docker.internal:8000", // Dockerで起動したDynamoDB Localを指す
						// URL:           "http://localhost:8000", // Dockerで起動したDynamoDB Localを指す
						SigningRegion: "ap-northeast-1",
					}, nil
				}
				return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
			})),
		)
		log.Println("Using DynamoDB Local")
	} else {
		// クラウドDynamoDB用の設定
		cfg, err = config.LoadDefaultConfig(context.TODO())
		log.Println("Using DynamoDB in AWS Cloud")
	}

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	return dynamodb.NewFromConfig(cfg)
}

func (repo *MonsterRepository) FindByNo(ctx context.Context, no int) (entity.Monster, error) {
	result, err := repo.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &repo.tableName,
		Key: map[string]types.AttributeValue{
			"No": &types.AttributeValueMemberS{Value: fmt.Sprintf("%d", no)},
		},
	})
	if err != nil {
		return entity.Monster{}, err
	}

	if len(result.Item) == 0 {
		return entity.Monster{}, nil
	}

	game8ScoresLen := len(result.Item["Game8Scores"].(*types.AttributeValueMemberL).Value)
	game8Scores := make([]entity.Game8MonsterScore, 0, game8ScoresLen)

	dbMonster := result.Item

	noStr := dbMonster["No"].(*types.AttributeValueMemberS).Value
	no, err = strconv.Atoi(noStr)
	if err != nil {
		return entity.Monster{}, err
	}
	originMonsterNoStr := dbMonster["OriginMonsterNo"].(*types.AttributeValueMemberS).Value
	originMonsterNo, err := strconv.Atoi(originMonsterNoStr)
	if err != nil {
		return entity.Monster{}, err
	}

	game8URLStr := dbMonster["Game8URL"].(*types.AttributeValueMemberS).Value
	game8URL, err := url.Parse(game8URLStr)
	if err != nil {
		return entity.Monster{}, err
	}

	for _, scores := range dbMonster["Game8Scores"].(*types.AttributeValueMemberL).Value {
		score := scores.(*types.AttributeValueMemberM).Value
		game8Scores = append(game8Scores, entity.Game8MonsterScore{
			Name:        score["Name"].(*types.AttributeValueMemberS).Value,
			LeaderPoint: score["LeaderPoint"].(*types.AttributeValueMemberS).Value,
			SubPoint:    score["SubPoint"].(*types.AttributeValueMemberS).Value,
			AssistPoint: score["AssistPoint"].(*types.AttributeValueMemberS).Value,
		})
	}

	monster := entity.Monster{
		No:              no,
		Name:            dbMonster["Name"].(*types.AttributeValueMemberS).Value,
		OriginMonsterNo: originMonsterNo,
		Game8URL:        game8URL,
		Game8Scores:     game8Scores,
	}

	return monster, nil
}

func (repo *MonsterRepository) Save(ctx context.Context, monster entity.Monster) error {
	game8Scores := make([]types.AttributeValue, 0, len(monster.Game8Scores))

	for _, score := range monster.Game8Scores {
		game8Scores = append(game8Scores, &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"Name":        &types.AttributeValueMemberS{Value: score.Name},
				"LeaderPoint": &types.AttributeValueMemberS{Value: score.LeaderPoint},
				"SubPoint":    &types.AttributeValueMemberS{Value: score.SubPoint},
				"AssistPoint": &types.AttributeValueMemberS{Value: score.AssistPoint},
			},
		})
	}

	m := map[string]types.AttributeValue{
		"No":   &types.AttributeValueMemberS{Value: fmt.Sprintf("%d", monster.No)},
		"Name": &types.AttributeValueMemberS{Value: monster.Name},
		"OriginMonsterNo": &types.AttributeValueMemberS{
			Value: fmt.Sprintf("%d", monster.OriginMonsterNo),
		},
		"Game8URL":    &types.AttributeValueMemberS{Value: monster.Game8URL.String()},
		"Game8Scores": &types.AttributeValueMemberL{Value: game8Scores},
	}

	_, err := repo.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &repo.tableName,
		Item:      m,
	})
	if err != nil {
		return err
	}

	return nil
}
