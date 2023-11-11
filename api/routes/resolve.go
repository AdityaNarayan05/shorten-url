package routes

import (
	"github.com/AdityaNarayan05/shorten-url/database"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

// ResolveURL handles the resolution of a short URL to its original URL.
func ResolveURL(c *fiber.Ctx) error {
	// Extract the short URL from the request parameters
	shortURL := c.Params("url")

	// Create a Redis client for the main database
	r := database.CreateClient(0)
	defer r.Close()

	// Retrieve the original URL from the database
	value, err := r.Get(database.Ctx, shortURL).Result()

	if err == redis.Nil {
		// If the short URL is not found in the database, return a 404 response
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "short URL not found in DB",
		})
	} else if err != nil {
		// If there's an error connecting to the database, return a 500 response
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "cannot connect to DB",
		})
	}

	// Create a Redis client for the rate limiter
	rInr := database.CreateClient(1)
	defer rInr.Close()

	// Increment the counter in the rate limiter
	_ = rInr.Incr(database.Ctx, "counter")

	// Redirect the client to the original URL with a 301 status code
	return c.Redirect(value, 301)
}