variable "api_token" {
  type = string
  description = "your logzio API token"
}

provider "logzio" {
  api_token = var.api_token
}

resource "logzio_endpoint" "slack_endpoint" {
  title = "slack_endpoint"
  description = "hello"
  endpoint_type = "slack"
  slack {
    url = "https://jsonplaceholder.typicode.com/todos/1"
  }
}

resource "logzio_endpoint" "custom_endpoint" {
  title = "custom_endpoint"
  endpoint_type = "custom"
  description = "hello"
  custom {
    url = "https://jsonplaceholder.typicode.com/todos/1"
    method = "POST"
    headers = {
      this = "is"
      a = "header"
    }
    body_template = jsonencode({
      this: "is"
      my: "template"
      nested: {
        also: "working"
      }
    })
  }
}

resource "logzio_endpoint" "pagerduty_endpoint" {
  title = "pagerduty_endpoint"
  endpoint_type = "pagerduty"
  description = "hello"
  pagerduty {
    service_key = "my_service_key"
  }
}

resource "logzio_endpoint" "bigpanda_endpoint" {
  title = "bigpanda_endpoint"
  endpoint_type = "bigpanda"
  description = "hello"
  bigpanda {
    api_token = "my_api_token"
    app_key = "my_app_key"
  }
}

resource "logzio_endpoint" "datadog_endpoint" {
  title = "datadog_endpoint"
  endpoint_type = "datadog"
  description = "hello"
  datadog {
    api_key = "my_api_key"
  }
}

resource "logzio_endpoint" "victorops_endpoint" {
  title = "victorops_endpoint"
  endpoint_type = "victorops"
  description = "hello"
  victorops {
    routing_key = "my_routing_key"
    message_type = "my_message_type"
    service_api_key = "my_service_api_key"
  }
}

resource "logzio_endpoint" "opsgenie_endpoint" {
  title = "opsgenie_endpoint"
  endpoint_type = "opsgenie"
  description = "hello"
  opsgenie {
    api_key = "my_api_key"
  }
}

resource "logzio_endpoint" "servicenow_endpoint" {
  title = "servicenow_endpoint"
  endpoint_type = "servicenow"
  description = "hello"
  servicenow {
    username = "my_username"
    password = "my_password"
    url = "https://jsonplaceholder.typicode.com/todos/1"
  }
}

resource "logzio_endpoint" "microsoftteams_endpoint" {
  title = "microsoftteams_endpoint"
  endpoint_type = "microsoftteams"
  description = "hello"
  microsoftteams {
    url = "https://jsonplaceholder.typicode.com/todos/1"
  }
}