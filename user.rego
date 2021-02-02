package opa.user

import input
default allow=false

#GET /user/data
allow {
  input.method == "GET"
  input.path = ["user", "data"]
  data.roles["admin"][_] == input.payload.username
}
