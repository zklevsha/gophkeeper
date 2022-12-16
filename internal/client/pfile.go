package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/zklevsha/gophkeeper/internal/helpers"
	"github.com/zklevsha/gophkeeper/internal/pb"
	"github.com/zklevsha/gophkeeper/internal/structs"
)

// 1 MB
const PFILE_MAX_SIZE = 1000000

// pfileAdd sends pfile to server
func pfileAdd(mstorage *structs.MemStorage, ctx context.Context, gclient *structs.Gclient) {
	err := reqCheck(mstorage)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// getting list of pfiles
	pfiles, err := listDir(mstorage.PfilesDir)
	if err != nil {
		log.Printf("ERROR: cant list pfile directory (%s): %s\n",
			mstorage.PfilesDir, err.Error())
		return
	}
	if len(pfiles) == 0 {
		log.Printf("pfile directory is empty (%s)", mstorage.PfilesDir)
		return
	}

	// reading client`s input and loading file
	pfilePath := inputSelect("pfile to send", pfiles)
	tags, err := getTags(getInput(`metainfo: {"key":"value",...}`, isTags, false))
	if err != nil {
		log.Printf("ERROR: cant parse tags: %s\n", err.Error())
		return
	}
	data, err := loadFile(pfilePath)
	if err != nil {
		log.Printf("ERROR: cant load file: %s", err.Error())
		return
	}
	pfileName := filepath.Base(pfilePath)
	pfile := structs.Pfile{Name: pfileName, Data: data, Tags: tags}

	// sending data to server
	pdata, err := helpers.ToPdata("pfile", pfile, mstorage.MasterKey)
	if err != nil {
		log.Printf("cant convert pfile to pdata: %s", err.Error())
		return
	}
	_, err = gclient.Pdata.AddPdata(ctx, &pb.AddPdataRequest{Pdata: pdata})
	if err != nil {
		log.Printf("cannot send pfile to server: %s", err.Error())
		return
	}

	log.Println("pfile was added")
}

// pfile retrrives user`s private file from server
func pfileGet(mstorage *structs.MemStorage, ctx context.Context, gclient *structs.Gclient) {
	err := reqCheck(mstorage)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// getting list of existing pfile entries
	entries, err := listPnames(ctx, gclient, "pfile")
	if err != nil {
		log.Printf("ERROR: cant retrive list of existing pfile entries: %s", err.Error())
	}
	if len(entries) == 0 {
		log.Println("You dont have any pstring entries")
		return
	}
	var pnames []string
	for pname := range entries {
		pnames = append(pnames, pname)
	}

	// parsing input
	pname := inputSelect("Pfile name ", pnames)
	pdataID := entries[pname]

	// getting pdata
	resp, err := gclient.Pdata.GetPdata(ctx, &pb.GetPdataRequest{PdataID: pdataID})
	if err != nil {
		log.Printf("ERROR: cant retrive pdata from server: %s\n", err.Error())
		return
	}

	// decoding to pfile
	cleaned, err := helpers.FromPdata(resp.Pdata, mstorage.MasterKey)
	if err != nil {
		log.Printf("ERROR: cant decode pdata: %s\n", err.Error())
		return
	}
	pfile := cleaned.(structs.Pfile)

	//saving pfile
	pfileName := fmt.Sprintf("%s-%d", pfile.Name, time.Now().UnixNano())
	pfilePath := path.Join(mstorage.PfilesDir, pfileName)
	err = os.WriteFile(pfilePath, pfile.Data, 0600)
	if err != nil {
		log.Printf("cant save pfile %s: %s", pfilePath, err.Error())
		return
	}

	// prinring output
	log.Printf("pfile saved at %s", pfilePath)
	if len(pfile.Tags) == 0 {
		return
	}
	log.Println("pfile tags:")
	tags_pretty, err := json.MarshalIndent(pfile.Tags, "", " ")
	if err != nil {
		log.Printf("ERROR cant encode tags to JSON : %s\n", err.Error())
	} else {
		log.Println(string(tags_pretty))
	}

}

