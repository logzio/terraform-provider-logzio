resource "logzio_endpoint" "%s" {
  title = "my_custom_title"
  endpoint_type = "Custom"
  description = "this_is_my_description"
  custom {
    url = "https://www.test.com"
    method = "POST"
    headers = {
      this = "is"
      a = "header"
    }
    body_template = {
      this = "is"
      my = "template"
    }
  }
}