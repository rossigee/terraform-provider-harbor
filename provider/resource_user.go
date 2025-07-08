package provider

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"full_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"admin": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := client.UserBody(d)

	_, header, _, err := apiClient.SendRequest("POST", models.PathUsers, &body, 201)
	if err != nil {
		return err
	}

	id, err := client.GetID(header)
	if err != nil {
		return nil
	}

	d.SetId(id)
	return resourceUserRead(d, m)
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	
	currentId := d.Id()
	if currentId == "" {
		return nil
	}
	
	// Check if the ID is in the expected numeric format (/users/123)
	// If not, it's likely a username from Crossplane's external-name
	isNumericId := regexp.MustCompile(`^/users/\d+$`).MatchString(currentId)
	
	// If it's not a numeric ID, we need to search for the user by username
	if !isNumericId {
		// Extract username from the ID (remove /users/ prefix if present)
		searchUsername := currentId
		if strings.HasPrefix(searchUsername, "/users/") {
			searchUsername = strings.TrimPrefix(searchUsername, "/users/")
		}
		
		// Search for the user by username
		allUsersResp, _, _, err := apiClient.SendRequest("GET", models.PathUsers, nil, 200)
		if err != nil {
			return fmt.Errorf("Failed to search for users: %v", err)
		}
		
		var allUsers []models.UserBody
		err = json.Unmarshal([]byte(allUsersResp), &allUsers)
		if err != nil {
			return fmt.Errorf("Failed to parse users list: %v", err)
		}
		
		// Look for user with matching username
		var targetUser *models.UserBody
		for _, user := range allUsers {
			if user.Username == searchUsername {
				targetUser = &user
				break
			}
		}
		
		if targetUser == nil {
			// User doesn't exist - clear the ID to trigger creation
			d.SetId("")
			return nil
		}
		
		// Update the ID to use the numeric user ID for future operations
		userIdPath := fmt.Sprintf("/users/%d", targetUser.UserID)
		d.SetId(userIdPath)
		
		// Set the user data from the found user
		d.Set("username", targetUser.Username)
		d.Set("full_name", targetUser.Realname)
		d.Set("email", targetUser.Email)
		d.Set("admin", targetUser.SysadminFlag)
		d.Set("comment", targetUser.Comment)
		
		return nil
	}
	
	// ID is already in numeric format, proceed with normal read
	resp, _, _, err := apiClient.SendRequest("GET", currentId, nil, 200)
	if err != nil {
		return err
	}
	
	var jsonData models.UserBody
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return fmt.Errorf("Resource not found %s", d.Id())
	}

	d.Set("username", jsonData.Username)
	d.Set("full_name", jsonData.Realname)
	d.Set("email", jsonData.Email)
	d.Set("admin", jsonData.SysadminFlag)
	d.Set("comment", jsonData.Comment)

	return nil
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := client.UserBody(d)
	_, _, _, err := apiClient.SendRequest("PUT", d.Id(), body, 200)
	if err != nil {
		return err
	}

	_, _, _, err = apiClient.SendRequest("PUT", d.Id()+"/sysadmin", body, 200)
	if err != nil {
		return err
	}

	if d.HasChange("password") == true {
		_, _, _, err = apiClient.SendRequest("PUT", d.Id()+"/password", body, 200)
		if err != nil {
			return err
		}
	}

	return resourceUserRead(d, m)
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	_, _, _, err := apiClient.SendRequest("DELETE", d.Id(), nil, 200)
	if err != nil {
		return err
	}
	return nil
}
