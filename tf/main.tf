data "aws_caller_identity" "current" {}

resource "aws_s3_bucket" "lambda_source" {
  bucket = "spicy-omelet-lambda-source"
  acl    = "public-read" // todo: yes this need to be fixed
}

# if copy pasting may need to 1) `touch null` and 2) `zip null.zip null`
# $ echo "package main" > null.go
# $ zip null.zip null.go
# $ aws s3 cp null.zip s3://spicy-omelet-lambda-source/null.zip
resource "aws_s3_bucket_object" "lambda_src_null_zip" {
  key                    = "null.zip"
  bucket                 = aws_s3_bucket.lambda_source.id
  source                 = "${path.module}/null.zip"
  server_side_encryption = "aws:kms"
}

resource "aws_lambda_function" "lambda" {
  function_name = "circleci-hi-lambda"
  role          = aws_iam_role.lambda.arn
  description   = "circleci-hi-lambda"
  handler       = "circleci-handler"
  runtime       = "go1.x"
  timeout       = 10

  s3_bucket = aws_s3_bucket.lambda_source.bucket
  s3_key    = "null.zip"

  lifecycle {
    ignore_changes = [
      layers,
      environment
    ]
  }
}

// IAM
data "aws_iam_policy_document" "assume_role" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "lambda" {
  name               = "lambda-role"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

resource "aws_iam_role_policy_attachment" "service_role" {
  role       = "lambda-role"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# API Gateway
resource "aws_apigatewayv2_api" "gateway" {
  name          = "circleci-hi-lambda-API"
  protocol_type = "HTTP"
  description   = "This is the api gateway trigger for the hi lambda"
  tags          = {}
}

resource aws_apigatewayv2_route notifications {
  api_id    = aws_apigatewayv2_api.gateway.id
  route_key = "POST /notifications"
  target    = "integrations/${aws_apigatewayv2_integration.notifications.id}"
  # target = aws_apigatewayv2_integration.notifications.id
}

resource aws_apigatewayv2_integration notifications {
  api_id           = aws_apigatewayv2_api.gateway.id
  integration_type = "AWS_PROXY"

  description          = "notifications route for slack"
  integration_method   = "POST"
  # passthrough_behavior = "WHEN_NO_MATCH"
  integration_uri      = aws_lambda_function.lambda.invoke_arn
}

# Adds a prod stage to thid but its still not a working terraform implementation i used the console to create
resource "aws_apigatewayv2_stage" "prod" {
  api_id = aws_apigatewayv2_api.gateway.id
  name   = "prod"
}

// todo get this bound to the gateway - at the moment running will just break permissions I think
resource "aws_lambda_permission" "lambda_permission" {
  statement_id = "AllowMyDemoAPIInvoke"
  action       = "lambda:InvokeFunction"
  // todo this is ARN not function name so this will fail
  function_name = aws_lambda_function.lambda.arn
  principal     = "apigateway.amazonaws.com"

  # The /*/*/* part allows invocation from any stage, method and resource path
  # within API Gateway REST API. the last one indicates where to send requests to.
  # see more detail https://docs.aws.amazon.com/lambda/latest/dg/services-apigateway.html
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*/${aws_lambda_function.lambda.function_name}"
}