# F5 Go API
The goal of this project is to develop a REST API application using Go, in preparation for the remainder of the F5 internship. The theme of the project is an API which models the elevators in the F5 Tower in Seattle, with 2 sets of 6 elevators which can access a range of floors. This project will involve an environment where clients send HTTP requests to the server, whose endpoints are listed in the following table:

| HTTP Verb | Endpoint Name | Endpoint Description                                                                                                    | Arguments                                         | Return Value                |
|-----------|---------------|-------------------------------------------------------------------------------------------------------------------------|---------------------------------------------------|-----------------------------|
| POST      | /callelevator | User calls an elevator to a given floor, and the elevator goes down to the destination floor                            | Starting floor, destination floor                 | Success/Error               |
| GET       | /elevatorinfo | Provides information about an elevator (its current floor, its range of accessible floors, and whether it's in transit) | Elevator name                                     | Elevator information/Error  |
| GET       | /allinfo      | Provides information about all the elevators                                                                            | None                                              | Elevators information/Error |
| GET       | /ping         | Returns string "PONG!"                                                                                                  | None                                              | String                      |
| UPDATE    | /update       | Update an elevator's info                                                                                               | Elevator name, New lower & upper bounds | Success/Error               |

# Approach
1. Create an application which can run elevator logic locally
2. Set up a server which can handle requests to the elevator logic
3. Run the server locally and start with the GET requests to ensure functionality
4. Add the other endpoints and ensure local functionality
5. Get the server running on Docker
6. Write Postman tests (or use Gorat if there's time?)
