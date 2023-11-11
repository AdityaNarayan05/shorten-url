package routes

import(
	"time"
)
type request struct {
	URL                string              'json:"url"'
	CustomShort	       string              'json:"short"'
	Expiry             time.Duration       'json:"expiry"'
}

type response struct {
	URL                string              'json:"url"'
	CustomShort        string              'json:"short"'
	Expiry             time.Duration       'json:"expiry"'
	XRateRemaining     int                 'json:"rate_limit"'
	XRateLimitRest     time.Duration       'json:"rate_limit_reset"'
}


func ShortenURL(c *fiber.ctx)  error{

	body := new(request)
	if err := c.BodyParser(&body); err != nil {
		return c.status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"cannot parse JSON"})
	}

	//implement rate limiting
	

	//check if the input is actual URL or not
	if !govalidator.IsURL(body.URL){
		return c.status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Invalid URL"})
	}

	//check for domain error
	if !helpers.RemoveDomainError(body.URL){
		return c.status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Invalid Domain Error"})
	}

	//enforce https, SSL
	body.URL=helpers.EnforceHTTP(body.URL)
} 



