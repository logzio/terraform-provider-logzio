variable "api_token" {
  type = string
  description = "your logzio API token"
}

provider "logzio" {
  api_token = var.api_token
}

resource logzio_grafana_notification_policy test_np {
  contact_point = "default-email"
  group_by = ["p8s_logz_name"]
  group_wait      = "50s"
  group_interval  = "7m"
  repeat_interval = "4h"

  policy {
    matcher {
      label = "some_label"
      match = "="
      value = "some_value"
    }
    contact_point = "default-email"
    continue      = true

    group_wait      = "50s"
    group_interval  = "7m"
    repeat_interval = "4h"
    mute_timings = ["some-mute-timing"]


    policy {
      matcher {
        label = "another_label"
        match = "="
        value = "another_value"
      }
      contact_point = "default-email"
    }
  }
}