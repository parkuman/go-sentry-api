package sentry

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/parkuman/go-sentry-api/datatype"
)

// Tag is used for a event
type Tag struct {
	Value *string `json:"value,omitempty"`
	Key   *string `json:"key,omitempty"`
}

// User is the user that was affected
type User struct {
	Username *string  `json:"username,omitempty"`
	Email    *string  `json:"email,omitempty"`
	ID       *string  `json:"id,omitempty"`
	Name     *string  `json:"name,omitempty"`
	Role     *string  `json:"role,omitempty"`
	RoleName *string  `json:"roleName,omitempty"`
	Projects []string `json:"projects,omitempty"`
}

// Entry is the entry for the message/stacktrace/etc...
type Entry struct {
	Type string          `json:"type,omitempty"`
	Data json.RawMessage `json:"data,omitempty"`
}

// GetInterface will convert the entry into a go interface
func (e *Entry) GetInterface() (string, interface{}, error) {
	var destination interface{}

	switch e.Type {
	case "message":
		destination = new(datatype.Message)
	case "stacktrace":
		destination = new(datatype.Stacktrace)
	case "exception":
		destination = new(datatype.Exception)
	case "request":
		destination = new(datatype.Request)
	case "template":
		destination = new(datatype.Template)
	case "user":
		destination = new(datatype.User)
	case "query":
		destination = new(datatype.Query)
	case "breadcrumbs":
		destination = new(datatype.Breadcrumb)
	}

	err := json.Unmarshal(e.Data, &destination)
	return e.Type, destination, err
}

// Event is the event that was created on the app and sentry reported on
type Event struct {
	EventID         string                  `json:"eventID,omitempty"`
	UserReport      *string                 `json:"userReport,omitempty"`
	NextEventID     *string                 `json:"nextEventID,omitempty"`
	PreviousEventID *string                 `json:"previousEventID,omitempty"`
	Message         *string                 `json:"message,omitempty"`
	ID              *string                 `json:"id,omitempty"`
	Size            *int                    `json:"size,omitempty"`
	Platform        *string                 `json:"platform,omitempty"`
	Type            *string                 `json:"type,omitempty"`
	Metadata        *map[string]string      `json:"metadata,omitempty"`
	Tags            *[]Tag                  `json:"tags,omitempty"`
	DateCreated     *time.Time              `json:"dateCreated,omitempty"`
	DateReceived    *time.Time              `json:"dateReceived,omitempty"`
	User            *User                   `json:"user,omitempty"`
	Entries         []Entry                 `json:"entries,omitempty"`
	Packages        *map[string]string      `json:"packages,omitempty"`
	SDK             *map[string]interface{} `json:"sdk,omitempty"`
	Contexts        *map[string]interface{} `json:"contexts,omitempty"`
	Context         *map[string]interface{} `json:"context,omitempty"`
	Release         *Release                `json:"release,omitempty"`
	GroupID         *string                 `json:"groupID,omitempty"`
}

// EventStatsSet represents a set of statistics for an event.
type EventStatsSet struct {
	Data          []EventStatsPoint           `json:"data"`
	Confidence    []EventStatsConfidencePoint `json:"confidence"`
	Order         int                         `json:"order"`
	IsMetricsData bool                        `json:"isMetricsData"`
	Start         int64                       `json:"start"`
	End           int64                       `json:"end"`
	Meta          json.RawMessage             `json:"meta"`
}

// EventStatsPoint represents a single data point in an event's statistics.
type EventStatsPoint struct {
	Timestamp int64
	Values    []EventStatsPointValue
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (p *EventStatsPoint) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if len(raw) != 2 {
		return fmt.Errorf("expected array of length 2, got %d", len(raw))
	}
	if err := json.Unmarshal(raw[0], &p.Timestamp); err != nil {
		return err
	}
	return json.Unmarshal(raw[1], &p.Values)
}

// EventStatsPointValue represents the value of a data point.
type EventStatsPointValue struct {
	Count json.Number `json:"count"`
}

// EventStatsConfidencePoint represents a confidence data point.
type EventStatsConfidencePoint struct {
	Timestamp int64
	Values    []EventStatsConfidencePointValue
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (p *EventStatsConfidencePoint) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if len(raw) != 2 {
		return fmt.Errorf("expected array of length 2, got %d", len(raw))
	}
	if err := json.Unmarshal(raw[0], &p.Timestamp); err != nil {
		return err
	}
	return json.Unmarshal(raw[1], &p.Values)
}

// EventStatsConfidencePointValue represents the value of a confidence point.
type EventStatsConfidencePointValue struct {
	Count *string `json:"count"`
}

// for building the query string of the /organizations/:org/events-stats endpoint
type EventsStatsRequest struct {
	Dataset      string
	End          *time.Time
	Environment  string
	ExcludeOther bool
	Field        []string
	Interval     string
	OrderBy      string
	Partial      bool
	PerPage      int
	Project      []string // Project IDs
	Query        string
	Referrer     string
	Sampling     string
	Sort         string
	Start        *time.Time
	UTC          bool
	YAxis        []string
}

