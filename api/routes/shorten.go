package routes

import (
	"time"
	"strconv"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/AdityaNarayan05/shorten-url/database"
	"github.com/AdityaNarayan05/shorten-url/helpers"
	"github.com/satori/go.uuid"
	"github.com/asaskevich/govalidator"
)

// request represents the incoming JSON request structure.
type request struct {
	URL     string        `json:"url"`
	CustomShort string    `json:"short"`
	Expiry  time.Duration `json:"expiry"`
}

// response represents the JSON response structure.
type response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit"`
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
}

// ShortenURL handles the endpoint for shortening URLs.
func ShortenURL(c *fiber.Ctx) error {
	// Parse the incoming JSON request
	body := new(request)
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	// Implement rate limiting
	r2 := database.CreateClient(1)
	defer r2.Close()
	val, err := r2.Get(database.Ctx, c.IP()).Result()

	if err == redis.Nil {
		_ = r2.Set(database.Ctx, c.IP, os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else {
		val, _ := r2.Get(database.Ctx, c.IP()).Result()
		valInt, _ := strconv.Atoi(val)

		if valInt <= 0 {
			limit, _ := r2.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":            "Rate Limit exceeded",
				"rate_limit_reset": limit / time.Second / time.Minute,
			})
		}
	}

	// Check if the input is a valid URL
	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Invalid URL"})
	}

	// Check for domain error
	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Invalid Domain Error"})
	}

	// Enforce HTTPS, SSL
	body.URL = helpers.EnforceHTTP(body.URL)

	var id string

	if body.CustomShort == "" {
		id = uuid.NewV4().String()[:6]
	} else {
		id = body.CustomShort
	}

	r := database.CreateClient(0)
	defer r.Close()

	val, _ = r.Get(database.Ctx, c.IP()).Result()
	if val != "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "URL custom short is already in use",
		})
	}

	if body.Expiry == 0 {
		body.Expiry = 24
	}

	err = r.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to connect to server",
		})
	}

	resp := response{
		URL:             body.URL,
		CustomShort:     "",
		Expiry:          body.Expiry,
		XRateRemaining:  10,
		XRateLimitReset: 30,
	}

	r2.Decr(database.Ctx, c.IP())

	val, _ = r2.Get(database.Ctx, c.IP()).Result()
	resp.XRateRemaining, _ = strconv.Atoi(val)

	ttl, _ := r2.TTL(database.Ctx, c.IP()).Result()
	resp.XRateLimitReset = ttl / time.Second / time.Minute

	resp.CustomShort = os.Getenv("Domain") + "/" + id

	return c.Status(fiber.StatusOK).JSON(resp)
}
