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

	"github.com/zklevsha/gophkeeper/internal/pb"
)

const pfileMaxSize = 1000000

// Pfile represents user`s private file
type Pfile struct {
	Name string
	Data []byte
	Tags map[string]string
}

// pfileAdd sends pfile to server
func pfileAdd(ctx context.Context, mstorage *MemStorage, gclient *Gclient) {
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
	pfile := Pfile{Name: pfileName, Data: data, Tags: tags}

	// sending data to server
	pdata, err := toPdata("pfile", pfile, mstorage.MasterKey)
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
func pfileGet(ctx context.Context, mstorage *MemStorage, gclient *Gclient) {
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
	cleaned, err := fromPdata(resp.Pdata, mstorage.MasterKey)
	if err != nil {
		log.Printf("ERROR: cant decode pdata: %s\n", err.Error())
		return
	}
	pfile := cleaned.(Pfile)

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
	tagsPretty, err := json.MarshalIndent(pfile.Tags, "", " ")
	if err != nil {
		log.Printf("ERROR cant encode tags to JSON : %s\n", err.Error())
	} else {
		log.Println(string(tagsPretty))
	}

}

// pfileUpdate update pfile entry
func pfileUpdate(ctx context.Context, mstorage *MemStorage, gclient *Gclient) {
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
	cleaned, err := fromPdata(getResp.Pdata, mstorage.MasterKey)
	if err != nil {
		log.Printf("ERROR: cant convert pdata to pfile: %s", err.Error())
		return
	}
	pfileOld := cleaned.(Pfile)

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
	tagsJSON, err := json.Marshal(pfileOld.Tags)
	if err != nil {
		log.Printf("ERROR: cant parse old tags: %s\n", err.Error())
		return
	}
	tagsStr := getInput(fmt.Sprintf("new tags[%s]", tagsJSON), isTags, false)
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
	pfile := Pfile{Name: pfileName, Data: data, Tags: tagsNew}
	pdata, err := toPdata("pfile", pfile, mstorage.MasterKey)
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
	if fStat.Size() > pfileMaxSize {
		return []byte{}, fmt.Errorf("file size cannot be larger than %d bytes", pfileMaxSize)
	}

	// loading file
	data, err := os.ReadFile(fpath)
	if err != nil {
		return []byte{}, fmt.Errorf("cant read file %s", err.Error())
	}
	return data, nil
}

// pfileDetele deletes pstring entry
func pfileDelete(ctx context.Context, mstorage *MemStorage, gclient *Gclient) {
	err := reqCheck(mstorage)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// getting list of available pnames
	entries, err := listPnames(ctx, gclient, "pfile")
	if err != nil {
		log.Printf("ERROR: cant retrive list of existing pfile entries: %s", err.Error())
	}
	if len(entries) == 0 {
		log.Printf("You dont have any pfile entries")
		return
	}
	// parsing input
	var pnames []string
	for pname := range entries {
		pnames = append(pnames, pname)
	}
	pname := inputSelect("pstring name: ", pnames)
	if !getYN(fmt.Sprintf("do you want delete %s?", pname)) {
		log.Println("Canceled")
		return
	}
	_, err = gclient.Pdata.DeletePdata(ctx, &pb.DeletePdataRequest{PdataID: entries[pname]})
	if err != nil {
		log.Printf("ERROR: cant delete pdata: %s", err.Error())
		return
	}
	log.Printf("pfile %s was deleted", pname)
}