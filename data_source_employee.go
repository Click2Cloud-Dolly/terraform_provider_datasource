package Employee

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"net/http"

	"strconv"

	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)


type Response []struct {
	ID        int         `json:"ID"`
	CreatedAt time.Time   `json:"CreatedAt"`
	UpdatedAt time.Time   `json:"UpdatedAt"`
	DeletedAt interface{} `json:"DeletedAt"`
	UserID    int         `json:"UserId"`
	Username  string      `json:"Username"`
	Location  string      `json:"Location"`
	Position  string      `json:"Position"`
}
//func dataSourceEmployee() *schema.Resource {
//	return &schema.Resource{
//		ReadContext: dataSourceEmployeeRead,
//		Schema: map[string]*schema.Schema{
//			"user_id": &schema.Schema{
//				Type:     schema.TypeInt,
//				Computed: true,
//			},
//			"user_name": &schema.Schema{
//				Type:     schema.TypeString,
//				Computed: true,
//			},
//			"mob_no": &schema.Schema{
//				Type:     schema.TypeInt,
//				Computed: true,
//			},
//			"location": &schema.Schema{
//				Type:     schema.TypeString,
//				Computed: true,
//			},
//			"position": &schema.Schema{
//				Type:     schema.TypeString,
//				Computed: true,
//			},
//		},
//
//	}
//}
func dataSourceEmployee() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEmployeeRead,
		Schema: map[string]*schema.Schema{
			"employee": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"user_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"mob_no": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"location": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"position": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},

	}
}

func dataSourceEmployeeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	//client := &http.Client{Timeout: 10 * time.Second}
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	r, err := http.Get( "http://localhost:8080/")
	if err != nil {
		return diag.FromErr(err)
	}

	//r, err := client.Do(req)
	//if err != nil {
	//	return diag.FromErr(err)
	//}
	log.Printf("Response:%v",r)
	defer r.Body.Close()


	//employee := make([]map[string]interface{}, 0)

	//err = json.NewDecoder(r.Body).Decode(&employee)
	//if err != nil {
	//	return diag.FromErr(err)
	//}
	//log.Printf("Response:%v",employee)
	//
	//if err := d.Set("user_id", employee[2]["UserID"]); err != nil {
	//	return diag.FromErr(err)
	//}
	//if err := d.Set("user_name", employee[2]["Username"]); err != nil {
	//	return diag.FromErr(err)
	//}
	//
	//if err := d.Set("location", employee[2]["Location"]); err != nil {
	//	return diag.FromErr(err)
	//}
	//
	//if err := d.Set("position", employee[2]["Position"]); err != nil {
	//	return diag.FromErr(err)
	//}




		//if err != nil {
		//	return diag.FromErr(err)
		//}
		//log.Printf("Response:%v", response)
		//
		//if err := d.Set("user_id", response[2]["UserID"]); err != nil {
		//	return diag.FromErr(err)
		//}
		//if err := d.Set("user_name", response[2]["Username"]); err != nil {
		//	return diag.FromErr(err)
		//}
		//
		//if err := d.Set("location", response[2]["Location"]); err != nil {
		//	return diag.FromErr(err)
		//}
		//
		//if err := d.Set("position", response[2]["Position"]); err != nil {
		//	return diag.FromErr(err)
		//}




	responseData, err  := ioutil.ReadAll(r.Body)
	if err != nil{
		log.Fatal(err)
	}


	var responseObject Response
	json.Unmarshal(responseData, &responseObject)


	var s []map[string]interface{}
	for _, rg := range responseObject{
		mapping := map[string]interface{}{

			"user_id":  rg.UserID,
			"user_name":  rg.Username,
			"location":  rg.Location,
			"position":  rg.Position,
		}

		s = append(s,mapping)
		log.Printf("Response:%v",s)
	}

	if err := d.Set("employee", s); err != nil{
		return nil
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
func writeToFile(filePath string, data interface{}) error {
	var out string
	switch data.(type) {
	case string:
		out = data.(string)
		break
	case nil:
		return nil
	default:
		bs, err := json.MarshalIndent(data, "", "\t")
		if err != nil {
			return fmt.Errorf("MarshalIndent data %#v got an error: %#v", data, err)
		}
		out = string(bs)
	}
	if strings.HasPrefix(filePath, "~") {
		home, err := GetUserHomeDir()
		if err != nil {
			return err
		}
		if home != "" {
			filePath = strings.Replace(filePath, "~", home, 1)
		}
	}
	if _, err := os.Stat(filePath); err == nil {
		if err := os.Remove(filePath); err != nil {
			return err
		}
	}
	return ioutil.WriteFile(filePath, []byte(out), 422)
}
func GetUserHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("Get current user got an error: %#v.", err)
	}
	return usr.HomeDir, nil
}



