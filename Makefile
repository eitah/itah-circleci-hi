# set -euo pipefail
# set -x

S3_BUCKET="spicy-omelet-lambda-source"
S3_KEY="circleci-handler/dist.zip"

.PHONY: *

build: circleci

circleci:
	GOOS=$(GOOS) go build -o out/circleci-handler ./cmd/circleci-handler

zip: GOOS=linux
zip: circleci
	zip -Dj ./out/dist.zip ./out/circleci-handler

upload: zip
	aws s3 cp ./out/dist.zip "s3://${S3_BUCKET}/${S3_KEY}"

deploy: upload
	aws lambda update-function-code --function-name circleci-hi-lambda --s3-bucket $(S3_BUCKET) --s3-key $(S3_KEY) --publish --region $(AWS_REGION)


test:
	go test ${PKG}

cover:
	go test -coverprofile coverage.out ${PKG}
	go tool cover -html=coverage.out

clean:
	rm -rf out