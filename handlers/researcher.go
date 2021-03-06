package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/omerfruk/my-blog/database"
	"github.com/omerfruk/my-blog/models"
	"github.com/omerfruk/my-blog/service"
	"io/ioutil"
	"os"
	"strings"
	"time"
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
	sess, err := store.Get(c)
	if err != nil {
		fmt.Println(err)
	}
	//eger session acık ise girilecek yer
	if sess.Get("temp") != nil {
		// Destroy session
		if err := sess.Destroy(); err != nil {
			panic(err)
		}
		return c.Redirect("/")
	} else { //sessions açık değilse girilecek  method
		return c.Render("login", true)
	}
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

// session olusturma
var store = session.New(session.Config{
	Expiration:   5 * time.Minute,
	CookieName:   "session_id",
	KeyGenerator: utils.UUID,

	//storage problemi olusuyor sonrasında bununla ilgilenilecek
	/*Storage: postgres.New(postgres.Config{
		Host:       "localhost",
		Port:       5432,
		Username:   "postgres",
		Password:   "123",
		Database:   "blog",
		Table:      "fiber_storage",
		GCInterval: 10 * time.Hour,
	}),*/
})

//login kontrol
func LogControl(c *fiber.Ctx) error {
	var request models.RequestBody
	var temp models.User

	err := c.BodyParser(&request)
	if err != nil {
		fmt.Println(err)
	}
	err = database.DB().Where("mail = ? and password = ?", request.Email, request.Password).First(&temp).Error
	if err != nil {
		return c.Redirect("down")
	}

	sess, err := store.Get(c)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(sess.ID())
	defer sess.Save()

	fmt.Println(sess)
	sess.Set("temp", temp.Fullname)
	sess.Get("temp")

	return c.Redirect("/admin")
}

//yanlıs giris ekrani
var sayac int = 0

func DownPage(c *fiber.Ctx) error {
	sayac++
	if sayac >= 3 {
		sayac = 0
		c.Redirect("/singup")
	}
	conn := "Your Login Failed"
	return c.Render("down", conn)

}

// jwt uretimi
/*func LogControlJWT(c *fiber.Ctx) error {
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

//admin sayfasi
func AdminPage(c *fiber.Ctx) error {
	var user []models.User
	inf := service.GetUser("omer faruk")
	topbar := service.GetTopBar("ÖmFar.")
	err := database.DB().Find(&user).Error
	if err != nil {
		fmt.Println(err)
	}
	admin := models.AdminPage{
		Topbar: topbar,
		Info:   inf,
		User:   user,
	}
	return c.Render("admin", admin)
}

//admin sayfasindan kullanicilari editlemek icin
func EditUser(c *fiber.Ctx) error {
	key := c.Params("key")
	var request models.RequestSingUp

	err := c.BodyParser(&request)
	if request.FullName == "" && request.Password == "" && request.Email == "" {
		return c.Redirect("/admin")
	}
	if err != nil {
		fmt.Println(err)
	}
	service.UpdateUser(key, request.Email, request.FullName, request.Phone)
	return c.Redirect("/admin")
}

//admin sayfasindan kullanici silmek icin
func DltUSer(c *fiber.Ctx) error {
	key := c.Params("key")
	fmt.Println(key)

	service.DeleteUser(key)

	return c.Redirect("/admin")
}

//admin sayufasidan  kullanici silmek icin
func CreateUser(c *fiber.Ctx) error {
	var temp models.RequestSingUp
	err := c.BodyParser(&temp)
	if err != nil {
		fmt.Println(err)
	}
	service.SingUPUser(temp.Phone, temp.FullName, temp.Email, "")
	return c.Redirect("/admin")
}

func Gallery(c *fiber.Ctx) error {
	var tempPhoto []models.FotoGallery
	var tempFooter models.FooterBar
	database.DB().Find(&tempPhoto)
	fmt.Println(tempPhoto)
	database.DB().Find(&tempFooter)
	fmt.Println(tempFooter)
	f := models.FotoStruct{
		FotoGallery: tempPhoto,
		FooterBar:   tempFooter,
	}
	return c.Render("gallery", f)

}

//index sayfasnini renderlemek icin
func IndexRender(c *fiber.Ctx) error {

	topbar := service.GetTopBar("ÖmFar.")
	entry := service.GetEntry("WELCOME TO MY PAGE")
	intro := service.GetInstructions("Let's learn something about technology")
	footer := service.GetFooter("OMER FARUK TASDEMIR")
	portfolio := service.GetPortfolio()
	sess, err := store.Get(c)
	if err != nil {
		fmt.Println(err)
	}
	if sess.Get("temp") != nil {
		topbar.Login = "Logout"
	} else {
		topbar.Login = "Login"
	}
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
	return c.Render("info", I)
}
