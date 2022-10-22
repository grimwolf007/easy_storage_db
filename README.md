#Setup

#Teardown

#Road Map
 - [ ]  Milestone 1
   - [x] Find a golang web framework - Gin
   - [x] Make a basic webserver
   - [x] Make a pipeline that compiles it into a container and tests it
   - [x] Make it s3 compatible - Minio
 - [ ] Milestone 2
   - [ ] Create webapp
     - [ ] Upload an object to minio and store details in postgres
     - [ ] Re-upload an object and store N revisions
     - [ ] Edit object details
     - [ ] Download object
     - [ ] Clear object Details
     - [ ] Delete object and object details 
       - Soft delete? only fully delete when low on space
       - Force delete for things that need to be removed
       - Admin can see soft deleted items if they need to be force deleted or semi-restored
     - [ ] Basic Logging, access, uploads
     - [ ] Authentication
   - [ ] Add sorting and filtering
   - [ ] Add tagging
   - [ ] Display PNGs as images
 - [ ] Milestone 3
   - [ ] Create tool to auto upload all files and tag them accordingly
