# ETIAssg1
This file will explain the design decisions for this program, as well as instructions on how to set up and run them.

# Design consideration of microservices
In this application, we use microservices to add new modules without redesigning the system structure. Unlike monolithic servers, microservices are more flexible, allowing for a wider variation of data types to be used. They are also more highly reliable when compared to a monolithic server. During production, we have to ensure multiple factors is achieved. 
Scalability: Our carpool service will have a volatile user count, fluctuating with people joining or leaving the service. As such,we have a need to introduce load balancing and automatic scaling systems to manage fluctuating workloads among microservices. Each microservice should be scalable on its own, allowing for efficient allocation of resources based on demand.
Communication Between Services: Usng REST API, we can scale our systems better as it depends on client-server interactions, allowing users to get, add, edit and delete records.
Security: Using microservices we can separate our users to ensure only relevant data pertaining to them can be retrieved, as well as minimizing the chances of a data breach between user categories.

# Architecture diagram
![carpool](https://github.com/suicopath/ETIAssg1/assets/84904561/4aecb194-a919-413a-bbf6-64dc34074c7e)
In our current design, we have a console front-end which can be accessed through our API gateway. Using REST API, users can access our services separated into different micro services. We have microservices dedicated to Passengers and Carowners,as well as another for storing trips. Each microservice is has their own data so there is miimal risk of data leakage.

# Instructions for setting up and running the microservices
Download the files in the repository.
1. Open and run the SQL script in MySQL. This will set up the database with pre-registered data.
2. Run carpoolserver.go with a chosen code editor of your choice.
3. Run carpoolclient.go in a separate window.
4. Open the Command Prompt, and enter the curl statement.
5. The UI should appear in the carpoolserver terminal.
6. To select a category, enter the corresponding number.
7. Multiple options should appear. Select the option you need by entering the corresponding number.
8. Fill in the information required by the program.
9. When you are done, press 9 to exit the program.

