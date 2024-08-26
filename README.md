### API responsible for Saving Notes

Attributes Captured

- Title
- Description
- Date
- Status
  - Active
  - InActive

### DB_TABLE: APP_USER.T_USER_NOTES

Columns:

- ID
- DESCRIPTION
- CREATED_AT
- STATUS

### Need following env entries in .env file

export API_PORT=8080
export DB_PORT=1521
export DB_HOST=localhost
export DB_NAME=XEPDB1
export DB_USER=app_user
export DB_PASS=