func (r *EventsStatsRequest) ToQueryString() string {
	query := url.Values{}

	query.Add("excludeOther", strconv.Itoa(btoi(r.ExcludeOther)))
	query.Add("partial", strconv.Itoa(btoi(r.Partial)))
	query.Add("per_page", strconv.Itoa(r.PerPage))
	query.Add("utc", strconv.FormatBool(r.UTC))

	if r.Dataset != "" {
		query.Add("dataset", r.Dataset)
	}
	if r.End != nil {
		query.Add("end", r.End.UTC().Format("2006-01-02T15:04:05.999Z"))
	}
	if r.Environment != "" {
		query.Add("environment", r.Environment)
	}
	for _, f := range r.Field {
		query.Add("field", f)
	}
	if r.Interval != "" {
		query.Add("interval", r.Interval)
	}
	if r.OrderBy != "" {
		query.Add("orderby", r.OrderBy)
	}
	for _, p := range r.Project {
		query.Add("project", p)
	}
	if r.Query != "" {
		query.Add("query", r.Query)
	}
	if r.Referrer != "" {
		query.Add("referrer", r.Referrer)
	}
	if r.Sampling != "" {
		query.Add("sampling", r.Sampling)
	}
	if r.Sort != "" {
		query.Add("sort", r.Sort)
	}
	if r.Start != nil {
		query.Add("start", r.Start.UTC().Format("2006-01-02T15:04:05.999Z"))
	}
	for _, y := range r.YAxis {
		query.Add("yAxis", y)
	}
	return query.Encode()
}

// btoi is a helper to convert bool to int (0 or 1)
func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

// for building the query string of the /organizations/:org/events endpoint
type EventsRequest struct {
	Dataset     string
	End         *time.Time
	Environment string
	Field       []string
	PerPage     int
	Project     []string // Project slugs
	Query       string
	Referrer    string
	Sampling    string
	Sort        []string
	Start       *time.Time
	UTC         bool
}

func (r *EventsRequest) ToQueryString() string {
	query := url.Values{}

	query.Add("per_page", strconv.Itoa(r.PerPage))
	query.Add("utc", strconv.FormatBool(r.UTC))

	if r.Dataset != "" {
		query.Add("dataset", r.Dataset)
	}
	if r.End != nil {
		query.Add("end", r.End.UTC().Format("2006-01-02T15:04:05.999Z"))
	}
	if r.Environment != "" {
		query.Add("environment", r.Environment)
	}
	for _, f := range r.Field {
		query.Add("field", f)
	}
	for _, p := range r.Project {
		query.Add("project", p)
	}
	if r.Query != "" {
		query.Add("query", r.Query)
	}
	if r.Referrer != "" {
		query.Add("referrer", r.Referrer)
	}
	if r.Sampling != "" {
		query.Add("sampling", r.Sampling)
	}
	for _, s := range r.Sort {
		query.Add("sort", s)
	}
	if r.Start != nil {
		query.Add("start", r.Start.UTC().Format("2006-01-02T15:04:05.999Z"))
	}
	return query.Encode()
}

// EventData represents a single data object in the events response with dynamic fields.
type EventData map[string]interface{}

// EventsResponse represents the response from the organization events endpoint with dynamic fields.
type EventsResponse struct {
	Data []EventData     `json:"data"`
	Meta json.RawMessage `json:"meta"`
}

// GetProjectEvent will fetch a event on a project
func (c *Client) GetProjectEvent(o Organization, p Project, eventID string) (Event, error) {
	var event Event
	err := c.do("GET", fmt.Sprintf("projects/%s/%s/events/%s", *o.Slug, *p.Slug, eventID), &event, nil)
	return event, err
}

// GetLatestEvent will fetch the latest event for a issue
func (c *Client) GetLatestEvent(i Issue) (Event, error) {
	var event Event
	err := c.do("GET", fmt.Sprintf("issues/%s/events/latest", *i.ID), &event, nil)
	return event, err
}

// GetOldestEvent will fetch the latest event for a issue
func (c *Client) GetOldestEvent(i Issue) (Event, error) {
	var event Event
	err := c.do("GET", fmt.Sprintf("issues/%s/events/oldest", *i.ID), &event, nil)
	return event, err
}

// GetEvents will fetch all events for a given org and project
func (c *Client) GetEvents(o Organization, params *EventsRequest) (EventsResponse, error) {
	events := EventsResponse{}
	_, err := c.doWithPaginationQuery("GET", fmt.Sprintf("organizations/%s/events", *o.Slug), &events, nil, params)
	return events, err
}

// GetEventsStats will fetch stats for events for a given org and project
func (c *Client) GetEventsStats(o Organization, params *EventsStatsRequest) (map[string]EventStatsSet, error) {
	stats := make(map[string]EventStatsSet)
	_, err := c.doWithPaginationQuery("GET", fmt.Sprintf("organizations/%s/events-stats", *o.Slug), &stats, nil, params)
	return stats, err
}
