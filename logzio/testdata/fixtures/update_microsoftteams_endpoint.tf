resource "logzio_endpoint" "microsoftteams" {
  title = "updated_microsoftteams_title"
  endpoint_type = "microsoftteams"
  description = "this is my description"
  microsoftteams {
    url = "https://jsonplaceholder.typicode.com/todos/2"
  }
}