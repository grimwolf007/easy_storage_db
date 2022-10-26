# Install
- Clone the repo
- Download it as a zip and extract it

# Setup
`docker compose up`
 - In the future I will make a script to store the passwords more securely and inject them into the docker-compose

# Teardown
`docker compose down`

# Uninstall
`docker compose down`
`docker compose rm`
remove used images
`docker rmi minio`
`docker rmi postgres`
`docker rmi adminer`
`docker rmi easy_storage_webapp`


# Road Map
 - [x]  Milestone 1
   - [x] Find a golang web framework - Gin
   - [x] Make a basic webserver
   - [x] Make a pipeline that compiles it into a container and tests it
   - [x] Make it s3 compatible - Minio
 - [ ] Milestone 2
   - [x] Store minio service account environment in a file
   - [ ] Create webapp
     - [x] upload a file to the webapp
     - [ ] Upload an object to minio
       - [ ] and store details in postgres
       - [ ] and check if the name is already used
       - [ ] and check if file already exists and you just want to add the name as a new tag
     - [ ] Re-upload an object
       - [ ] and store N revisions
     - [ ] Edit object details
     - [ ] Download object
     - [ ] Clear object Details
     - [ ] Delete object and object details 
       - Soft delete? only fully delete when low on space
       - Force delete for things that need to be removed
       - Admin can see soft deleted items if they need to be force deleted or semi-restored
 - [ ] Milestone 3
   - [ ] use imagemagic to see pes images thumbnails
   - [ ] Add tagging
   - [ ] Add sorting and filtering
   - [ ] Automatic Backups
   - [ ] Create tool to auto upload all files and tag them accordingly
   - [ ] Basic Logging, access, uploads
   - [ ] Authentication
 - [ ] Milestone 4
   - [ ] Add ZAP fuzzy testing to pipeline
   - [ ] Add docker container scanning to pipeline
   - [ ] Add golang static scanning to pipeline




# Database Schema
file <--> thumbnail
file <--> name
file <==> tags
file <-=> taglist
tag-list <=-> tag 
