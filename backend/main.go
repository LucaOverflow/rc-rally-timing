package main

import (
	"log"
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	autoAddTransponderOnPassing(app)

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		// serves static files from the provided public dir (if exists)
		se.Router.GET("/{path...}", apis.Static(os.DirFS("./pb_public"), false))

		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

func autoAddTransponderOnPassing(app *pocketbase.PocketBase) {
	app.OnRecordCreateRequest("passings").BindFunc(func(e *core.RecordRequestEvent) error {
		_, err := app.FindRecordById("transponder", e.Record.GetString("transponder"))

		if err != nil {
			collection, err := app.FindCollectionByNameOrId("transponder")
			if err != nil {
				return e.InternalServerError("Could not find transponder collection", nil)
			}

			newTransponder := core.NewRecord(collection)
			newTransponder.Set("id", e.Record.GetString("transponder"))
			newTransponder.Set("type", e.Record.GetString("transponder_type"))

			err = app.Save(newTransponder)
			if err != nil {
				return e.InternalServerError("Could not create new transponder record", nil)
			}
		}

		return e.Next()
	})
}
