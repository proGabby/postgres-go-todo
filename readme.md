# Postgres-Go-Todo

Postgres-Go-Todo is a Todo application built with Go and the Gorilla toolkit. It incorporates key features such as user authentication, permission handling, middlewares, robust error handling, and utilizes a PostgreSQL database. The project follows the MVC (Model-View-Controller) architecture to organize and structure the codebase.

## Features

- **User Authentication:** Secure user authentication system to protect user accounts.
- **Permission Handling:** Granular permission handling based on user roles.
- **Middlewares:** Implementation of essential middlewares for various functionalities.
- **Error Handling:** Robust error handling mechanisms to improve application reliability.
- **PostgreSQL Database:** Utilizes PostgreSQL as the backend database for data storage.
- **MVC Architecture:** Organized codebase following the MVC architectural pattern.

## Getting Started

To start the project, follow these steps:

1. Clone the repository:

   ```bash
   git clone https://github.com/proGabby/postgres-go-todo.git
   cd postgres-go-todo
   ```

2. Create a .env file in the project root and add the following environment variables:

    ```env
    JWT_SECRET_KEY=your_jwt_secret_key
    DB_CONNECTION_STRING=your_db_connection_string
    ```

    Replace `your_jwt_secret_key` and `your_db_connection_string` with your preferred values.
   

3. Initialize Go modules:

   ```bash
   go mod init
   ```

4. Run the main Go file:

   ```bash
   go run main.go
   ```

This will initialize the project and start the application. Visit [http://localhost:8080](http://localhost:8080) in your browser to access the Todo application.

## Project Structure

The project follows the MVC architecture for code organization:

- `models/`: Contains data models and interacts with the database.
- `views/`: Handles the presentation logic and user interface.
- `controllers/`: Manages the application's business logic and orchestrates interactions.
- `middlewares/`: Includes various middlewares for authentication, logging, etc.
- `utils/`: Holds utility functions and helper modules.
- `main.go`: Entry point of the application.

Feel free to explore each directory for more details on the project structure.
