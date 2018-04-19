package evry.com.lunch.authz

test_allow_read {
    allow with input as {"method": "GET", "path": "/menu"}
}

test_allow_admin_write {
    allow with input as {"method": "POST", "path": "/menu", "jwt": {"role": "lunchAdmin"}}
}

test_deny_non_admin_write {
    not allow with input as {"method": "POST", "path": "/menu"}
}