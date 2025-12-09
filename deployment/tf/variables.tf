variable "aws_region" {
  description = "AWS region for all resources."

  type    = string
  default = "eu-west-1"
}

variable "function_name" {
  description = "Lambda function name"
  type        = string
  default     = "svc-health-aggregator-lambda"