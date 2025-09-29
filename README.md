# A take-home assessment for the Go Engineer role at Chatsheet
Here’s your take-home assessment for the Go Engineer role at Chatsheet. As mentioned, you’ll have 72 hours from receipt to complete this. Please submit your work through this form:

Go Engineers Assessment Form

## Task:
Implement a small full-stack app that lets a user connect their LinkedIn account using Unipile's native authentication(not the hosted wizard). Once the account is connected, store the returned account_id for that user in a database. Please add both the cookie auth and the username/password login

### Requirements (high-level):
#### Backend
- Expose an endpoint to trigger the Unipile LinkedIn connect flow.
- On success, capture and persist the account_id for the user.
- Provide a simple endpoint to fetch the user's stored accounts.
  
#### Database
- Keep a table with users and their linked accounts (user_id, provider, account_id).
- Use any relational DB (SQLite, Postgres, MySQL).

#### Frontend
- A basic form to submit either credentials or cookie information for LinkedIn.
- A page to display the saved account_id values for the current user.

### Deliverables:
- Source code for backend and frontend.(github repo link)
- Database schema/migrations.(github repo link)
- A hosted website url with the finished app (website url)

### Note:
Please use github as version tracking and display progress and commits as you progress through the assessment. Please submit a repo link that include both the source code + migration, and a url that host your website


## How to Start
### env
- golang v1.23.12+
- node.js v22.18.0+
- docker

### Development Mode
#### Database
```bash
# Start DB
docker-compose up
```

#### Frontend (svelte + vite)
Under ./web/myapp
```bash
npm install

# Run the frontend dev server on http://localhost:5173
npm run dev
# or build dist package for all-in-one fullstack
npm run build
```

#### Backend (golang)
Put your unipile ***api_key*** & ***api_base_url*** into ./config/config.yml first.
```bash
# The server will be run on http://localhost:8080
go run ./cmd/myapp/main.go
```

#### Test
- If you use frontend dev server, please access http://localhost:5173
- If you choose all-in-one fullstack, please access http://localhost:8080

### Docker Mode
Put your unipile ***api_key*** & ***api_base_url*** into ./config/config.yml first.
```bash
# build
docker build -f Dockerfile-mono -t chatsheet:latest .

# run
docker run -it --rm --name chatsheet -p 8080:8080 -p 5432:5432 chatsheet:latest
```
And the access to http://localhost:8080
