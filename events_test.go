package sentry

import (
	"testing"
	"time"

	"github.com/parkuman/go-sentry-api/datatype"
)

func TestEventsResource(t *testing.T) {

	client := newTestClient(t)

	org, err := client.GetOrganization(getDefaultOrg())
	if err != nil {
		t.Fatal(err)
	}

	team, cleanup := createTeamHelper(t)
	defer cleanup()

	project, cleanupproj := createProjectHelper(t, team)
	defer cleanupproj()

	createMessagesHelper(t, client, org, project, 10)

	issues, _, err := client.GetIssues(org, project, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(issues) == 0 {
		t.Fatalf("no issues found for project %s", project.Name)
	}

	t.Run("Read out all events for a issue", func(t *testing.T) {

		for _, issue := range issues {

			t.Logf("Issue: %s", *issue.ID)

			_, err = client.GetLatestEvent(issue)
			if err != nil {
				t.Error(err)
			}

			_, err = client.GetOldestEvent(issue)
			if err != nil {
				t.Error(err)
			}

		}

		t.Run("Convert all entries to proper types if we can", func(t *testing.T) {

			oldest, err := client.GetOldestEvent(issues[0])
			if err != nil {
				t.Error(err)
			}

			if len(oldest.Entries) == 0 {
				t.Fatal("no entries found and is nil")
			}

			for _, entry := range oldest.Entries {
				name, inter, err := entry.GetInterface()
				if err != nil {
					t.Error(err)
				}
				t.Logf("Name: %s, Interface: %v, Error: %v", name, inter, err)

				if name == "message" {
					if inter == nil {
						t.Error("iter message is nil")
					}

					if m := inter.(*datatype.Message); m != nil {
						t.Logf("Messages: %v", m.Message)
					} else {
						t.Errorf("Message is nil %v", m)
					}
				}
			}
		})

	})

}

func TestEventsEndpoints(t *testing.T) {
	client := newTestClient(t)
	org, err := client.GetOrganization(getDefaultOrg())
	if err != nil {
		t.Fatal(err)
	}

	team, cleanupTeam := createTeamHelper(t)
	defer cleanupTeam()

	project, cleanupProject := createProjectHelper(t, team)
	defer cleanupProject()

	createMessagesHelper(t, client, org, project, 5)

	// Give Sentry a moment to process the events
	time.Sleep(2 * time.Second)

	t.Run("GetEvents", func(t *testing.T) {
		params := &EventsRequest{
			Project: []string{*project.Slug},
		}
		events, err := client.GetEvents(org, params)
		if err != nil {
			t.Fatalf("GetEvents failed: %s", err)
		}
		if len(events.Data) == 0 {
			t.Fatal("GetEvents returned no events")
		}
	})

	t.Run("GetEventsStats", func(t *testing.T) {
		params := &EventsStatsRequest{
			Project: []string{project.ID},
			YAxis:   []string{"count()"},
		}
		stats, err := client.GetEventsStats(org, params)
		if err != nil {
			t.Fatalf("GetEventsStats failed: %s", err)
		}
		if len(stats) == 0 {
			t.Fatal("GetEventsStats returned no stats")
		}
	})
}

