.PHONY: build deploy

build:
	cd src && \
	GOOS=linux GOARCH=amd64 go build -o bootstrap
	sam build
	cp -r src/templates/ .aws-sam/build/GetMonsterRatingFunction/templates/

deploy:
	sam deploy \
	--stack-name puzzdra-monster-rating-sam \
	--region ap-northeast-1 \
	--confirm-changeset \
	--capabilities CAPABILITY_IAM \
	--no-disable-rollback \
	--parameter-overrides "GetMonsterRatingFunctionAuth=NONE"

# endpointを取得する
endpoints:
	sam list endpoints --stack-name puzzdra-monster-rating-sam --region ap-northeast-1

# 削除コマンド
delete:
	sam delete --stack-name puzzdra-monster-rating-sam --region ap-northeast-1

# ローカル実行
local-start-api:
	sam local start-api --env-vars local.env.json
local-invoke:
	sam local invoke GetMonsterRatingFunction --event ./aws/sam/events/event.json

# DynamoDB =====================================================================

# ローカルでDynamoDBを起動する
local-dynamodb-table:
	aws dynamodb create-table \
		--table-name MonsterRatings \
		--attribute-definitions AttributeName=No,AttributeType=S \
		--key-schema AttributeName=No,KeyType=HASH \
		--provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
		--endpoint-url http://localhost:8000

# ローカルでテストデータ投入（プロジェクトルートで実行することを前提）
local-dynamodb-seed-test-data:
	aws dynamodb batch-write-item \
	--region ap-northeast-1 \
	--endpoint-url http://localhost:8000 \
	--request-items file://dynamodb/testdata/monster_ratings.json

# dynamodbにテストデータを投入する
dynamodb-seed-test-data:
	aws dynamodb batch-write-item \
	--region ap-northeast-1 \
	--request-items file://dynamodb/testdata/monster_ratings.json

# テストデータ取得
local-dynamodb-scan:
	aws dynamodb scan \
	--table-name MonsterRatings \
	--region ap-northeast-1 \
	--endpoint-url http://localhost:8000

# 以下参考
# aws dynamodb get-item \
#   --table-name MyDynamoDBTable \
#   --region ap-northeast-1 \
#   --key '{
#     "id": {"S": "12345"}
#   }'

# aws dynamodb query \
#   --table-name MyDynamoDBTable \
#   --region ap-northeast-1 \
#   --key-condition-expression "id = :id" \
#   --expression-attribute-values '{
#     ":id": {"S": "12345"}
#   }'
# ```
