resource "logzio_endpoint" "microsoftteams" {
  title = "my_microsoftteams_title"
  endpoint_type = "microsoftteams"
  description = "this is my description"
  microsoftteams {
    url = "https://jsonplaceholder.typicode.com/todos/1"
  }
}