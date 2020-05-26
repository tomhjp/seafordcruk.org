provider "aws" {
  version = "~> 2.60"
  region  = "eu-west-1"
  profile = "default"
}

variable "tag" {}

resource "aws_lambda_function" "seafordcruk-org-contact-us-lambda" {
  function_name    = "seafordcruk-org-contact-us-lambda"
  filename         = "${path.module}/handler.zip"
  handler          = "handler"
  role             = aws_iam_role.seafordcruk-org-contact-us-lambda-role.arn
  runtime          = "go1.x"
  source_code_hash = filebase64sha256("${path.module}/handler.zip")
  tags             = { version = var.tag }
}

resource "aws_cloudwatch_log_group" "seafordcruk-org-contact-us-lambda" {
  name              = "/aws/lambda/${aws_lambda_function.seafordcruk-org-contact-us-lambda.function_name}"
  retention_in_days = 30
}

resource "aws_iam_role" "seafordcruk-org-contact-us-lambda-role" {
  name = "seafordcruk-org-contact-us-lambda-role"
  path = "/service-role/"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_policy" "seafordcruk-org-contact-us-lambda-logging-policy" {
  name        = "seafordcruk-org-contact-us-lambda-policy"
  path        = "/"
  description = "IAM policy for logging and sending email from the seafordcruk.org 'contact us' form handler"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:*:*:*",
      "Effect": "Allow"
    },
    {
      "Action": [
        "ses:SendEmail"
      ],
      "Resource": "arn:aws:ses:*:703450008913:identity/*@gmail.com",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "seafordcruk-org-contact-us-lambda-logs" {
  role       = aws_iam_role.seafordcruk-org-contact-us-lambda-role.name
  policy_arn = aws_iam_policy.seafordcruk-org-contact-us-lambda-logging-policy.arn
}