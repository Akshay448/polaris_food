# Polaris Food Project Documentation
To keep things simple and get started calling the apis
1. clone the repo
2. import the postman [collection](Polaris-Food.postman_collection.json) in your postman client
3. run the apis using [docker](#Running in docker) or run through golang to see live changes 
   with sqlitedb which can be accessed with any jetbrains ide
4. Once the api server starts at localhost:8080, start using the apis through postman
5. To understand the database structure with relations, refer to [draw.io](https://drive.google.com/file/d/1vPhfVjy2-TqiDGi45_u8HJC4tAGrdHIV/view?usp=sharing)

## Project Structure

- **/polaris-food**
    - **/api** - REST API endpoint handlers
        - **/v1** - Version 1 of the API
            - **/handlers** - HTTP handlers
            - **/middlewares** - API middlewares
    - **/cmd** - Main applications for this project
        - **/food-delivery-main** - The web server executable
    - **/config** - Configuration files for db and external service
    - **/internal**
        - **/auth** - Authentication and authorization logic
        - **/database** - Database interactions
            - **/models** - Data models representing tables
        - **/service** - Business logic
            - **/user** - User-related business logic
            - **/rider** - Rider-related business logic
            - **/restaurant** - Restaurant-related business logic
            - **/order** - Order-related business logic
        - **/util** - Utility functions and shared code across services
    - **/tests** - Unit tests and integration tests
        - **/unit**
        - **/integration**
    - **Dockerfile** - To containerize the application
    - **README.md**

## Points to note before running the project
* A simple sqlite database (file fooddelivery.db) is being used here for demo purposes, it's not 
scalable solution but just to represent the working code
* If run through docker, it will be difficult to access the db, because docker is putting 
everything inside the container including db file, and vendor directory which makes it heavy
* if you want to access the db simultaneously while calling the apis, tun the server in local 
  system using go installed in your system, mentioned the commands below


## Running in docker

### Build and run the application using Docker:
takes about 5 minutes to build the images, please be patient

```bash
docker build -t food-delivery-app .
docker run -p 8080:8080 --name my-food-delivery-app food-delivery-app`
```
After the docker image is run, access the api endpoints at localhost:8080

## Running locally with go
### Prerequisites - go 1.21 or later is installed

```bash
go build -o ./out/food-delivery-main ./cmd/food-delivery-main
chmod +x ./out/food-delivery-main
./out/food-delivery-main
```
After the project is run, access the api endpoints at localhost:8080


## API Endpoints
The API is organized around RESTful principles, providing access to resources such as users,
riders, restaurants, orders, and ratings.

Access API documentation here, generated through postman - 

### User Endpoints

* POST /api/v1/register/user: Register a new user.
* GET /api/v1/users/:id/orders: Fetch a user's order history.
* GET /api/v1/users/:id/coupons: Retrieve coupons available to a user.

### Rider Endpoints
* POST /api/v1/register/rider: Register a new rider.
* GET /api/v1/riders/nearest/:restaurantId: Find the nearest rider to a restaurant.
* PUT /api/v1/riders/:id/location: Update a rider's location.
* GET /api/v1/riders/:id/orders: Fetch a rider's order history.

### Restaurant Endpoints
* POST /api/v1/register/restaurant: Register a new restaurant.
* GET /api/v1/restaurants/suggest: Suggest restaurants to a user.
* GET /api/v1/restaurants/:id/menu: Provide the menu of a restaurant.

### Order Endpoints
* POST /api/v1/orders/create: Create a new order.
* POST /api/v1/orders/update: Update an existing order.

### Ratings Endpoints
* POST /api/v1/ratings/submit: Submit a rating.

## Services
The services are designed to interact with each other or with external services or with a database
Currently interacting with database only
* UserService
    * RegisterUser
    * GetUserOrders
    * GetUserCoupons
* RiderService
    * RegisterRider
    * UpdateRiderLocation
    * GetRiderOrders
    * GetNearestAvailableRider
* OrderService
    * CreateOrder
    * UpdateOrder
* RestaurantService
    * RegisterRestaurants
    * GetMenu
    * SuggestRestaurants
* RatingService
    * SubmitRatings

## Database Tables structure
### Coupon Table

| Column Name    | Data Type        | Constraints |
|----------------|------------------|-------------|
| ID             | INTEGER          | PRIMARY KEY |
| CreatedAt      | DATETIME         |             |
| UpdatedAt      | DATETIME         |             |
| DeletedAt      | DATETIME         |             |
| Code           | TEXT             |             |
| Description    | TEXT             |             |
| DiscountType   | TEXT             |             |
| DiscountValue  | REAL             |             |
| ValidFrom      | DATETIME         |             |
| ValidUntil     | DATETIME         |             |
| MinOrderValue  | REAL             |             |
| Active         | BOOLEAN          |             |


### FoodCategory

| Column Name  | Data Type | Constraints |
|--------------|-----------|-------------|
| ID           | INTEGER   | PRIMARY KEY |
| CreatedAt    | DATETIME  |             |
| UpdatedAt    | DATETIME  |             |
| DeletedAt    | DATETIME  |             |
| Name         | TEXT      |             |
| Description  | TEXT      |             |

### MenuItem

| Column Name   | Data Type | Constraints                |
|---------------|-----------|----------------------------|
| ID            | INTEGER   | PRIMARY KEY                |
| CreatedAt     | DATETIME  |                            |
| UpdatedAt     | DATETIME  |                            |
| DeletedAt     | DATETIME  |                            |
| RestaurantID  | INTEGER   | FOREIGN KEY (Restaurant)   |
| Name          | TEXT      |                            |
| Description   | TEXT      |                            |
| Price         | REAL      |                            |
| CategoryID    | INTEGER   | FOREIGN KEY (FoodCategory) |


### Order

| Column Name      | Data Type | Constraints              |
|------------------|-----------|--------------------------|
| ID               | INTEGER   | PRIMARY KEY              |
| CreatedAt        | DATETIME  |                          |
| UpdatedAt        | DATETIME  |                          |
| DeletedAt        | DATETIME  |                          |
| UserID           | INTEGER   | FOREIGN KEY (USER)       |
| RestaurantID     | INTEGER   | FOREIGN KEY (Restaurant) |
| RiderID          | INTEGER   | FOREIGN KEY (Rider)      |
| Status           | TEXT      |                          |
| TotalPrice       | REAL      |                          |
| DeliveryAddress  | TEXT      |                          |
| CouponId         | INTEGER   |                          |

### OrderItem

| Column Name  | Data Type | Constraints            |
|--------------|-----------|------------------------|
| ID           | INTEGER   | PRIMARY KEY            |
| CreatedAt    | DATETIME  |                        |
| UpdatedAt    | DATETIME  |                        |
| DeletedAt    | DATETIME  |                        |
| OrderID      | INTEGER   | FOREIGN KEY (Order)    |
| MenuItemID   | INTEGER   | FOREIGN KEY (MenuItem) |
| Quantity     | INTEGER   |                        |
| Price        | REAL      |                        |


### Rating

| Column Name | Data Type | Constraints         |
|-------------|-----------|---------------------|
| ID          | INTEGER   | PRIMARY KEY         |
| CreatedAt   | DATETIME  |                     |
| UpdatedAt   | DATETIME  |                     |
| DeletedAt   | DATETIME  |                     |
| OrderID     | INTEGER   | FOREIGN KEY (Order) |
| RatedByID   | INTEGER   | FOREIGN KEY (User)  |
| RatedToID   | INTEGER   | FOREIGN KEY (User)  |
| Stars       | REAL      |                     |
| Comment     | TEXT      |                     |


### Restaurant

| Column Name   | Data Type | Constraints |
|---------------|-----------|-------------|
| ID            | INTEGER   | PRIMARY KEY |
| CreatedAt     | DATETIME  |             |
| UpdatedAt     | DATETIME  |             |
| DeletedAt     | DATETIME  |             |
| Name          | TEXT      |             |
| Address       | TEXT      |             |
| DeliveryTime  | INTEGER   |             |
| IsOpen        | BOOLEAN   |             |
| Latitude      | REAL      |             |
| Longitude     | REAL      |             |


### RiderProfile

| Column Name         | Data Type | Constraints        |
|---------------------|-----------|--------------------|
| ID                  | INTEGER   | PRIMARY KEY        |
| AvailabilityStatus  | BOOLEAN   |                    |
| IsDelivering        | BOOLEAN   |                    |
| UserID              | INTEGER   | FOREIGN KEY (User) |
| CreatedAt           | DATETIME  |                    |
| UpdatedAt           | DATETIME  |                    |
| Latitude            | REAL      |                    |
| Longitude           | REAL      |                    |


### User

| Column Name    | Data Type | Constraints |
|----------------|-----------|-------------|
| ID             | INTEGER   | PRIMARY KEY |
| CreatedAt      | DATETIME  |             |
| UpdatedAt      | DATETIME  |             |
| DeletedAt      | DATETIME  |             |
| Username       | TEXT      |             |
| Email          | TEXT      | UNIQUE      |
| PasswordHash   | TEXT      |             |
| Role           | TEXT      |             |
| AverageRating  | REAL      |             |


### UserCoupon

| Column Name | Data Type | Constraints          |
|-------------|-----------|----------------------|
| ID          | INTEGER   | PRIMARY KEY          |
| CreatedAt   | DATETIME  |                      |
| UpdatedAt   | DATETIME  |                      |
| DeletedAt   | DATETIME  |                      |
| UserID      | INTEGER   | FOREIGN KEY (User)   |
| CouponID    | INTEGER   | FOREIGN KEY (Coupon) |
| IsUsed      | BOOLEAN   |                      |
