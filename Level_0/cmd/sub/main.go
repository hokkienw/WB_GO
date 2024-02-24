package main

import(
    "level_0/pkg/cache"
    "level_0/pkg/database"
    "level_0/pkg/server"
    "level_0/pkg/stanpkg"
    "log"
    "time"
)

func main(){
db, err := database.InitDB()
if err != nil {
	log.Fatalf("Error initializing database: %v", err)
}
defer db.Close()

err = database.CreateTableDB(db)
if err != nil {
	log.Fatalf("Error creating table in database: %v", err)
}

cache := cache.New(15*time.Minute, 20*time.Minute)

err = database.GetFromDB(db, cache)
if err != nil {
	log.Fatalf("Error while getting data from database: %v", err)
}

go server.RunServer(cache)

stanpkg.SubscribeAndListen(cache, db)

}