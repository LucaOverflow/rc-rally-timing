package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	autoAddTransponderOnPassing(app)

	requireDecoderOnStageActivation(app)
	validateDifferentDecoderOnStageActivation(app)
	validateDecoderNotInUseOnStageActivation(app)
	removeDecoderOnStageDeactivation(app)

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

func requireDecoderOnStageActivation(app *pocketbase.PocketBase) {
	app.OnRecordUpdateRequest("stages").BindFunc(func(e *core.RecordRequestEvent) error {
		if e.Record.GetBool("active") {
			if e.Record.GetString("start_decoder") == "" || e.Record.GetString("stop_decoder") == "" {
				return e.BadRequestError("Decoder must be set when activating the stage", nil)
			}
		}

		return e.Next()
	})
}

func validateDifferentDecoderOnStageActivation(app *pocketbase.PocketBase) {
	app.OnRecordUpdateRequest("stages").BindFunc(func(e *core.RecordRequestEvent) error {
		if e.Record.GetBool("active") {
			if e.Record.GetString("start_decoder") == e.Record.GetString("stop_decoder") {
				return e.BadRequestError("Start and stop decoder must be diffrent", nil)
			}
		}

		return e.Next()
	})
}

func validateDecoderNotInUseOnStageActivation(app *pocketbase.PocketBase) {
	app.OnRecordUpdateRequest("stages").BindFunc(func(e *core.RecordRequestEvent) error {
		if e.Record.GetBool("active") {
			decoderStringArray := e.Record.GetString("start_decoder") + "," + e.Record.GetString("stop_decoder")

			_, err := app.FindFirstRecordByFilter("stages",
				"start_decoder = {:startDecoder} || start_decoder = {:stopDecoder} || stop_decoder ~ {:startDecoder} || stop_decoder = {:stopDecoder}",
				dbx.Params{
					"decoder":      decoderStringArray,
					"startDecoder": e.Record.GetString("start_decoder"),
					"stopDecoder":  e.Record.GetString("stop_decoder"),
				})

			if err != sql.ErrNoRows {
				return e.BadRequestError("Decoder already in use on a different stage", nil)
			}
		}

		return e.Next()
	})
}

func removeDecoderOnStageDeactivation(app *pocketbase.PocketBase) {
	app.OnRecordUpdateRequest("stages").BindFunc(func(e *core.RecordRequestEvent) error {
		if !e.Record.GetBool("active") {
			e.Record.Set("start_decoder", "")
			e.Record.Set("stop_decoder", "")
		}

		return e.Next()
	})
}
