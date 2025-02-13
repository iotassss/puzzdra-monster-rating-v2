package repository

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/iotassss/puzzdra-monster-rating-v2/internal/entity"
)

type MonsterRepository struct {
	client    DynamoDBAPI
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

func NewMonsterRepositoryWithRemoteDB() *MonsterRepository {
	remoteDBClient := newRemoteDynamoDBClient()

	tableNameMonsters := os.Getenv("DYNAMODB_TABLE_NAME_REMOTE")

	return &MonsterRepository{
		client:    remoteDBClient,
		tableName: tableNameMonsters,
	}
}

func DynamoDBOutputToEntity(dynamoDBOutput *dynamodb.GetItemOutput) (entity.Monster, error) {
	dbMonster := dynamoDBOutput.Item
	if len(dbMonster) == 0 {
		return entity.Monster{}, nil
	}

	noStr := dbMonster["No"].(*types.AttributeValueMemberN).Value
	no, err := strconv.Atoi(noStr)
	if err != nil {
		return entity.Monster{}, err
	}
	originMonsterNoStr := dbMonster["OriginMonsterNo"].(*types.AttributeValueMemberN).Value
	originMonsterNo, err := strconv.Atoi(originMonsterNoStr)
	if err != nil {
		return entity.Monster{}, err
	}

	game8URLStr := dbMonster["Game8URL"].(*types.AttributeValueMemberS).Value
	game8URL, err := url.Parse(game8URLStr)
	if err != nil {
		return entity.Monster{}, err
	}

	game8ScoresLen := len(dbMonster["Game8Scores"].(*types.AttributeValueMemberL).Value)
	game8Scores := make([]entity.Game8MonsterScore, 0, game8ScoresLen)
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

func (repo *MonsterRepository) FindByNo(ctx context.Context, no int) (entity.Monster, error) {
	result, err := repo.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &repo.tableName,
		Key: map[string]types.AttributeValue{
			"No": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", no)},
		},
	})
	if err != nil {
		return entity.Monster{}, err
	}

	return DynamoDBOutputToEntity(result)
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

func (r *MonsterRepository) SaveAllMonsters(ctx context.Context, monsters []*entity.Monster) error {
	const batchSize = 25 // DynamoDB BatchWriteItem supports up to 25 items per request

	for i := 0; i < len(monsters); i += batchSize {
		end := i + batchSize
		if end > len(monsters) {
			end = len(monsters)
		}

		batch := monsters[i:end]
		if err := r.writeBatch(ctx, batch); err != nil {
			return fmt.Errorf("failed to save batch %d-%d: %w", i, end, err)
		}
	}

	return nil
}

func (r *MonsterRepository) writeBatch(ctx context.Context, batch []*entity.Monster) error {
	writeRequests := make([]types.WriteRequest, len(batch))
	for i, monster := range batch {
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
		monsterItem := map[string]types.AttributeValue{
			"No":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", monster.No)},
			"Name": &types.AttributeValueMemberS{Value: monster.Name},
			"OriginMonsterNo": &types.AttributeValueMemberN{
				Value: fmt.Sprintf("%d", monster.OriginMonsterNo),
			},
			"Game8URL": &types.AttributeValueMemberS{Value: monster.Game8URL.String()},
			"Game8Scores": &types.AttributeValueMemberL{
				Value: game8Scores,
			},
			"BatchProcessStatus": &types.AttributeValueMemberN{Value: "1"},
		}

		writeRequests[i] = types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: monsterItem,
			},
		}
	}

	maxRetries := 5
	for attempt := 0; attempt < maxRetries; attempt++ {
		if len(writeRequests) == 0 {
			return nil
		}

		input := &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				r.tableName: writeRequests,
			},
		}

		out, err := r.client.BatchWriteItem(ctx, input)
		if err != nil {
			return err
		}

		// 失敗したアイテムを再試行
		writeRequests = out.UnprocessedItems[r.tableName]

		// すべて成功したら終了
		if len(writeRequests) == 0 {
			return nil
		}

		// バックオフを入れる (指数バックオフなどを検討)
		time.Sleep(time.Duration(attempt*100) * time.Millisecond)
	}

	return fmt.Errorf("batch write failed after %d attempts", maxRetries)
}

func (r *MonsterRepository) ScanAll(ctx context.Context) ([]*entity.Monster, error) {
	var items []*entity.Monster
	var lastEvaluatedKey map[string]types.AttributeValue

	for {
		// スキャンリクエスト
		output, err := r.client.Scan(ctx, &dynamodb.ScanInput{
			TableName:         aws.String(r.tableName),
			ExclusiveStartKey: lastEvaluatedKey,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to scan items: %w", err)
		}

		// レスポンスからエンティティに変換
		for _, item := range output.Items {
			monster, err := DynamoDBOutputToEntity(&dynamodb.GetItemOutput{Item: item})
			if err != nil {
				return nil, fmt.Errorf("failed to convert item: %w", err)
			}
			items = append(items, &monster)
		}

		// LastEvaluatedKeyが存在しない場合、終了
		if output.LastEvaluatedKey == nil {
			break
		}
		lastEvaluatedKey = output.LastEvaluatedKey
	}

	return items, nil
}
