.PHONY: build deploy

build:
	sam build -t ./aws/sam/template.yaml
	cp -r templates/ .aws-sam/build/GetMonsterRatingFunction/templates/

deploy:
	sam deploy \
	--stack-name puzzdra-monster-rating-sam \
	--region ap-northeast-1 \
	--confirm-changeset \
	--capabilities CAPABILITY_IAM \
	--no-disable-rollback \
	--parameter-overrides "GetMonsterRatingFunctionAuth=NONE" \
	--config-file ./aws/sam/samconfig.toml

# dynamodbにテストデータを投入する
dynamodb-seed-test-data:
	aws dynamodb batch-write-item \
	--region ap-northeast-1 \
	--request-items '{
		"MonsterRatings": [
			{
				"PutRequest": {
					"Item": {
						"id": {"S": "12345"},
						"Name": {"S": "Sample Name"},
						"Age": {"N": "30"}
					}
				}
			},
			{
				"PutRequest": {
					"Item": {
						"id": {"S": "67890"},
						"Name": {"S": "Another Name"},
						"Age": {"N": "25"}
					}
				}
			}
		]
	}'

# endpointを取得する
endpoints:
	sam list endpoints --stack-name puzzdra-monster-rating-sam --region ap-northeast-1

# 削除コマンド
delete:
	sam delete --stack-name puzzdra-monster-rating-sam --region ap-northeast-1

# テストデータ取得
dynamodb-scan:
	aws dynamodb scan \
	--table-name MonsterRatings \
	--region ap-northeast-1
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

# ローカル実行
local-start-api:
	sam local start-api --env-vars local.env.json
local-invoke:
	sam local invoke GetMonsterRatingFunction --event ./aws/sam/events/event.json
