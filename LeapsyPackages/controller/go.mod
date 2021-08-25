module controller

go 1.14

replace leapsy.com/db => ../db

replace leapsy.com/model => ../model

replace leapsy.com/properties => ../../properties

require (
	github.com/gofiber/fiber v1.14.6
	go.mongodb.org/mongo-driver v1.7.1
	leapsy.com/db v0.0.0-00010101000000-000000000000
	leapsy.com/model v0.0.0-00010101000000-000000000000
	leapsy.com/properties v0.0.0-00010101000000-000000000000
)
