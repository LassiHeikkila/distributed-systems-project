# accountservice
`accountservice` is responsible for handling user account related data, 
i.e.: 
- sign up
- login/logout
- request authentication checking

# API
Needs to have endpoints to:
- Create account
- Read/get account details (by id or username)
- Update account details
- Delete account
- Login (get session token)
- Logout (invalidate session token)
- Logout everywhere (invalidate all session tokens for user)

# Questions
- During sign up, do we want to validate email by sending a mail with link to click?
  - Maybe in the future