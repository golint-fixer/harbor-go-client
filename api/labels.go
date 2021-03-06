package api

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/moooofly/harbor-go-client/utils"
)

func init() {
	utils.Parser.AddCommand("labels_list",
		"List labels according to the query strings.",
		"This endpoint let user list labels by name, scope and project_id",
		&labelslist)
	utils.Parser.AddCommand("label_create",
		"Post creates a label",
		"This endpoint let user creates a label.",
		&labelcreate)
	utils.Parser.AddCommand("label_del_by_id",
		"Delete the label specified by ID.",
		"Delete the label specified by ID.",
		&labeldel)
	utils.Parser.AddCommand("label_get_by_id",
		"Get the label specified by ID.",
		"This endpoint let user get the label by specific ID.",
		&labelget)
	utils.Parser.AddCommand("label_update",
		"Update the label properties.",
		"This endpoint let user update label properties.",
		&labelupdate)
}

type labelsList struct {
	Name      string `short:"n" long:"name" description:"The label name as filter." default:""`
	Scope     string `short:"s" long:"scope" description:"(REQUIRED) The label scope. Valid values are 'g' and 'p'. 'g' for global labels and 'p' for project labels." required:"yes"`
	ProjectID int    `short:"i" long:"project_id" description:"Relevant project ID, Required when scope is 'p'." default:"0"`
	Page      int    `short:"p" long:"page" description:"The page nubmer, default is 1." default:"1"`
	PageSize  int    `short:"z" long:"page_size" description:"The size of per page, default is 10, maximum is 100." default:"10"`
}

var labelslist labelsList

func (x *labelsList) Execute(args []string) error {
	GetLabels(utils.URLGen("/api/labels"))
	return nil
}

// GetLabels let user list labels by name, scope and project_id
//
// params:
//  name       - The label name.
//  scope      - (REQUIRED) The label scope. Valid values are g and p. g for global labels and p for project labels.
//  project_id - Relevant project ID, required when scope is p.
//  page       - The page nubmer, default is 1.
//  page_size  - The size of per page, default is 10, maximum is 100.
//
// operation format:
//  GET /labels
//
// e.g. curl -X GET --header 'Accept: application/json' 'https://localhost/api/labels?scope=g&page=1&page_size=10'
//
func GetLabels(baseURL string) {
	targetURL := baseURL + "?scope=" + labelslist.Scope +
		"&name=" + labelslist.Name +
		"&project_id=" + strconv.Itoa(labelslist.ProjectID) +
		"&page=" + strconv.Itoa(labelslist.Page) +
		"&page_size=" + strconv.Itoa(labelslist.PageSize)

	fmt.Println("==> GET", targetURL)

	// Read beegosessionID from .cookie.yaml
	c, err := utils.CookieLoad()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	utils.Request.Get(targetURL).
		Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
		End(utils.PrintStatus)
}

type labelCreate struct {
	ID           int    `short:"i" long:"id" description:"The ID of label. If not set, automatically generated by harbor." default:"0" json:"id"`
	Name         string `short:"n" long:"name" description:"(REQUIRED) The name of label." required:"yes" json:"name"`
	Description  string `short:"d" long:"description" description:"(REQUIRED) The description of label." required:"yes" json:"description"`
	Color        string `short:"c" long:"color" description:"The color code of label. (e.g. Format: #A9B6BE)" default:"#000000" json:"color"`
	Scope        string `short:"s" long:"scope" description:"The scope of label, 'g' for global labels and 'p' for project labels." default:"g" json:"scope"`
	ProjectID    int    `short:"p" long:"project_id" description:"The project ID if the label is a project label. Required when scope is 'p'." default:"0" json:"project_id"`
	CreationTime string `long:"creation_time" description:"The creation time of label. default time.Now()" default:"" json:"creation_time"`
	UpdateTime   string `long:"update_time" description:"The update time of label. default time.Now()" default:"" json:"update_time"`
	Deleted      bool   `long:"deleted" description:"The label is deleted or not." json:"deleted"`
}

var labelcreate labelCreate

func (x *labelCreate) Execute(args []string) error {
	PostLabelCreate(utils.URLGen("/api/labels"))
	return nil
}

// PostLabelCreate let user creates a label.
//
// params:
//   id            - The ID of label.
//   name          - (REQUIRED) The name of label.
//   description   - (REQUIRED) The description of label.
//   color         - The color code of label. (e.g. Format: #A9B6BE)
//   scope         - The scope of label. ('p' indicates project scope, 'g' indicates global scope)
//   project_id    - Which project id this label belongs to when created. ('0' indicates global label, others indicate specific project)
//   creation_time - The creation time of label. default time.Now()
//   update_time   - The update time of label. default time.Now()
//   deleted       - not sure
//
// format:
//   POST /labels
/*
curl -X POST --header 'Content-Type: application/json' --header 'Accept: text/plain' -d '{ \
   "id": 100, \
   "name": "label-name-100", \
   "description": "label-description-100", \
   "color": "#000000", \
   "scope": "g", \
   "project_id": 0, \
   "deleted": true \
 }' 'https://localhost/api/labels'
*/
func PostLabelCreate(baseURL string) {
	if labelcreate.CreationTime == "" || labelcreate.UpdateTime == "" {
		now := time.Now().Format("2006-01-02T15:04:05Z")
		labelcreate.CreationTime = now
		labelcreate.UpdateTime = now
	}

	targetURL := baseURL
	fmt.Println("==> POST", targetURL)

	// Read beegosessionID from .cookie.yaml
	c, err := utils.CookieLoad()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	t, err := json.Marshal(&labelcreate)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("==> label add:", string(t))

	utils.Request.Post(targetURL).
		Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
		Send(string(t)).
		End(utils.PrintStatus)
}

