package main

import (
  "encoding/json"

  "github.com/hashicorp/terraform/helper/schema"

  jira "gopkg.in/andygrunwald/go-jira.v1"
)

func resourceProject() *schema.Resource {
  return &schema.Resource{
    Create: resourceProjectCreate,
    Read:   resourceProjectRead,
    Update: resourceProjectUpdate,
    Delete: resourceProjectDelete,

    Schema: map[string]*schema.Schema{
      "description": &schema.Schema{
        Type:     schema.TypeString,
        Optional: true,
      },
      "email": &schema.Schema{
        Type:     schema.TypeString,
        Optional: true,
      },
      "key": &schema.Schema{
        Type:     schema.TypeString,
        Required: true,
      },
      "name": &schema.Schema{
        Type:     schema.TypeString,
        Required: true,
      },
    },
  }
}

type Project struct {
  Description string `json:"description"`
  Email       string `json:"email"`
  Key         string `json:"key"`
  Name        string `json:"name"`
}

func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
  key := d.Get("key").(string)
  d.SetId(key)

  client := m.(*jira.Client)

  new_project := Project{
    Description: d.Get("description").(string),
    Email: d.Get("email").(string),
    Key: key,
    Name: d.Get("name").(string),
  }

  // Project creation is not implemented by go-jira
  j, err := json.Marshal(new_project)

  if err != nil {
    return err
  }

  client.NewRequest("POST", "rest/api/2/project", j)

  return nil
}

func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
  client := m.(*jira.Client)

  project, resp, err := client.Project.Get(d.Id())

  if resp.StatusCode == 404 {
    d.SetId("")
    return nil
  }

  if err != nil {
    return err
  }

  d.Set("description", project.Description)
  d.Set("email", project.Email)
  d.Set("key", project.Key)
  d.Set("name", project.Name)

  return nil
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
  return nil
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
  return nil
}
