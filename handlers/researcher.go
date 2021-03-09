package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/omerfruk/my-blog/database"
	"github.com/omerfruk/my-blog/models"
	"github.com/omerfruk/my-blog/service"
	"io/ioutil"
	"os"
	"strings"
)

//db den okuyup sayfaya yansıtmak için
func ResearchersAll(c *fiber.Ctx) error {
	var res []models.Research
	area := c.Params("key")
	temp := strings.ToLower(area)
	res = service.GetResearch(temp)
	return c.Render("researcher", fiber.Map{
		"res": res,
	})
}

//json dosyasından okumak için
func ResearchersAllJson(c *fiber.Ctx) error {
	var research []models.Research
	jsonFile, err := os.Open("fakedata/fake_data.json")
	if err != nil {
		fmt.Println(err)
	}
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(byteValue, &research)
	if err != nil {
		fmt.Println(err)
	}
	var titles [3]string
	var text [3]string
	j := 0
	area := c.Params("key")
	for i := 0; i < len(research); i++ {
		if research[i].Area == area {
			//fmt.Println(research[i].Area)
			titles[j] = research[i].Title
			//fmt.Println(titles[j])
			text[j] = research[i].Text
			//fmt.Println(text[j])
			j++
		}
	}
	return c.Render("researcher", fiber.Map{
		"Info":         "Hello, " + area + " !",
		"explanation":  titles[0],
		"explanation2": titles[1],
		"explanation3": titles[2],
		"Text":         text[0],
		"Text2":        text[1],
		"Text3":        text[2],
	})

}

//login sayfasının renderi

func Login(c *fiber.Ctx) error {

	return c.Render("login", true)
}

//gingup render
func SingUp(c *fiber.Ctx) error {

	return c.Render("singup", true)
}

//sing up post render
func SingUpPost(c *fiber.Ctx) error {
	var request models.RequestSingUp
	err := c.BodyParser(&request)
	if err != nil {
		fmt.Println(err)
	}
	service.SingUPUser(request.Phone, request.FullName, request.Email, request.Password)
	return c.Redirect("/login")
}

func LogControl(c *fiber.Ctx) error {
	var request models.RequestBody
	var temp models.User

	err := c.BodyParser(&request)
	if err != nil {
		fmt.Println(err)
	}
	err = database.DB().Where("mail = ? and password = ?", request.Email, request.Password).First(&temp).Error
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	return c.Redirect("/admin")
}

// jwt token olusturuldu uretimi
/*func LogControl(c *fiber.Ctx) error {
	var request models.RequestBody
	var temp models.User
	app := fiber.New()
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte("secret"),
	}))
	err := c.BodyParser(&request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(request.Email + "   " + request.Password)
	err = database.DB().Where("mail = ? and password = ?", request.Email, request.Password).First(&temp).Error
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	fmt.Println(temp.Fullname)
	fmt.Println(temp.Authority)
	token := jwt.New(jwt.SigningMethodHS256)

	claims:=token.Claims.(jwt.MapClaims)
	claims["name"]=temp.Fullname
	claims["admin"]=temp.Authority
	claims["exp"]=time.Now().Add(time.Hour * 24 * 30).Unix()

	fmt.Println(claims)

	t,err := token.SignedString([]byte("secret"))
	if err != nil{
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(fiber.Map{
		"token":t,
	})
}
func Restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.SendString("Welcome " + name)
}*/

func AdminPage(c *fiber.Ctx) error {
	var temp []models.User
	err := database.DB().Find(&temp).Error
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(temp)
	return c.Render("admin", temp)
}

func EditUser(c *fiber.Ctx) error {
	key := c.Params("key")
	var request models.RequestSingUp

	err := c.BodyParser(&request)
	if err != nil {
		fmt.Println(err)
	}
	service.UpdateUser(key, request.Email, request.FullName, request.Phone)
	return c.Redirect("/admin")
}

func DltUSer(c *fiber.Ctx) error {
	key := c.Params("key")
	fmt.Println(key)

	service.DeleteUser(key)

	return c.Redirect("/admin")
}

func IndexRender(c *fiber.Ctx) error {
	topbar := service.GetTopBar("ÖmFar.")
	entry := service.GetEntry("WELCOME TO MY PAGE")
	intro := service.GetInstructions("Let's learn something about technology")
	footer := service.GetFooter("OMER FARUK TASDEMIR")
	portfolio := service.GetPortfolio()
	a := models.Anasayfa{
		Portfolio: portfolio,
		Entry:     entry,
		Topbar:    topbar,
		Intro:     intro,
		Footer:    footer,
	}
	return c.Render("index", a)
}

//info Render bolumu
func InfoRender(c *fiber.Ctx) error {
	topbar := service.GetTopBar("ÖmFar.")
	user := service.GetUser("Ömer Faruk")
	footer := service.GetFooter("OMER FARUK TASDEMIR")

	I := models.Info{
		Topbar: topbar,
		User:   user,
		Footer: footer,
	}
	fmt.Println(I.Footer.InstaSrc)
	return c.Render("info", I)
}

func Session(c *fiber.Ctx) error {
	store := session.New()
	sess, err := store.Get(c)
	if err != nil {
		panic(err)
	}

	// save session
	defer sess.Save()

	// Get value
	//name := sess.Get("name")

	// Set key/value
	sess.Set("name", "john")

	// Delete key
	sess.Delete("name")

	// Destroy session
	if err := sess.Destroy(); err != nil {
		panic(err)
	}
	fmt.Println(sess)
	return c.Render("/", true)
}
