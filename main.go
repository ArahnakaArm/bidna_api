package main

import (
	"fmt"
	"gofiber/db"
	"gofiber/routes"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s \n", err))
	}
}

func main() {
	app := fiber.New()
	app.Use(cors.New())

	/* 	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	   	if err != nil {
	   		panic(err)
	   	}
	   	defer func() {
	   		if err = client.Disconnect(context.TODO()); err != nil {
	   			panic(err)
	   		}
	   	}()

	   	dbctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	   	err = client.Connect(dbctx) */
	/* 	// Ping the primary
	   	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
	   		panic(err)
	   	} */
	/* 	fmt.Println("Successfully connected and pinged.")

	   	db := client.Database("products")
	   	collection := db.Collection("products") */
	/* 	filter := bson.M{"name": "test1"}
	   	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	   	cur ,err := collection.Find(ctx, bson.M{}).Decode(&result)
	   	if err != nil {

	   	}
	*/
	/* 	cursor, err := collection.Find(ctx, bson.M{})
	   	if err != nil {
	   		log.Fatal(err)
	   	}
	   	defer cursor.Close(ctx)
	   	for cursor.Next(ctx) {
	   		var episode bson.M
	   		if err = cursor.Decode(&episode); err != nil {
	   			log.Fatal(err)
	   		}
	   		fmt.Println(episode)
	   	} */

	/* 	var episodes []Result
	   	cursor, err := collection.Find(dbctx, bson.M{})
	   	if err != nil {
	   		panic(err)
	   	}
	   	if err = cursor.All(dbctx, &episodes); err != nil {
	   		panic(err.Error())
	   	}

	   	fmt.Println(episodes)

	   	b, err := json.Marshal(episodes)
	   	if err != nil {
	   		fmt.Println(err)
	   		return
	   	}
	   	fmt.Println(string(b)) */
	/* 	print(result.Name)
	 */
	/* findCursor, findErr := collection.Find(context.TODO(), bson.D{})
	if findErr != nil {
		panic(findErr)
	}
	var findResults []bson.M
	if findErr = findCursor.All(context.TODO(), &findResults); findErr != nil {
		panic(findErr)
	}
	for _, result := range findResults {
		fmt.Println(result)
	} */

	/* 	app.Get("/login", func(c *fiber.Ctx) error {
	   		url := googleOauthConfig.AuthCodeURL(oauthStateString)
	   		return c.Redirect(url, fiber.StatusTemporaryRedirect)

	   	})
	*/
	/* app.Get("/callback", func(c *fiber.Ctx) error {
		if c.Query("state") != oauthStateString {
			fmt.Printf("state is not valid")
			return c.Redirect("/", fiber.StatusTemporaryRedirect)
		}

		token, err := googleOauthConfig.Exchange(oauth2.NoContext, c.Query("code"))

		if err != nil {
			fmt.Printf("could not get token %s\n", err.Error())
		}

		response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)

		if err != nil {
			fmt.Println("could not get request")
			return c.Redirect("/", fiber.StatusTemporaryRedirect)
		}

		defer response.Body.Close()

		content, err := ioutil.ReadAll(response.Body)

		if err != nil {
			fmt.Printf("could not parse ")
			return c.Redirect("/", fiber.StatusTemporaryRedirect)
		}

		return c.Send(content)

	}) */

	/* fmt.Println("Connected to user MongoDB!") */
	database := db.ConnectMongoDB()

	routes.AddCommonRoute(app)
	routes.AddProductsRoute(app, database)
	routes.AddUsersRoute(app)
	routes.AddGoogleAuthRoute(app)

	app.Listen(fmt.Sprintf("%s:%s", viper.GetString("app.ip"), viper.GetString("app.port")))
	/* http.ListenAndServe(":8000", nil) */
}

type Result struct {
	Name   string   `bson:"name,omitempty"`
	Colors []string `bson:"colors,omitempty"`
}
