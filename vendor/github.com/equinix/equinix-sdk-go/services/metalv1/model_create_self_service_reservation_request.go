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

// checks if the CreateSelfServiceReservationRequest type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CreateSelfServiceReservationRequest{}

// CreateSelfServiceReservationRequest struct for CreateSelfServiceReservationRequest
type CreateSelfServiceReservationRequest struct {
	Item                 []SelfServiceReservationItemRequest        `json:"item,omitempty"`
	Notes                *string                                    `json:"notes,omitempty"`
	Period               *CreateSelfServiceReservationRequestPeriod `json:"period,omitempty"`
	StartDate            *time.Time                                 `json:"start_date,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _CreateSelfServiceReservationRequest CreateSelfServiceReservationRequest

// NewCreateSelfServiceReservationRequest instantiates a new CreateSelfServiceReservationRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCreateSelfServiceReservationRequest() *CreateSelfServiceReservationRequest {
	this := CreateSelfServiceReservationRequest{}
	return &this
}

// NewCreateSelfServiceReservationRequestWithDefaults instantiates a new CreateSelfServiceReservationRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCreateSelfServiceReservationRequestWithDefaults() *CreateSelfServiceReservationRequest {
	this := CreateSelfServiceReservationRequest{}
	return &this
}

// GetItem returns the Item field value if set, zero value otherwise.
func (o *CreateSelfServiceReservationRequest) GetItem() []SelfServiceReservationItemRequest {
	if o == nil || IsNil(o.Item) {
		var ret []SelfServiceReservationItemRequest
		return ret
	}
	return o.Item
}

// GetItemOk returns a tuple with the Item field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateSelfServiceReservationRequest) GetItemOk() ([]SelfServiceReservationItemRequest, bool) {
	if o == nil || IsNil(o.Item) {
		return nil, false
	}
	return o.Item, true
}

// HasItem returns a boolean if a field has been set.
func (o *CreateSelfServiceReservationRequest) HasItem() bool {
	if o != nil && !IsNil(o.Item) {
		return true
	}

	return false
}

// SetItem gets a reference to the given []SelfServiceReservationItemRequest and assigns it to the Item field.
func (o *CreateSelfServiceReservationRequest) SetItem(v []SelfServiceReservationItemRequest) {
	o.Item = v
}

// GetNotes returns the Notes field value if set, zero value otherwise.
func (o *CreateSelfServiceReservationRequest) GetNotes() string {
	if o == nil || IsNil(o.Notes) {
		var ret string
		return ret
	}
	return *o.Notes
}

// GetNotesOk returns a tuple with the Notes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateSelfServiceReservationRequest) GetNotesOk() (*string, bool) {
	if o == nil || IsNil(o.Notes) {
		return nil, false
	}
	return o.Notes, true
}

// HasNotes returns a boolean if a field has been set.
func (o *CreateSelfServiceReservationRequest) HasNotes() bool {
	if o != nil && !IsNil(o.Notes) {
		return true
	}

	return false
}

// SetNotes gets a reference to the given string and assigns it to the Notes field.
func (o *CreateSelfServiceReservationRequest) SetNotes(v string) {
	o.Notes = &v
}

// GetPeriod returns the Period field value if set, zero value otherwise.
func (o *CreateSelfServiceReservationRequest) GetPeriod() CreateSelfServiceReservationRequestPeriod {
	if o == nil || IsNil(o.Period) {
		var ret CreateSelfServiceReservationRequestPeriod
		return ret
	}
	return *o.Period
}

// GetPeriodOk returns a tuple with the Period field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateSelfServiceReservationRequest) GetPeriodOk() (*CreateSelfServiceReservationRequestPeriod, bool) {
	if o == nil || IsNil(o.Period) {
		return nil, false
	}
	return o.Period, true
}

// HasPeriod returns a boolean if a field has been set.
func (o *CreateSelfServiceReservationRequest) HasPeriod() bool {
	if o != nil && !IsNil(o.Period) {
		return true
	}

	return false
}

// SetPeriod gets a reference to the given CreateSelfServiceReservationRequestPeriod and assigns it to the Period field.
func (o *CreateSelfServiceReservationRequest) SetPeriod(v CreateSelfServiceReservationRequestPeriod) {
	o.Period = &v
}

// GetStartDate returns the StartDate field value if set, zero value otherwise.
func (o *CreateSelfServiceReservationRequest) GetStartDate() time.Time {
	if o == nil || IsNil(o.StartDate) {
		var ret time.Time
		return ret
	}
	return *o.StartDate
}

// GetStartDateOk returns a tuple with the StartDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateSelfServiceReservationRequest) GetStartDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.StartDate) {
		return nil, false
	}
	return o.StartDate, true
}

// HasStartDate returns a boolean if a field has been set.
func (o *CreateSelfServiceReservationRequest) HasStartDate() bool {
	if o != nil && !IsNil(o.StartDate) {
		return true
	}

	return false
}

// SetStartDate gets a reference to the given time.Time and assigns it to the StartDate field.
func (o *CreateSelfServiceReservationRequest) SetStartDate(v time.Time) {
	o.StartDate = &v
}

func (o CreateSelfServiceReservationRequest) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CreateSelfServiceReservationRequest) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Item) {
		toSerialize["item"] = o.Item
	}
	if !IsNil(o.Notes) {
		toSerialize["notes"] = o.Notes
	}
	if !IsNil(o.Period) {
		toSerialize["period"] = o.Period
	}
	if !IsNil(o.StartDate) {
		toSerialize["start_date"] = o.StartDate
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *CreateSelfServiceReservationRequest) UnmarshalJSON(bytes []byte) (err error) {
	varCreateSelfServiceReservationRequest := _CreateSelfServiceReservationRequest{}

	err = json.Unmarshal(bytes, &varCreateSelfServiceReservationRequest)

	if err != nil {
		return err
	}

	*o = CreateSelfServiceReservationRequest(varCreateSelfServiceReservationRequest)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(bytes, &additionalProperties); err == nil {
		delete(additionalProperties, "item")
		delete(additionalProperties, "notes")
		delete(additionalProperties, "period")
		delete(additionalProperties, "start_date")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableCreateSelfServiceReservationRequest struct {
	value *CreateSelfServiceReservationRequest
	isSet bool
}

func (v NullableCreateSelfServiceReservationRequest) Get() *CreateSelfServiceReservationRequest {
	return v.value
}

func (v *NullableCreateSelfServiceReservationRequest) Set(val *CreateSelfServiceReservationRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableCreateSelfServiceReservationRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableCreateSelfServiceReservationRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCreateSelfServiceReservationRequest(val *CreateSelfServiceReservationRequest) *NullableCreateSelfServiceReservationRequest {
	return &NullableCreateSelfServiceReservationRequest{value: val, isSet: true}
}

func (v NullableCreateSelfServiceReservationRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCreateSelfServiceReservationRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
