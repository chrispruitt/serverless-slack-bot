resource "aws_iam_role" "main" {
  name               = "${terraform.workspace}-${var.name}-lambda"
  assume_role_policy = data.aws_iam_policy_document.assume.json
}

resource "aws_iam_role_policy" "main" {
  name   = "${terraform.workspace}-${var.name}-lambda"
  role   = aws_iam_role.main.id
  policy = var.iam_policy
}

data "aws_iam_policy_document" "assume" {
  statement {
    actions = [
      "sts:AssumeRole",
    ]

    principals {
      type = "Service"

      identifiers = [
        "events.amazonaws.com",
        "lambda.amazonaws.com",
      ]
    }
  }
}

