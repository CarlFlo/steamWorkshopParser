package server

import (
	"context"

	protos "github.com/CarlFlo/steamWorkshopDownloader/protos/workshopParser"
	"github.com/CarlFlo/steamWorkshopDownloader/steamWorkshop"
)

type MyWorkshopParserServer struct {
	protos.UnimplementedWorkshopParserServer
}

func (wps MyWorkshopParserServer) ParseWorkshopItem(ctx context.Context, req *protos.Request) (*protos.Response, error) {

	output, err := steamWorkshop.RunProgram(req.Url)

	return &protos.Response{
		Result: output,
	}, err
}
