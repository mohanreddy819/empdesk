
variable "region" {
  description = "The region to deploy resources in"
  type        = string
  default     = "us-central1"
}
variable "project_id" {
  description = "The GCP project ID"
  type        = string
}
variable "zone" {
  description = "The zone to deploy resources in"
  type        = string
  default     = "us-central1-a"
  
}