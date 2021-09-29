resource "logzio_endpoint" "microsoftteams" {
  title = "microsoftteams_empty_url"
  endpoint_type = "microsoftteams"
  description = "this is my description"
  microsoftteams {
    url = ""
  }
}