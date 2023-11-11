package routes

import(
	"time"
	"github.com/AdityaNarayan05/shorten-url/database"
	"os"
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
	r2 := database.CreateClient(1)
	defer r2.Close()
	val,err := r2.Get(database.Ctx,c.IP()).Result()

	if err==redis.Nil{
		_=r2.Set(database.Ctx,c.IP,os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else {
		val,_ := r2.Get(database.Ctx,c.IP().Result())
		valInt,_ := strconv.Atoi(val)
	
		if valInt <= 0{
			limit,_ :=r2.TTL(database.Ctx,c.IP()).Result()
			return c.status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":"Rate Limit exceeded",
				"rate_limit_reset":limit / time.timeNanosecond/time.Minute,
			})
		}
	}

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

	r2.Decr(database.Ctx,c.IP());
} 