// pfileUpdate update pfile entry
func pfileUpdate(mstorage *structs.MemStorage, ctx context.Context, gclient *structs.Gclient) {
	err := reqCheck(mstorage)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// getting current pfile
	entries, err := listPnames(ctx, gclient, "pfile")
	if err != nil {
		log.Printf("ERROR: cant retrive list of existing pfile entries: %s", err.Error())
	}
	if len(entries) == 0 {
		log.Printf("You dont have any pfile entries")
		return
	}
	var pnames []string
	for pname := range entries {
		pnames = append(pnames, pname)
	}
	pname := inputSelect("Pstring to update: ", pnames)
	pdataID := entries[pname]
	getResp, err := gclient.Pdata.GetPdata(ctx, &pb.GetPdataRequest{PdataID: pdataID})
	if err != nil {
		log.Printf("ERROR: cant retrive pdata from server: %s\n", err.Error())
		return
	}
	cleaned, err := helpers.FromPdata(getResp.Pdata, mstorage.MasterKey)
	if err != nil {
		log.Printf("ERROR: cant convert pdata to pfile: %s", err.Error())
		return
	}
	pfileOld := cleaned.(structs.Pfile)

	// getting list of local pfiles
	pfiles, err := listDir(mstorage.PfilesDir)
	if err != nil {
		log.Printf("ERROR: cant list pfile directory (%s): %s\n",
			mstorage.PfilesDir, err.Error())
		return
	}
	if len(pfiles) == 0 {
		log.Printf("pfile directory is empty (%s)", mstorage.PfilesDir)
		return
	}

	//getting new tags from user
	tagsJson, err := json.Marshal(pfileOld.Tags)
	if err != nil {
		log.Printf("ERROR: cant parse old tags: %s\n", err.Error())
		return
	}
	tagsStr := getInput(fmt.Sprintf("new tags[%s]", tagsJson), isTags, false)
	var tagsNew map[string]string
	if tagsStr != "" {
		tagsNew, err = getTags(tagsStr)
		if err != nil {
			log.Printf("ERROR: cant convert tags %s\n", err.Error())
			return
		}
	} else {
		tagsNew = pfileOld.Tags
	}

	// getting new pfile
	var data []byte
	var pfileName string
	pfiles = append(pfiles, "dont change pfile")
	pfilePath := inputSelect("new pfile", pfiles)
	if pfilePath == "dont change pfile" {
		data = pfileOld.Data
		pfileName = pfileOld.Name
	} else {
		data, err = loadFile(pfilePath)
		if err != nil {
			log.Printf("ERROR: cant load file: %s", err.Error())
			return
		}
		pfileName = filepath.Base(pfilePath)
	}

	// Sending new pfile to server
	pfile := structs.Pfile{Name: pfileName, Data: data, Tags: tagsNew}
	pdata, err := helpers.ToPdata("pfile", pfile, mstorage.MasterKey)
	if err != nil {
		log.Printf("cant convert pfile to pdata: %s", err.Error())
		return
	}
	pdata.ID = getResp.Pdata.ID
	_, err = gclient.Pdata.UpdatePdata(ctx, &pb.UpdatePdataRequest{Pdata: pdata})
	if err != nil {
		log.Printf("cannot send pfile to server: %s", err.Error())
		return
	}

	log.Println("pfile was updated")

}

// loadFile checks file size and loads it from disk
func loadFile(fpath string) ([]byte, error) {
	// checking file size
	fStat, err := os.Stat(fpath)
	if err != nil {
		return []byte{}, fmt.Errorf("cannot get file info: %s", err.Error())
	}
	if fStat.Size() > PFILE_MAX_SIZE {
		return []byte{}, fmt.Errorf("file size cannot be larger than %d bytes", PFILE_MAX_SIZE)
	}

	// loading file
	data, err := os.ReadFile(fpath)
	if err != nil {
		return []byte{}, fmt.Errorf("cant read file %s", err.Error())
	}
	return data, nil
}
