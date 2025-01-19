package repository

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"

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
