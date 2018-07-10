package main

import (
  "log"

  "github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
  return &schema.Provider{
    Schema: map[string]*schema.Schema{
      "jira_auth_method": &schema.Schema{
        Type:        schema.TypeString,
        Required:    true,
        DefaultFunc: schema.EnvDefaultFunc("JIRA_AUTH_METHOD", nil),
      },
      "jira_base_url": &schema.Schema{
        Type:        schema.TypeString,
        Required:    true,
        DefaultFunc: schema.EnvDefaultFunc("JIRA_BASE_URL", nil),
      },
      "jira_oauth_consumer_key": &schema.Schema{
        Type:        schema.TypeString,
        Optional:    true,
        DefaultFunc: schema.EnvDefaultFunc("JIRA_OAUTH_CONSUMER_KEY", nil),
      },
      "jira_oauth_private_key": &schema.Schema{
        Type:        schema.TypeString,
        Optional:    true,
        DefaultFunc: schema.EnvDefaultFunc("JIRA_OAUTH_PRIVATE_KEY", nil),
      },
      "jira_username": &schema.Schema{
        Type:        schema.TypeString,
        Required:    true,
        DefaultFunc: schema.EnvDefaultFunc("JIRA_USERNAME", nil),
      },
      "jira_password": &schema.Schema{
        Type:        schema.TypeString,
        Required:    true,
        DefaultFunc: schema.EnvDefaultFunc("JIRA_PASSWORD", nil),
      },
    },
    ResourcesMap: map[string]*schema.Resource{
      "jira_project": resourceProject(),
    },
    ConfigureFunc: providerConfigure,
  }
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
  config := Config{
    AuthMethod: d.Get("jira_auth_method").(string),
    BaseUrl: d.Get("jira_base_url").(string),
    Username: d.Get("jira_username").(string),
    Password: d.Get("jira_password").(string),
    OauthConsumerKey: d.Get("jira_oauth_consumer_key").(string),
  }

  log.Println("[INFO] Initializing Jira client")
  client := config.Client()

  return client, nil
}
