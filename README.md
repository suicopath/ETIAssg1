# ETIAssg1
This file will explain the design decisions for this program, as well as instructions on how to set up and run them.

# Design consideration of microservices
In this application, we use microservices to add new modules without redesigning the system structure. Unlike monolithic servers, microservices are more flexible, allowing for a wider variation of data types to be used. They are also more flexible, and highly reliable compared to a monolithic server. During production, we have to ensure multiple factors is achieved. 
Scalability: Our carpool service will have a volatile user count, fluctuating with people joining or leaving the service. As such,we have a need to introduce load balancing and automatic scaling systems to manage fluctuating workloads among microservices. Each microservice should be scalable on its own, allowing for efficient allocation of resources based on demand.
Communication Between Services: Usng REST API we can scale our systems better as it depends on client-server interactions, allowing users to get, add, edit and delete records.
Security: Using microservices we can separate our users to ensure only relevant data pertaining to them can be retrieved, as well as minimizing the chances of a data breach between user categories.

# Architecture diagram
![carpool](https://github.com/suicopath/ETIAssg1/assets/84904561/4aecb194-a919-413a-bbf6-64dc34074c7e)
In our current design, we have a console front-end which can be accessed through our API gateway. Using REST API, users can access our services separated into different micro services. We have a microservice dedicated to Passengers and 

# Instructions for setting up and running your microservices
Download the files in the repository.
1. Run carpoolserver.go
2. Run carpoolclient.go in a separate window.
3. Open the Command Prompt, and enter the curl statement.
4. The UI should appear in the carpoolserver terminal.
5. To select a category, enter the corresponding number.
6. Multiple options should appear. Select the option you need by entering the corresponding number.
