package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/zklevsha/gophkeeper/internal/pb"
	"github.com/zklevsha/gophkeeper/internal/structs"
)

// reqCheck checks that user is logged in and has master key loaded
func reqCheck(mstorage *structs.MemStorage) error {
	if mstorage.MasterKey.Key == "" {
		return fmt.Errorf("master-key does not exists\n" +
			"add it via key-generate/key-load commands")

	}
	if mstorage.Token == "" {
		return errors.New("login required (login)")
	}
	return nil
}

// listPnames getting list of available pnames
// list returned as <pname>: <pid> map
func listPnames(ctx context.Context, gclient *structs.Gclient, ptype string) (map[string]int64, error) {
	listResponse, err := gclient.Pdata.ListPdata(ctx, &pb.ListPdataRequest{Ptype: ptype})
	if err != nil {
		return map[string]int64{}, err
	}

	entries := make(map[string]int64)
	for _, e := range listResponse.PdataEtnry {
		entries[e.Name] = e.ID
	}

	return entries, nil
}
