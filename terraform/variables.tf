variable "name" {
  type    = string
  default = "slackbot"
}

variable "s3_bucket" {
  type        = string
  description = "Name of the s3 bucket where the lambda artifact is located."
}

variable "s3_key" {
  type        = string
  description = "S3 key of the lambda artifact."
}

variable "slack_oauth_access_token" {
  type        = string
  description = "Slack app oath access token."
}

variable "slack_verification_token" {
  type        = string
  description = "Slack app verification token."
}

variable "iam_policy" {
  type        = string
  description = "Policy for the slackbot lambda in the format of the json attribute from an aws_iam_policy_document."
}