type labelDel struct {
	ID int `short:"i" long:"id" description:"(REQUIRED) Label ID." required:"yes"`
}

var labeldel labelDel

func (x *labelDel) Execute(args []string) error {
	DeleteLabel(utils.URLGen("/api/labels"))
	return nil
}

// DeleteLabel deletes the label specified by ID.
//
// params:
//   id - Label ID.
//
// operation format:
//  DELETE /labels/{id}
//
// e.g. curl -X DELETE --header 'Accept: text/plain' 'https://localhost/api/labels/100'
//
func DeleteLabel(baseURL string) {
	targetURL := baseURL + "/" + strconv.Itoa(labeldel.ID)

	fmt.Println("==> DELETE", targetURL)

	// Read beegosessionID from .cookie.yaml
	c, err := utils.CookieLoad()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	utils.Request.Delete(targetURL).
		Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
		End(utils.PrintStatus)
}

type labelGet struct {
	ID int `short:"i" long:"id" description:"(REQUIRED) Label ID." required:"yes"`
}

var labelget labelGet

func (x *labelGet) Execute(args []string) error {
	GetLabel(utils.URLGen("/api/labels"))
	return nil
}

// GetLabel gets the label specified by ID.
//
// params:
//   id - Label ID.
//
// operation format:
//  GET /labels/{id}
//
// e.g. curl -X GET --header 'Accept: text/plain' 'https://localhost/api/labels/100'
//
func GetLabel(baseURL string) {
	targetURL := baseURL + "/" + strconv.Itoa(labelget.ID)

	fmt.Println("==> GET", targetURL)

	// Read beegosessionID from .cookie.yaml
	c, err := utils.CookieLoad()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	utils.Request.Get(targetURL).
		Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
		End(utils.PrintStatus)
}

type labelUpdate struct {
	ID          int    `short:"i" long:"id" description:"(REQUIRED) Label ID." required:"yes" json:"id"`
	Name        string `short:"n" long:"name" description:"(REQUIRED) The name of label." required:"yes" json:"name"`
	Description string `short:"d" long:"description" description:"(REQUIRED) The description of label." required:"yes" json:"description"`
	Color       string `short:"c" long:"color" description:"The color code of label. (e.g. Format: #A9B6BE)" default:"#000000" json:"color"`
	Scope       string `short:"s" long:"scope" description:"The scope of label, 'g' for global labels and 'p' for project labels." default:"g" json:"scope"`
	ProjectID   int    `short:"p" long:"project_id" description:"The project ID if the label is a project label. Required when scope is 'p'." default:"0" json:"project_id"`
	//CreationTime string `long:"creation_time" description:"The creation time of label. default time.Now()" default:"" json:"creation_time"`
	//UpdateTime   string `long:"update_time" description:"The update time of label. default time.Now()" default:"" json:"update_time"`
	Deleted bool `long:"deleted" description:"The label is deleted or not." json:"deleted"`
}

var labelupdate labelUpdate

func (x *labelUpdate) Execute(args []string) error {
	PutLabelUpdate(utils.URLGen("/api/labels"))
	return nil
}

// PutLabelUpdate let user update label properties.
//
// params:
//   id            - The ID of label.
//   name          - (REQUIRED) The name of label.
//   description   - (REQUIRED) The description of label.
//   color         - The color code of label. (e.g. Format: #A9B6BE)
//   scope         - The scope of label. ('p' indicates project scope, 'g' indicates global scope)
//   project_id    - Which project id this label belongs to when created. ('0' indicates global label, others indicate specific project)
//   creation_time - The creation time of label. default time.Now()
//   update_time   - The update time of label. default time.Now()
//   deleted       - not sure
//
// operation format:
//   PUT /labels/{id}
/*
curl -X PUT --header 'Content-Type: application/json' --header 'Accept: text/plain' -d '{ \
   "id": 0, \
   "name": "label-name-100", \
   "description": "label-description-100", \
   "color": "#000000", \
   "scope": "g", \
   "project_id": 0, \
   "deleted": true \
 }' 'https://localhost/api/labels/100'
*/
func PutLabelUpdate(baseURL string) {
	// NOTE:
	// Though as swagger shows, both creation_time and creation_time can be updated, but actually not
	/*
		if labelupdate.UpdateTime == "" {
			now := time.Now().Format("2006-01-02T15:04:05Z")
			labelupdate.UpdateTime = now
		}
	*/

	targetURL := baseURL + "/" + strconv.Itoa(labelupdate.ID)
	fmt.Println("==> PUT", targetURL)

	// Read beegosessionID from .cookie.yaml
	c, err := utils.CookieLoad()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	t, err := json.Marshal(&labelupdate)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("==> label add:", string(t))

	utils.Request.Put(targetURL).
		Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
		Set("Content-Type", "application/json").
		Send(string(t)).
		End(utils.PrintStatus)
}
