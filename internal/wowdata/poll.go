package wowdata

import "context"

func (app *App) Run(ctx context.Context) error {
	// plan work: what dynamic data needs to be fetched? what static data is missing or needs refetch? what imports are staged?
	// do work: call a bnet api, check lastmodified, write to s3 if new; or import from s3 to pg
	return nil
}
