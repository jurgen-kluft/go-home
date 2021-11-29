// To parse and unparse this JSON data, add this code to your project and do:
//
//    r, err := UnmarshalCcalendar(bytes)
//    bytes, err = r.Marshal()

//  {
// 	"name": "Home",
//  "url": "http://"
// 	"url": "zXxPr7iy8AU0fwN57TbTW_eIeIDlyGSarB0Efu9en8iR1L0DrXpG8apaR6WldVSo49Zpp8kXo4t25izstjaE1YqU3Scbr9k9hmAeXxCu0AZH1rx5j2f4eW_V-ackWM5CrRm7U-FlgSDzOPVMFfalFQWGdDzrOunfxqnxrZWt77TY"
//
//  "url": "file://MTY5ODMwNTE0MzE2OTgzMCxyjL_SGbifBa-CPGxwhRTU4PcTmof8n8ja2YvzGiPH.ics"
//  "url": "W_ElpkUQPAmXmJB-6sNknBwfVrBLyEeOzMo9KEIrtwEv5Osb4EaBCoMl8YpUFPr56hFCYdzzNHCwlyl8hpRwcWhYNh_f9T34BJs9UK0bqPhOIpdQxArAkrMVoQ=="
//   },
//   {
// 	"name": "Season",
//  "url": "http://"
// 	"url": "UNBDFqQf0AmljLgoiqi__ds5wOChFiUBouzJ0dxLWY0z9PqzVThhDWg8Z7CzRrn8-B5H8SwmLF5Peiq6XJzh-vJ0Z1zoZ4_Z5KBx3wCJ-mXvW1-8Kci2XC4w-9N7nLAOuLWNDwZHQ6HyNUPXEdGS7YWN4Sqfddku7BZm7GzdgQ8VzmxP20rsgXcmjGvPjOaQ4PpaJbGx9SQy67z1glA3boJ9UuDHxYOS"
//
//  "url": "file://MTY5ODMwNTE0MzE2OTgzMCxyjL_SGbifBa-CPGxwhRTF01br2oW796m-fnT_v2B8u4CQY3oi8wRnziliGDSySbKLzVFvtPeFd7FGZE7rEy4.ics"
//  "url": "E5TdywPHr3bj2USXucchGNEFGgn3-FAfMbhTSsXXuwzPgmm8LA5M9VdzHHORQf6oAs_j8fQBBtoaX2Lvk8tS0kO4jJUUvdFlJtjVBN9bZkAIkosH_slcua_A0RSHnHk9pByDH6Eq36uvpyMBzrDYx0thwos69zrfq5wf1DVNWDLWIkd4O6I="
//   },
//   {
// 	"name": "TimeOfDay2",
//  "url": "http://"
// 	"url": "dKJsQjiixHPeLrTVkQPgJsAbT60A5Lf2YqEu1ZkZkOafy2MThqYFk2XahH0jYvG6LfsWTbdmU2Kfvm-jneNuCGLNFw0F8BAfOdlnywzN2P3V-3YgWCqnfGzTDxMQxxbxr4JCbcArf3VG8sKY1j4G2JI5T0BCqxQXaCPn_YtnHI7gq38xnpwlburfzRRLVx_ddyMRtZ62b0GciR248HSTQRNxTm2B"
//
//  "url": "file://MTY5ODMwNTE0MzE2OTgzMCxyjL_SGbifBa-CPGxwhRSI-eyL6c-3ZyJgzA7r_A9e0adqYeMe7Otbm6lrhe05ZfJhPRkgxxVAmmqif0EVL_w.ics"
//  "url": "FwTntbDZWE_B5P6n4XHBEebQSLN6v4WaCckrbNq7QgGyl_e_GPWdqde8FvpHoB4aQLpdnhO5lU2IlAtyPDc8oE23CvQgmGh9hqrBmZK_PCJmofhK1kk0dREEiIG0N09aHYtDA-eSGrG6-RZF-ZUdSDTdwkqCjZj2eqPCxsj7bj0Snb_hrTI="
//   }

package config

import (
	"encoding/json"
	"fmt"
)

func CalendarConfigFromJSON(data []byte) (*CalendarConfig, error) {
	r := &CalendarConfig{}
	err := json.Unmarshal(data, r)
	return r, err
}

// FromJSON converts a json string to a CalendarConfig instance
func (r *CalendarConfig) FromJSON(data []byte) error {
	c := CalendarConfig{}
	err := json.Unmarshal(data, &c)
	*r = c
	return err
}

// ToJSON converts a CalendarConfig to a JSON string
func (r *CalendarConfig) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

type CalendarConfig struct {
	Calendars []Ccal    `json:"calendars"`
	Sensors   []Csensor `json:"sensors"`
	Rules     []Crule   `json:"rules"`
}

func (c *CalendarConfig) Print() {
	for _, cal := range c.Calendars {
		cal.print()
	}
	for _, sensor := range c.Sensors {
		sensor.print()
	}
	for _, rule := range c.Rules {
		rule.print()
	}
}

type Ccal struct {
	Name string      `json:"name"`
	URL  CryptString `json:"url"`
}

func (c Ccal) print() {
	fmt.Printf("calendar.name = %s\n", c.Name)
	fmt.Printf("calendar.url = %s\n", c.URL)
}

type Csensor struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	AttrType  string `json:"attrtype"`
	AttrValue string `json:"attrvalue"`
}

func (c Csensor) print() {
	fmt.Printf("sensor.name = %s\n", c.Name)
	fmt.Printf("sensor.type = %s\n", c.Type)
	fmt.Printf("sensor.attrvalue = %s\n", c.AttrValue)
	fmt.Printf("sensor.attrtype = %s\n", c.AttrType)
}

type Crule struct {
	Key    string   `json:"key"`
	Value  string   `json:"value"`
	IfThen []IfThen `json:"if"`
}

type IfThen struct {
	Key   string `json:"key"`
	State string `json:"state"`
}

func (c Crule) print() {
	fmt.Printf("rule.key = %s\n", c.Key)
	fmt.Printf("rule.value = %s\n", c.Value)
	fmt.Println("If:")
	for _, rule := range c.IfThen {
		fmt.Printf("    if key='%s' has state:'%s'\n", rule.Key, rule.State)
	}
}
