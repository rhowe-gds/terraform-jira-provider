package main

import (
  "encoding/json"

  "github.com/hashicorp/terraform/helper/schema"
  "github.com/hashicorp/terraform/helper/validation"

  jira "gopkg.in/andygrunwald/go-jira.v1"
)

func resourceProject() *schema.Resource {
  return &schema.Resource{
    Create: resourceProjectCreate,
    Read:   resourceProjectRead,
    Update: resourceProjectUpdate,
    Delete: resourceProjectDelete,

    Schema: map[string]*schema.Schema{
      "assignee_type": &schema.Schema{
        Type:         schema.TypeString,
        Optional:     true,
        Default:      "UNASSIGNED",
        ValidateFunc: validation.StringInSlice([]string{"PROJECT_LEAD", "UNASSIGNED"}, false),
      },
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
      "url": &schema.Schema{
        Type:     schema.TypeString,
        Optional: true,
      },
    },
  }
}

type Project struct {
  AssigneeType string `json:"assigneeType"`
  Description  string `json:"description"`
  Email        string `json:"email"`
  Key          string `json:"key"`
  Name         string `json:"name"`
  URL          string `json:"url"`
}

func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
  key := d.Get("key").(string)
  d.SetId(key)

  client := m.(*jira.Client)

  new_project := Project{
    AssigneeType: d.Get("assignee_type").(string),
    Description: d.Get("description").(string),
    Email: d.Get("email").(string),
    Key: key,
    Name: d.Get("name").(string),
    URL: d.Get("url").(string),
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

  d.Set("assignee_type", project.AssigneeType)
  d.Set("description", project.Description)
  d.Set("email", project.Email)
  d.Set("key", project.Key)
  d.Set("name", project.Name)
  d.Set("url", project.URL)

  return nil
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
  return nil
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
  return nil
}
