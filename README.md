# RPC-Auth

## How to use

To implement the project, you must run both of this commands:

```
go get golang.org/x/crypto/bcrypt

go get github.com/go-sql-driver/mysql
```

Also, you must have installed MySQL and Go, and create a new table in your local database using the following command:

```
CREATE TABLE users(
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50),
    password VARCHAR(120),
    firstname VARCHAR(50),
    lastname VARCHAR(50),
    birthdate VARCHAR(50)
);
```

Inside of the file of controller.go replace in the line 133 with your own credentials:

```
db, err = sql.Open("mysql", "<root>:<password>@/<dbname>")
// Replace with 
db, err = sql.Open("mysql", "myUsername:myPassword@/myDatabase")
```

[This is how you check your username and password if you're using XAMPP.](https://www.javierrguez.com/recuperar-contrasena-de-phpmyadmin-con-xampp/)


## How to run

To run the program, open your powershell or cmd in the folder tmpl of the project, and type this command:

```
go build
```

And then this one

```
controller
```

Finally, open your browser and go to http://localhost:8080/home
