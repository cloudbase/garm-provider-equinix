/*
Metal API

# Introduction Equinix Metal provides a RESTful HTTP API which can be reached at <https://api.equinix.com/metal/v1>. This document describes the API and how to use it.  The API allows you to programmatically interact with all of your Equinix Metal resources, including devices, networks, addresses, organizations, projects, and your user account. Every feature of the Equinix Metal web interface is accessible through the API.  The API docs are generated from the Equinix Metal OpenAPI specification and are officially hosted at <https://metal.equinix.com/developers/api>.  # Common Parameters  The Equinix Metal API uses a few methods to minimize network traffic and improve throughput. These parameters are not used in all API calls, but are used often enough to warrant their own section. Look for these parameters in the documentation for the API calls that support them.  ## Pagination  Pagination is used to limit the number of results returned in a single request. The API will return a maximum of 100 results per page. To retrieve additional results, you can use the `page` and `per_page` query parameters.  The `page` parameter is used to specify the page number. The first page is `1`. The `per_page` parameter is used to specify the number of results per page. The maximum number of results differs by resource type.  ## Sorting  Where offered, the API allows you to sort results by a specific field. To sort results use the `sort_by` query parameter with the root level field name as the value. The `sort_direction` parameter is used to specify the sort direction, either either `asc` (ascending) or `desc` (descending).  ## Filtering  Filtering is used to limit the results returned in a single request. The API supports filtering by certain fields in the response. To filter results, you can use the field as a query parameter.  For example, to filter the IP list to only return public IPv4 addresses, you can filter by the `type` field, as in the following request:  ```sh curl -H 'X-Auth-Token: my_authentication_token' \\   https://api.equinix.com/metal/v1/projects/id/ips?type=public_ipv4 ```  Only IP addresses with the `type` field set to `public_ipv4` will be returned.  ## Searching  Searching is used to find matching resources using multiple field comparissons. The API supports searching in resources that define this behavior. Currently the search parameter is only available on devices, ssh_keys, api_keys and memberships endpoints.  To search resources you can use the `search` query parameter.  ## Include and Exclude  For resources that contain references to other resources, sucha as a Device that refers to the Project it resides in, the Equinix Metal API will returns `href` values (API links) to the associated resource.  ```json {   ...   \"project\": {     \"href\": \"/metal/v1/projects/f3f131c8-f302-49ef-8c44-9405022dc6dd\"   } } ```  If you're going need the project details, you can avoid a second API request.  Specify the contained `href` resources and collections that you'd like to have included in the response using the `include` query parameter.  For example:  ```sh curl -H 'X-Auth-Token: my_authentication_token' \\   https://api.equinix.com/metal/v1/user?include=projects ```  The `include` parameter is generally accepted in `GET`, `POST`, `PUT`, and `PATCH` requests where `href` resources are presented.  To have multiple resources include, use a comma-separated list (e.g. `?include=emails,projects,memberships`).  ```sh curl -H 'X-Auth-Token: my_authentication_token' \\   https://api.equinix.com/metal/v1/user?include=emails,projects,memberships ```  You may also include nested associations up to three levels deep using dot notation (`?include=memberships.projects`):  ```sh curl -H 'X-Auth-Token: my_authentication_token' \\   https://api.equinix.com/metal/v1/user?include=memberships.projects ```  To exclude resources, and optimize response delivery, use the `exclude` query parameter. The `exclude` parameter is generally accepted in `GET`, `POST`, `PUT`, and `PATCH` requests for fields with nested object responses. When excluded, these fields will be replaced with an object that contains only an `href` field.

API version: 1.0.0
Contact: support@equinixmetal.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package metalv1

import (
	"encoding/json"
	"time"
)

// checks if the SpotMarketRequestCreateInputInstanceParameters type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &SpotMarketRequestCreateInputInstanceParameters{}

// SpotMarketRequestCreateInputInstanceParameters struct for SpotMarketRequestCreateInputInstanceParameters
type SpotMarketRequestCreateInputInstanceParameters struct {
	AlwaysPxe    *bool                  `json:"always_pxe,omitempty"`
	BillingCycle *string                `json:"billing_cycle,omitempty"`
	Customdata   map[string]interface{} `json:"customdata,omitempty"`
	Description  *string                `json:"description,omitempty"`
	Features     []string               `json:"features,omitempty"`
	Hostname     *string                `json:"hostname,omitempty"`
	Hostnames    []string               `json:"hostnames,omitempty"`
	// Whether the device should be locked, preventing accidental deletion.
	Locked                *bool      `json:"locked,omitempty"`
	NoSshKeys             *bool      `json:"no_ssh_keys,omitempty"`
	OperatingSystem       *string    `json:"operating_system,omitempty"`
	Plan                  *string    `json:"plan,omitempty"`
	PrivateIpv4SubnetSize *int32     `json:"private_ipv4_subnet_size,omitempty"`
	ProjectSshKeys        []string   `json:"project_ssh_keys,omitempty"`
	PublicIpv4SubnetSize  *int32     `json:"public_ipv4_subnet_size,omitempty"`
	Tags                  []string   `json:"tags,omitempty"`
	TerminationTime       *time.Time `json:"termination_time,omitempty"`
	// The UUIDs of users whose SSH keys should be included on the provisioned device.
	UserSshKeys          []string `json:"user_ssh_keys,omitempty"`
	Userdata             *string  `json:"userdata,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _SpotMarketRequestCreateInputInstanceParameters SpotMarketRequestCreateInputInstanceParameters

// NewSpotMarketRequestCreateInputInstanceParameters instantiates a new SpotMarketRequestCreateInputInstanceParameters object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSpotMarketRequestCreateInputInstanceParameters() *SpotMarketRequestCreateInputInstanceParameters {
	this := SpotMarketRequestCreateInputInstanceParameters{}
	return &this
}

// NewSpotMarketRequestCreateInputInstanceParametersWithDefaults instantiates a new SpotMarketRequestCreateInputInstanceParameters object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSpotMarketRequestCreateInputInstanceParametersWithDefaults() *SpotMarketRequestCreateInputInstanceParameters {
	this := SpotMarketRequestCreateInputInstanceParameters{}
	return &this
}

// GetAlwaysPxe returns the AlwaysPxe field value if set, zero value otherwise.
func (o *SpotMarketRequestCreateInputInstanceParameters) GetAlwaysPxe() bool {
	if o == nil || IsNil(o.AlwaysPxe) {
		var ret bool
		return ret
	}
	return *o.AlwaysPxe
}

// GetAlwaysPxeOk returns a tuple with the AlwaysPxe field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SpotMarketRequestCreateInputInstanceParameters) GetAlwaysPxeOk() (*bool, bool) {
	if o == nil || IsNil(o.AlwaysPxe) {
		return nil, false
	}
	return o.AlwaysPxe, true
}

// HasAlwaysPxe returns a boolean if a field has been set.
func (o *SpotMarketRequestCreateInputInstanceParameters) HasAlwaysPxe() bool {
	if o != nil && !IsNil(o.AlwaysPxe) {
		return true
	}

	return false
}

// SetAlwaysPxe gets a reference to the given bool and assigns it to the AlwaysPxe field.
func (o *SpotMarketRequestCreateInputInstanceParameters) SetAlwaysPxe(v bool) {
	o.AlwaysPxe = &v
}

// GetBillingCycle returns the BillingCycle field value if set, zero value otherwise.
func (o *SpotMarketRequestCreateInputInstanceParameters) GetBillingCycle() string {
	if o == nil || IsNil(o.BillingCycle) {
		var ret string
		return ret
	}
	return *o.BillingCycle
}

// GetBillingCycleOk returns a tuple with the BillingCycle field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SpotMarketRequestCreateInputInstanceParameters) GetBillingCycleOk() (*string, bool) {
	if o == nil || IsNil(o.BillingCycle) {
		return nil, false
	}
	return o.BillingCycle, true
}

// HasBillingCycle returns a boolean if a field has been set.
func (o *SpotMarketRequestCreateInputInstanceParameters) HasBillingCycle() bool {
	if o != nil && !IsNil(o.BillingCycle) {
		return true
	}

	return false
}

// SetBillingCycle gets a reference to the given string and assigns it to the BillingCycle field.
func (o *SpotMarketRequestCreateInputInstanceParameters) SetBillingCycle(v string) {
	o.BillingCycle = &v
}

// GetCustomdata returns the Customdata field value if set, zero value otherwise.
func (o *SpotMarketRequestCreateInputInstanceParameters) GetCustomdata() map[string]interface{} {
	if o == nil || IsNil(o.Customdata) {
		var ret map[string]interface{}
		return ret
	}
	return o.Customdata
}

// GetCustomdataOk returns a tuple with the Customdata field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SpotMarketRequestCreateInputInstanceParameters) GetCustomdataOk() (map[string]interface{}, bool) {
	if o == nil || IsNil(o.Customdata) {
		return map[string]interface{}{}, false
	}
	return o.Customdata, true
}

// HasCustomdata returns a boolean if a field has been set.
func (o *SpotMarketRequestCreateInputInstanceParameters) HasCustomdata() bool {
	if o != nil && !IsNil(o.Customdata) {
		return true
	}

	return false
}

// SetCustomdata gets a reference to the given map[string]interface{} and assigns it to the Customdata field.
func (o *SpotMarketRequestCreateInputInstanceParameters) SetCustomdata(v map[string]interface{}) {
	o.Customdata = v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *SpotMarketRequestCreateInputInstanceParameters) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SpotMarketRequestCreateInputInstanceParameters) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *SpotMarketRequestCreateInputInstanceParameters) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *SpotMarketRequestCreateInputInstanceParameters) SetDescription(v string) {
	o.Description = &v
}

// GetFeatures returns the Features field value if set, zero value otherwise.
func (o *SpotMarketRequestCreateInputInstanceParameters) GetFeatures() []string {
	if o == nil || IsNil(o.Features) {
		var ret []string
		return ret
	}
	return o.Features
}

// GetFeaturesOk returns a tuple with the Features field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SpotMarketRequestCreateInputInstanceParameters) GetFeaturesOk() ([]string, bool) {
	if o == nil || IsNil(o.Features) {
		return nil, false
	}
	return o.Features, true
}

// HasFeatures returns a boolean if a field has been set.
func (o *SpotMarketRequestCreateInputInstanceParameters) HasFeatures() bool {
	if o != nil && !IsNil(o.Features) {
		return true
	}

	return false
}

// SetFeatures gets a reference to the given []string and assigns it to the Features field.
func (o *SpotMarketRequestCreateInputInstanceParameters) SetFeatures(v []string) {
	o.Features = v
}

// GetHostname returns the Hostname field value if set, zero value otherwise.
func (o *SpotMarketRequestCreateInputInstanceParameters) GetHostname() string {
	if o == nil || IsNil(o.Hostname) {
		var ret string
		return ret
	}
	return *o.Hostname
}

// GetHostnameOk returns a tuple with the Hostname field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SpotMarketRequestCreateInputInstanceParameters) GetHostnameOk() (*string, bool) {
	if o == nil || IsNil(o.Hostname) {
		return nil, false
	}
	return o.Hostname, true
}

// HasHostname returns a boolean if a field has been set.
func (o *SpotMarketRequestCreateInputInstanceParameters) HasHostname() bool {
	if o != nil && !IsNil(o.Hostname) {
		return true
	}

	return false
}

// SetHostname gets a reference to the given string and assigns it to the Hostname field.
func (o *SpotMarketRequestCreateInputInstanceParameters) SetHostname(v string) {
	o.Hostname = &v
}

// GetHostnames returns the Hostnames field value if set, zero value otherwise.
func (o *SpotMarketRequestCreateInputInstanceParameters) GetHostnames() []string {
	if o == nil || IsNil(o.Hostnames) {
		var ret []string
		return ret
	}
	return o.Hostnames
}

// GetHostnamesOk returns a tuple with the Hostnames field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SpotMarketRequestCreateInputInstanceParameters) GetHostnamesOk() ([]string, bool) {
	if o == nil || IsNil(o.Hostnames) {
		return nil, false
	}
	return o.Hostnames, true
}

// HasHostnames returns a boolean if a field has been set.
func (o *SpotMarketRequestCreateInputInstanceParameters) HasHostnames() bool {
	if o != nil && !IsNil(o.Hostnames) {
		return true
	}

	return false
}

// SetHostnames gets a reference to the given []string and assigns it to the Hostnames field.
func (o *SpotMarketRequestCreateInputInstanceParameters) SetHostnames(v []string) {
	o.Hostnames = v
}

// GetLocked returns the Locked field value if set, zero value otherwise.
func (o *SpotMarketRequestCreateInputInstanceParameters) GetLocked() bool {
	if o == nil || IsNil(o.Locked) {
		var ret bool
		return ret
	}
	return *o.Locked
}

// GetLockedOk returns a tuple with the Locked field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SpotMarketRequestCreateInputInstanceParameters) GetLockedOk() (*bool, bool) {
	if o == nil || IsNil(o.Locked) {
		return nil, false
	}
	return o.Locked, true
}

// HasLocked returns a boolean if a field has been set.
func (o *SpotMarketRequestCreateInputInstanceParameters) HasLocked() bool {
	if o != nil && !IsNil(o.Locked) {
		return true
	}

	return false
}

// SetLocked gets a reference to the given bool and assigns it to the Locked field.
func (o *SpotMarketRequestCreateInputInstanceParameters) SetLocked(v bool) {
	o.Locked = &v
}

// GetNoSshKeys returns the NoSshKeys field value if set, zero value otherwise.
func (o *SpotMarketRequestCreateInputInstanceParameters) GetNoSshKeys() bool {
	if o == nil || IsNil(o.NoSshKeys) {
		var ret bool
		return ret
	}
	return *o.NoSshKeys
}

// GetNoSshKeysOk returns a tuple with the NoSshKeys field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SpotMarketRequestCreateInputInstanceParameters) GetNoSshKeysOk() (*bool, bool) {
	if o == nil || IsNil(o.NoSshKeys) {
		return nil, false
	}
	return o.NoSshKeys, true
}

// HasNoSshKeys returns a boolean if a field has been set.
func (o *SpotMarketRequestCreateInputInstanceParameters) HasNoSshKeys() bool {
	if o != nil && !IsNil(o.NoSshKeys) {
		return true
	}

	return false
}

// SetNoSshKeys gets a reference to the given bool and assigns it to the NoSshKeys field.
func (o *SpotMarketRequestCreateInputInstanceParameters) SetNoSshKeys(v bool) {
	o.NoSshKeys = &v
}

// GetOperatingSystem returns the OperatingSystem field value if set, zero value otherwise.
func (o *SpotMarketRequestCreateInputInstanceParameters) GetOperatingSystem() string {
	if o == nil || IsNil(o.OperatingSystem) {
		var ret string
		return ret
	}
	return *o.OperatingSystem
}

// GetOperatingSystemOk returns a tuple with the OperatingSystem field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SpotMarketRequestCreateInputInstanceParameters) GetOperatingSystemOk() (*string, bool) {
	if o == nil || IsNil(o.OperatingSystem) {
		return nil, false
	}
	return o.OperatingSystem, true
}

// HasOperatingSystem returns a boolean if a field has been set.
func (o *SpotMarketRequestCreateInputInstanceParameters) HasOperatingSystem() bool {
	if o != nil && !IsNil(o.OperatingSystem) {
		return true
	}

	return false
}

// SetOperatingSystem gets a reference to the given string and assigns it to the OperatingSystem field.
func (o *SpotMarketRequestCreateInputInstanceParameters) SetOperatingSystem(v string) {
	o.OperatingSystem = &v
}

// GetPlan returns the Plan field value if set, zero value otherwise.
func (o *SpotMarketRequestCreateInputInstanceParameters) GetPlan() string {
	if o == nil || IsNil(o.Plan) {
		var ret string
		return ret
	}
	return *o.Plan
}

// GetPlanOk returns a tuple with the Plan field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SpotMarketRequestCreateInputInstanceParameters) GetPlanOk() (*string, bool) {
	if o == nil || IsNil(o.Plan) {
		return nil, false
	}
	return o.Plan, true
}

// HasPlan returns a boolean if a field has been set.
func (o *SpotMarketRequestCreateInputInstanceParameters) HasPlan() bool {
	if o != nil && !IsNil(o.Plan) {
		return true
	}

	return false
}

// SetPlan gets a reference to the given string and assigns it to the Plan field.
func (o *SpotMarketRequestCreateInputInstanceParameters) SetPlan(v string) {
	o.Plan = &v
}

// GetPrivateIpv4SubnetSize returns the PrivateIpv4SubnetSize field value if set, zero value otherwise.
func (o *SpotMarketRequestCreateInputInstanceParameters) GetPrivateIpv4SubnetSize() int32 {
	if o == nil || IsNil(o.PrivateIpv4SubnetSize) {
		var ret int32
		return ret
	}
	return *o.PrivateIpv4SubnetSize
}

// GetPrivateIpv4SubnetSizeOk returns a tuple with the PrivateIpv4SubnetSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SpotMarketRequestCreateInputInstanceParameters) GetPrivateIpv4SubnetSizeOk() (*int32, bool) {
	if o == nil || IsNil(o.PrivateIpv4SubnetSize) {
		return nil, false
	}
	return o.PrivateIpv4SubnetSize, true
}

// HasPrivateIpv4SubnetSize returns a boolean if a field has been set.
func (o *SpotMarketRequestCreateInputInstanceParameters) HasPrivateIpv4SubnetSize() bool {
	if o != nil && !IsNil(o.PrivateIpv4SubnetSize) {
		return true
	}

	return false
}

// SetPrivateIpv4SubnetSize gets a reference to the given int32 and assigns it to the PrivateIpv4SubnetSize field.
func (o *SpotMarketRequestCreateInputInstanceParameters) SetPrivateIpv4SubnetSize(v int32) {
	o.PrivateIpv4SubnetSize = &v
}

// GetProjectSshKeys returns the ProjectSshKeys field value if set, zero value otherwise.
func (o *SpotMarketRequestCreateInputInstanceParameters) GetProjectSshKeys() []string {
	if o == nil || IsNil(o.ProjectSshKeys) {
		var ret []string
		return ret
	}
	return o.ProjectSshKeys
}

// GetProjectSshKeysOk returns a tuple with the ProjectSshKeys field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SpotMarketRequestCreateInputInstanceParameters) GetProjectSshKeysOk() ([]string, bool) {
	if o == nil || IsNil(o.ProjectSshKeys) {
		return nil, false
	}
	return o.ProjectSshKeys, true
}

// HasProjectSshKeys returns a boolean if a field has been set.
func (o *SpotMarketRequestCreateInputInstanceParameters) HasProjectSshKeys() bool {
	if o != nil && !IsNil(o.ProjectSshKeys) {
		return true
	}

	return false
}

// SetProjectSshKeys gets a reference to the given []string and assigns it to the ProjectSshKeys field.
func (o *SpotMarketRequestCreateInputInstanceParameters) SetProjectSshKeys(v []string) {
	o.ProjectSshKeys = v
}

// GetPublicIpv4SubnetSize returns the PublicIpv4SubnetSize field value if set, zero value otherwise.
func (o *SpotMarketRequestCreateInputInstanceParameters) GetPublicIpv4SubnetSize() int32 {
	if o == nil || IsNil(o.PublicIpv4SubnetSize) {
		var ret int32
		return ret
	}
	return *o.PublicIpv4SubnetSize
}

// GetPublicIpv4SubnetSizeOk returns a tuple with the PublicIpv4SubnetSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SpotMarketRequestCreateInputInstanceParameters) GetPublicIpv4SubnetSizeOk() (*int32, bool) {
	if o == nil || IsNil(o.PublicIpv4SubnetSize) {
		return nil, false
	}
	return o.PublicIpv4SubnetSize, true
}

// HasPublicIpv4SubnetSize returns a boolean if a field has been set.
func (o *SpotMarketRequestCreateInputInstanceParameters) HasPublicIpv4SubnetSize() bool {
	if o != nil && !IsNil(o.PublicIpv4SubnetSize) {
		return true
	}

	return false
}

// SetPublicIpv4SubnetSize gets a reference to the given int32 and assigns it to the PublicIpv4SubnetSize field.
func (o *SpotMarketRequestCreateInputInstanceParameters) SetPublicIpv4SubnetSize(v int32) {
	o.PublicIpv4SubnetSize = &v
}

// GetTags returns the Tags field value if set, zero value otherwise.
func (o *SpotMarketRequestCreateInputInstanceParameters) GetTags() []string {
	if o == nil || IsNil(o.Tags) {
		var ret []string
		return ret
	}
	return o.Tags
}

// GetTagsOk returns a tuple with the Tags field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SpotMarketRequestCreateInputInstanceParameters) GetTagsOk() ([]string, bool) {
	if o == nil || IsNil(o.Tags) {
		return nil, false
	}
	return o.Tags, true
}

// HasTags returns a boolean if a field has been set.
func (o *SpotMarketRequestCreateInputInstanceParameters) HasTags() bool {
	if o != nil && !IsNil(o.Tags) {
		return true
	}

	return false
}

// SetTags gets a reference to the given []string and assigns it to the Tags field.
func (o *SpotMarketRequestCreateInputInstanceParameters) SetTags(v []string) {
	o.Tags = v
}

// GetTerminationTime returns the TerminationTime field value if set, zero value otherwise.
func (o *SpotMarketRequestCreateInputInstanceParameters) GetTerminationTime() time.Time {
	if o == nil || IsNil(o.TerminationTime) {
		var ret time.Time
		return ret
	}
	return *o.TerminationTime
}

// GetTerminationTimeOk returns a tuple with the TerminationTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SpotMarketRequestCreateInputInstanceParameters) GetTerminationTimeOk() (*time.Time, bool) {
	if o == nil || IsNil(o.TerminationTime) {
		return nil, false
	}
	return o.TerminationTime, true
}

// HasTerminationTime returns a boolean if a field has been set.
func (o *SpotMarketRequestCreateInputInstanceParameters) HasTerminationTime() bool {
	if o != nil && !IsNil(o.TerminationTime) {
		return true
	}

	return false
}

// SetTerminationTime gets a reference to the given time.Time and assigns it to the TerminationTime field.
func (o *SpotMarketRequestCreateInputInstanceParameters) SetTerminationTime(v time.Time) {
	o.TerminationTime = &v
}

// GetUserSshKeys returns the UserSshKeys field value if set, zero value otherwise.
func (o *SpotMarketRequestCreateInputInstanceParameters) GetUserSshKeys() []string {
	if o == nil || IsNil(o.UserSshKeys) {
		var ret []string
		return ret
	}
	return o.UserSshKeys
}

// GetUserSshKeysOk returns a tuple with the UserSshKeys field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SpotMarketRequestCreateInputInstanceParameters) GetUserSshKeysOk() ([]string, bool) {
	if o == nil || IsNil(o.UserSshKeys) {
		return nil, false
	}
	return o.UserSshKeys, true
}

// HasUserSshKeys returns a boolean if a field has been set.
func (o *SpotMarketRequestCreateInputInstanceParameters) HasUserSshKeys() bool {
	if o != nil && !IsNil(o.UserSshKeys) {
		return true
	}

	return false
}

// SetUserSshKeys gets a reference to the given []string and assigns it to the UserSshKeys field.
func (o *SpotMarketRequestCreateInputInstanceParameters) SetUserSshKeys(v []string) {
	o.UserSshKeys = v
}

// GetUserdata returns the Userdata field value if set, zero value otherwise.
func (o *SpotMarketRequestCreateInputInstanceParameters) GetUserdata() string {
	if o == nil || IsNil(o.Userdata) {
		var ret string
		return ret
	}
	return *o.Userdata
}

// GetUserdataOk returns a tuple with the Userdata field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SpotMarketRequestCreateInputInstanceParameters) GetUserdataOk() (*string, bool) {
	if o == nil || IsNil(o.Userdata) {
		return nil, false
	}
	return o.Userdata, true
}

// HasUserdata returns a boolean if a field has been set.
func (o *SpotMarketRequestCreateInputInstanceParameters) HasUserdata() bool {
	if o != nil && !IsNil(o.Userdata) {
		return true
	}

	return false
}

// SetUserdata gets a reference to the given string and assigns it to the Userdata field.
func (o *SpotMarketRequestCreateInputInstanceParameters) SetUserdata(v string) {
	o.Userdata = &v
}

func (o SpotMarketRequestCreateInputInstanceParameters) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o SpotMarketRequestCreateInputInstanceParameters) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.AlwaysPxe) {
		toSerialize["always_pxe"] = o.AlwaysPxe
	}
	if !IsNil(o.BillingCycle) {
		toSerialize["billing_cycle"] = o.BillingCycle
	}
	if !IsNil(o.Customdata) {
		toSerialize["customdata"] = o.Customdata
	}
	if !IsNil(o.Description) {
		toSerialize["description"] = o.Description
	}
	if !IsNil(o.Features) {
		toSerialize["features"] = o.Features
	}
	if !IsNil(o.Hostname) {
		toSerialize["hostname"] = o.Hostname
	}
	if !IsNil(o.Hostnames) {
		toSerialize["hostnames"] = o.Hostnames
	}
	if !IsNil(o.Locked) {
		toSerialize["locked"] = o.Locked
	}
	if !IsNil(o.NoSshKeys) {
		toSerialize["no_ssh_keys"] = o.NoSshKeys
	}
	if !IsNil(o.OperatingSystem) {
		toSerialize["operating_system"] = o.OperatingSystem
	}
	if !IsNil(o.Plan) {
		toSerialize["plan"] = o.Plan
	}
	if !IsNil(o.PrivateIpv4SubnetSize) {
		toSerialize["private_ipv4_subnet_size"] = o.PrivateIpv4SubnetSize
	}
	if !IsNil(o.ProjectSshKeys) {
		toSerialize["project_ssh_keys"] = o.ProjectSshKeys
	}
	if !IsNil(o.PublicIpv4SubnetSize) {
		toSerialize["public_ipv4_subnet_size"] = o.PublicIpv4SubnetSize
	}
	if !IsNil(o.Tags) {
		toSerialize["tags"] = o.Tags
	}
	if !IsNil(o.TerminationTime) {
		toSerialize["termination_time"] = o.TerminationTime
	}
	if !IsNil(o.UserSshKeys) {
		toSerialize["user_ssh_keys"] = o.UserSshKeys
	}
	if !IsNil(o.Userdata) {
		toSerialize["userdata"] = o.Userdata
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *SpotMarketRequestCreateInputInstanceParameters) UnmarshalJSON(bytes []byte) (err error) {
	varSpotMarketRequestCreateInputInstanceParameters := _SpotMarketRequestCreateInputInstanceParameters{}

	err = json.Unmarshal(bytes, &varSpotMarketRequestCreateInputInstanceParameters)

	if err != nil {
		return err
	}

	*o = SpotMarketRequestCreateInputInstanceParameters(varSpotMarketRequestCreateInputInstanceParameters)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(bytes, &additionalProperties); err == nil {
		delete(additionalProperties, "always_pxe")
		delete(additionalProperties, "billing_cycle")
		delete(additionalProperties, "customdata")
		delete(additionalProperties, "description")
		delete(additionalProperties, "features")
		delete(additionalProperties, "hostname")
		delete(additionalProperties, "hostnames")
		delete(additionalProperties, "locked")
		delete(additionalProperties, "no_ssh_keys")
		delete(additionalProperties, "operating_system")
		delete(additionalProperties, "plan")
		delete(additionalProperties, "private_ipv4_subnet_size")
		delete(additionalProperties, "project_ssh_keys")
		delete(additionalProperties, "public_ipv4_subnet_size")
		delete(additionalProperties, "tags")
		delete(additionalProperties, "termination_time")
		delete(additionalProperties, "user_ssh_keys")
		delete(additionalProperties, "userdata")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableSpotMarketRequestCreateInputInstanceParameters struct {
	value *SpotMarketRequestCreateInputInstanceParameters
	isSet bool
}

func (v NullableSpotMarketRequestCreateInputInstanceParameters) Get() *SpotMarketRequestCreateInputInstanceParameters {
	return v.value
}

func (v *NullableSpotMarketRequestCreateInputInstanceParameters) Set(val *SpotMarketRequestCreateInputInstanceParameters) {
	v.value = val
	v.isSet = true
}

func (v NullableSpotMarketRequestCreateInputInstanceParameters) IsSet() bool {
	return v.isSet
}

func (v *NullableSpotMarketRequestCreateInputInstanceParameters) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableSpotMarketRequestCreateInputInstanceParameters(val *SpotMarketRequestCreateInputInstanceParameters) *NullableSpotMarketRequestCreateInputInstanceParameters {
	return &NullableSpotMarketRequestCreateInputInstanceParameters{value: val, isSet: true}
}

func (v NullableSpotMarketRequestCreateInputInstanceParameters) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableSpotMarketRequestCreateInputInstanceParameters) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
