package evry.com.lunch.authz

default allow = false

allow {
    input.method = "GET"
    input.path = "/menu"
}

allow {
    input.method = "POST"
    input.path = "/menu"
    input.jwt.role = "lunchAdmin"
}