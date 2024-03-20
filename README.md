# Mood Harbour
A project for my Software Engineering course. Frontend Repo at https://github.com/anirudhgray/mood-harbour-frontend
## Features
- [x] Built using Golang, Gin, Gorm and PostgreSQL.
- [x] Dockerised via docker-compose.
- [x] Auth: Login, Register, Forgot Password, Reset Password, Delete Account
- [x] Mood Tracking Features: Add Mood Entries at any time, and see your mood history.
- [x] Facial Expression Detection to detect your mood in real time.
- [x] Publish helpful resources, and vote on resources present in the community.
- [x] Get personalised recommendations via collaborative-filtering for mood related resources.
- [x] Send mails for Auth related matters.
- [x] Certain users can have admin access.
## Running Locally
0. Ensure that you have Docker installed on your system.
1. Clone the repo.
2. Create a `.env` in the project root, and copy `.env.example` into it. Generate a random API_SECRET, and obtain a mailtrap API key for email sending.
3. `cd mood-harbour-backend`
4. `make dev`
5. This will spin up the docker-compose, and the containers needed to run the dev project (the main server, the database server, and one container for the DB admin view).
