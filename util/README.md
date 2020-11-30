prerequisite:-
   install programming language go 1.14
   install postgres

   Make sure your postgres should be running with exactly same configured database(username,password,hostname,port and dbname) is in config.go
   

   To get all dependend pakages run 'go get' command

   To host this application: 
                     Navigate to diretory where main.go placed 
                     now 'go run main.go' command

Once application start, all tables eill get created.
To start with open http://localhost:12345/signIn in Postnam and provide valid username and password
you can change port number from main.go
List of urls are availabe in main.go


Sample demo screen-shots attached in Demo-screenshot.zip file. 

