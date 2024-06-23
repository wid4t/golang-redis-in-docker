package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var rdb *redis.Client
var ctx = context.Background()

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "redis-clusterip.server.svc.cluster.local:6379",
		DB:   0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
}

func main() {

	initRedis()

	app := fiber.New()

	api := app.Group("/module/redis", logger.New())

	api.Get("/check", func(c *fiber.Ctx) error {
		return c.SendString("Hello, i'm from module golang-redis-in-docker")
	})

	api.Get("/set/:key/:value", func(c *fiber.Ctx) error {
		key := c.Params("key")
		value := c.Params("value")

		err := rdb.Set(ctx, key, value, 0).Err()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return c.SendString(fmt.Sprintf("Key %s set to %s", key, value))
	})

	api.Get("/get/:key", func(c *fiber.Ctx) error {
		key := c.Params("key")

		value, err := rdb.Get(ctx, key).Result()
		if err == redis.Nil {
			return c.Status(fiber.StatusNotFound).SendString("Key not found")
		} else if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return c.SendString(fmt.Sprintf("Key %s has value %s", key, value))
	})

	log.Fatal(app.Listen(":3000"))

}
