// Package repository provides adapter implementations for domain repository interfaces.
// This package bridges the domain layer with external services like Google Calendar API.
package repository

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	gcal "google.golang.org/api/calendar/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"

	"github.com/stainedhead/go-goog-cli/internal/domain/calendar"
)

// ErrInvalidCalendarRequest is returned for invalid calendar request errors.
var ErrInvalidCalendarRequest = errors.New("invalid calendar request")

// GCalService wraps the Google Calendar service and provides access to
// separate repository implementations for events, calendars, ACL, and free/busy.
type GCalService struct {
	service *gcal.Service
}

// NewGCalService creates a new GCalService with the given OAuth2 token source.
// The token source is used to authenticate requests to the Google Calendar API.
func NewGCalService(ctx context.Context, tokenSource oauth2.TokenSource) (*GCalService, error) {
	httpClient := oauth2.NewClient(ctx, tokenSource)
	service, err := gcal.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("failed to create calendar service: %w", err)
	}
	return &GCalService{
		service: service,
	}, nil
}

// NewGCalServiceWithService creates a new GCalService with an existing service.
// This is primarily useful for testing with a mock service.
func NewGCalServiceWithService(service *gcal.Service) *GCalService {
	return &GCalService{
		service: service,
	}
}

// Events returns an EventRepository for event operations.
func (s *GCalService) Events() *GCalEventRepository {
	return &GCalEventRepository{service: s.service}
}

// Calendars returns a CalendarRepository for calendar operations.
func (s *GCalService) Calendars() *GCalCalendarRepository {
	return &GCalCalendarRepository{service: s.service}
}

// ACL returns an ACLRepository for ACL operations.
func (s *GCalService) ACL() *GCalACLRepository {
	return &GCalACLRepository{service: s.service}
}

// FreeBusy returns a FreeBusyRepository for free/busy operations.
func (s *GCalService) FreeBusy() *GCalFreeBusyRepository {
	return &GCalFreeBusyRepository{service: s.service}
}

// -----------------------------------------------------------------------------
// GCalEventRepository Implementation
// -----------------------------------------------------------------------------

// GCalEventRepository implements calendar.EventRepository using the Google Calendar API.
type GCalEventRepository struct {
	service *gcal.Service
}

// Ensure GCalEventRepository implements calendar.EventRepository.
var _ calendar.EventRepository = (*GCalEventRepository)(nil)

// List returns events from a calendar within the specified time range.
func (r *GCalEventRepository) List(ctx context.Context, calendarID string, timeMin, timeMax time.Time) ([]*calendar.Event, error) {
	if !timeMin.Before(timeMax) {
		return nil, calendar.ErrInvalidTimeRange
	}

	call := r.service.Events.List(calendarID).
		Context(ctx).
		TimeMin(timeMin.Format(time.RFC3339)).
		TimeMax(timeMax.Format(time.RFC3339)).
		SingleEvents(true).
		OrderBy("startTime")

	var events []*calendar.Event
	err := call.Pages(ctx, func(page *gcal.Events) error {
		for _, item := range page.Items {
			event := gcalEventToDomain(item)
			if event != nil {
				event.CalendarID = calendarID
				events = append(events, event)
			}
		}
		return nil
	})
	if err != nil {
		return nil, mapAPIError(err, "event")
	}

	return events, nil
}

// Get retrieves a single event by ID.
func (r *GCalEventRepository) Get(ctx context.Context, calendarID, eventID string) (*calendar.Event, error) {
	gcalEvent, err := r.service.Events.Get(calendarID, eventID).Context(ctx).Do()
	if err != nil {
		return nil, mapAPIError(err, "event")
	}

	event := gcalEventToDomain(gcalEvent)
	if event != nil {
		event.CalendarID = calendarID
	}
	return event, nil
}

// Create creates a new event in the specified calendar.
func (r *GCalEventRepository) Create(ctx context.Context, calendarID string, event *calendar.Event) (*calendar.Event, error) {
	gcalEvent := domainEventToGcal(event)
	if gcalEvent == nil {
		return nil, ErrInvalidCalendarRequest
	}

	created, err := r.service.Events.Insert(calendarID, gcalEvent).Context(ctx).Do()
	if err != nil {
		return nil, mapAPIError(err, "event")
	}

	result := gcalEventToDomain(created)
	if result != nil {
		result.CalendarID = calendarID
	}
	return result, nil
}

// Update updates an existing event.
func (r *GCalEventRepository) Update(ctx context.Context, calendarID string, event *calendar.Event) (*calendar.Event, error) {
	gcalEvent := domainEventToGcal(event)
	if gcalEvent == nil || event.ID == "" {
		return nil, ErrInvalidCalendarRequest
	}

	updated, err := r.service.Events.Update(calendarID, event.ID, gcalEvent).Context(ctx).Do()
	if err != nil {
		return nil, mapAPIError(err, "event")
	}

	result := gcalEventToDomain(updated)
	if result != nil {
		result.CalendarID = calendarID
	}
	return result, nil
}

// Delete removes an event from a calendar.
func (r *GCalEventRepository) Delete(ctx context.Context, calendarID, eventID string) error {
	err := r.service.Events.Delete(calendarID, eventID).Context(ctx).Do()
	if err != nil {
		return mapAPIError(err, "event")
	}
	return nil
}

// Move moves an event to a different calendar.
func (r *GCalEventRepository) Move(ctx context.Context, sourceCalendarID, eventID, destinationCalendarID string) (*calendar.Event, error) {
	moved, err := r.service.Events.Move(sourceCalendarID, eventID, destinationCalendarID).Context(ctx).Do()
	if err != nil {
		return nil, mapAPIError(err, "event")
	}

	result := gcalEventToDomain(moved)
	if result != nil {
		result.CalendarID = destinationCalendarID
	}
	return result, nil
}

// QuickAdd creates an event based on a simple text string.
func (r *GCalEventRepository) QuickAdd(ctx context.Context, calendarID, text string) (*calendar.Event, error) {
	created, err := r.service.Events.QuickAdd(calendarID, text).Context(ctx).Do()
	if err != nil {
		return nil, mapAPIError(err, "event")
	}

	result := gcalEventToDomain(created)
	if result != nil {
		result.CalendarID = calendarID
	}
	return result, nil
}

// Instances returns instances of a recurring event.
func (r *GCalEventRepository) Instances(ctx context.Context, calendarID, eventID string, timeMin, timeMax time.Time) ([]*calendar.Event, error) {
	call := r.service.Events.Instances(calendarID, eventID).Context(ctx)

	if !timeMin.IsZero() {
		call = call.TimeMin(timeMin.Format(time.RFC3339))
	}
	if !timeMax.IsZero() {
		call = call.TimeMax(timeMax.Format(time.RFC3339))
	}

	var events []*calendar.Event
	err := call.Pages(ctx, func(page *gcal.Events) error {
		for _, item := range page.Items {
			event := gcalEventToDomain(item)
			if event != nil {
				event.CalendarID = calendarID
				events = append(events, event)
			}
		}
		return nil
	})
	if err != nil {
		return nil, mapAPIError(err, "event")
	}

	return events, nil
}

// RSVP updates the current user's response to an event.
func (r *GCalEventRepository) RSVP(ctx context.Context, calendarID, eventID, response string) error {
	if !calendar.IsValidResponseStatus(response) {
		return ErrInvalidCalendarRequest
	}

	// Get the current event
	gcalEvent, err := r.service.Events.Get(calendarID, eventID).Context(ctx).Do()
	if err != nil {
		return mapAPIError(err, "event")
	}

	// Find self in attendees and update response
	for _, attendee := range gcalEvent.Attendees {
		if attendee.Self {
			attendee.ResponseStatus = response
			break
		}
	}

	// Update the event
	_, err = r.service.Events.Update(calendarID, eventID, gcalEvent).Context(ctx).Do()
	if err != nil {
		return mapAPIError(err, "event")
	}

	return nil
}

// -----------------------------------------------------------------------------
// GCalCalendarRepository Implementation
// -----------------------------------------------------------------------------

// GCalCalendarRepository implements calendar.CalendarRepository using the Google Calendar API.
type GCalCalendarRepository struct {
	service *gcal.Service
}

// Ensure GCalCalendarRepository implements calendar.CalendarRepository.
var _ calendar.CalendarRepository = (*GCalCalendarRepository)(nil)

// List returns all calendars accessible to the user.
func (r *GCalCalendarRepository) List(ctx context.Context) ([]*calendar.Calendar, error) {
	var calendars []*calendar.Calendar

	call := r.service.CalendarList.List().Context(ctx)
	err := call.Pages(ctx, func(page *gcal.CalendarList) error {
		for _, item := range page.Items {
			cal := gcalCalendarToDomain(item)
			if cal != nil {
				calendars = append(calendars, cal)
			}
		}
		return nil
	})
	if err != nil {
		return nil, mapAPIError(err, "calendar")
	}

	return calendars, nil
}

// Get retrieves a single calendar by ID.
func (r *GCalCalendarRepository) Get(ctx context.Context, calendarID string) (*calendar.Calendar, error) {
	gcalCal, err := r.service.CalendarList.Get(calendarID).Context(ctx).Do()
	if err != nil {
		return nil, mapAPIError(err, "calendar")
	}
	return gcalCalendarToDomain(gcalCal), nil
}

// Create creates a new calendar.
func (r *GCalCalendarRepository) Create(ctx context.Context, cal *calendar.Calendar) (*calendar.Calendar, error) {
	gcalCal := domainCalendarToGcal(cal)
	if gcalCal == nil {
		return nil, ErrInvalidCalendarRequest
	}

	created, err := r.service.Calendars.Insert(gcalCal).Context(ctx).Do()
	if err != nil {
		return nil, mapAPIError(err, "calendar")
	}

	// Return calendar info from CalendarList for full metadata
	return r.Get(ctx, created.Id)
}

// Update updates an existing calendar.
func (r *GCalCalendarRepository) Update(ctx context.Context, cal *calendar.Calendar) (*calendar.Calendar, error) {
	gcalCal := domainCalendarToGcal(cal)
	if gcalCal == nil || cal.ID == "" {
		return nil, ErrInvalidCalendarRequest
	}

	_, err := r.service.Calendars.Update(cal.ID, gcalCal).Context(ctx).Do()
	if err != nil {
		return nil, mapAPIError(err, "calendar")
	}

	// Return calendar info from CalendarList for full metadata
	return r.Get(ctx, cal.ID)
}

// Delete removes a calendar.
func (r *GCalCalendarRepository) Delete(ctx context.Context, calendarID string) error {
	err := r.service.Calendars.Delete(calendarID).Context(ctx).Do()
	if err != nil {
		return mapAPIError(err, "calendar")
	}
	return nil
}

// Clear clears all events from a calendar.
func (r *GCalCalendarRepository) Clear(ctx context.Context, calendarID string) error {
	err := r.service.Calendars.Clear(calendarID).Context(ctx).Do()
	if err != nil {
		return mapAPIError(err, "calendar")
	}
	return nil
}

// -----------------------------------------------------------------------------
// GCalACLRepository Implementation
// -----------------------------------------------------------------------------

// GCalACLRepository implements calendar.ACLRepository using the Google Calendar API.
type GCalACLRepository struct {
	service *gcal.Service
}

// Ensure GCalACLRepository implements calendar.ACLRepository.
var _ calendar.ACLRepository = (*GCalACLRepository)(nil)

// List returns all ACL rules for a calendar.
func (r *GCalACLRepository) List(ctx context.Context, calendarID string) ([]*calendar.ACLRule, error) {
	var rules []*calendar.ACLRule

	call := r.service.Acl.List(calendarID).Context(ctx)
	err := call.Pages(ctx, func(page *gcal.Acl) error {
		for _, item := range page.Items {
			rule := gcalACLToDomain(item)
			if rule != nil {
				rules = append(rules, rule)
			}
		}
		return nil
	})
	if err != nil {
		return nil, mapAPIError(err, "acl")
	}

	return rules, nil
}

// Get retrieves a single ACL rule by ID.
func (r *GCalACLRepository) Get(ctx context.Context, calendarID, ruleID string) (*calendar.ACLRule, error) {
	gcalRule, err := r.service.Acl.Get(calendarID, ruleID).Context(ctx).Do()
	if err != nil {
		return nil, mapAPIError(err, "acl")
	}
	return gcalACLToDomain(gcalRule), nil
}

// Insert creates a new ACL rule for a calendar.
func (r *GCalACLRepository) Insert(ctx context.Context, calendarID string, rule *calendar.ACLRule) (*calendar.ACLRule, error) {
	gcalRule := domainACLToGcal(rule)
	if gcalRule == nil {
		return nil, ErrInvalidCalendarRequest
	}

	created, err := r.service.Acl.Insert(calendarID, gcalRule).Context(ctx).Do()
	if err != nil {
		return nil, mapAPIError(err, "acl")
	}

	return gcalACLToDomain(created), nil
}

// Update updates an existing ACL rule.
func (r *GCalACLRepository) Update(ctx context.Context, calendarID string, rule *calendar.ACLRule) (*calendar.ACLRule, error) {
	gcalRule := domainACLToGcal(rule)
	if gcalRule == nil || rule.ID == "" {
		return nil, ErrInvalidCalendarRequest
	}

	updated, err := r.service.Acl.Update(calendarID, rule.ID, gcalRule).Context(ctx).Do()
	if err != nil {
		return nil, mapAPIError(err, "acl")
	}

	return gcalACLToDomain(updated), nil
}

// Delete removes an ACL rule from a calendar.
func (r *GCalACLRepository) Delete(ctx context.Context, calendarID, ruleID string) error {
	err := r.service.Acl.Delete(calendarID, ruleID).Context(ctx).Do()
	if err != nil {
		return mapAPIError(err, "acl")
	}
	return nil
}

// -----------------------------------------------------------------------------
// GCalFreeBusyRepository Implementation
// -----------------------------------------------------------------------------

// GCalFreeBusyRepository implements calendar.FreeBusyRepository using the Google Calendar API.
type GCalFreeBusyRepository struct {
	service *gcal.Service
}

// Ensure GCalFreeBusyRepository implements calendar.FreeBusyRepository.
var _ calendar.FreeBusyRepository = (*GCalFreeBusyRepository)(nil)

// Query returns free/busy information for the specified calendars and time range.
func (r *GCalFreeBusyRepository) Query(ctx context.Context, request *calendar.FreeBusyRequest) (*calendar.FreeBusyResponse, error) {
	if request == nil {
		return nil, ErrInvalidCalendarRequest
	}

	if !request.TimeMin.Before(request.TimeMax) {
		return nil, calendar.ErrInvalidTimeRange
	}

	// Build calendar items
	items := make([]*gcal.FreeBusyRequestItem, len(request.CalendarIDs))
	for i, id := range request.CalendarIDs {
		items[i] = &gcal.FreeBusyRequestItem{Id: id}
	}

	gcalRequest := &gcal.FreeBusyRequest{
		TimeMin: request.TimeMin.Format(time.RFC3339),
		TimeMax: request.TimeMax.Format(time.RFC3339),
		Items:   items,
	}

	gcalResponse, err := r.service.Freebusy.Query(gcalRequest).Context(ctx).Do()
	if err != nil {
		return nil, mapAPIError(err, "freebusy")
	}

	// Convert response
	response := &calendar.FreeBusyResponse{
		Calendars: make(map[string][]*calendar.TimePeriod),
	}

	for calID, calInfo := range gcalResponse.Calendars {
		var periods []*calendar.TimePeriod
		for _, busy := range calInfo.Busy {
			start, _ := time.Parse(time.RFC3339, busy.Start)
			end, _ := time.Parse(time.RFC3339, busy.End)
			period, err := calendar.NewTimePeriod(start, end)
			if err == nil {
				periods = append(periods, period)
			}
		}
		response.Calendars[calID] = periods
	}

	return response, nil
}

// -----------------------------------------------------------------------------
// Conversion Functions
// -----------------------------------------------------------------------------

// parseEventDateTime converts a Google Calendar EventDateTime to a time.Time.
// It returns the time, a boolean indicating if this is an all-day event, and any error.
func parseEventDateTime(dt *gcal.EventDateTime) (time.Time, bool, error) {
	if dt == nil {
		return time.Time{}, false, nil
	}

	// Check for all-day event (date only)
	if dt.Date != "" {
		t, err := time.Parse("2006-01-02", dt.Date)
		if err != nil {
			return time.Time{}, true, fmt.Errorf("invalid date format: %w", err)
		}
		return t, true, nil
	}

	// Parse datetime
	if dt.DateTime != "" {
		t, err := time.Parse(time.RFC3339, dt.DateTime)
		if err != nil {
			return time.Time{}, false, fmt.Errorf("invalid datetime format: %w", err)
		}
		return t, false, nil
	}

	return time.Time{}, false, nil
}

// parseRecurrence returns a copy of the recurrence rules.
func parseRecurrence(rrules []string) []string {
	if rrules == nil {
		return nil
	}
	result := make([]string, len(rrules))
	copy(result, rrules)
	return result
}

// gcalEventToDomain converts a Google Calendar Event to a domain Event.
func gcalEventToDomain(event *gcal.Event) *calendar.Event {
	if event == nil {
		return nil
	}

	start, allDay, _ := parseEventDateTime(event.Start)
	end, _, _ := parseEventDateTime(event.End)
	created, _ := time.Parse(time.RFC3339, event.Created)
	updated, _ := time.Parse(time.RFC3339, event.Updated)

	domainEvent := &calendar.Event{
		ID:          event.Id,
		Title:       event.Summary,
		Description: event.Description,
		Location:    event.Location,
		Start:       start,
		End:         end,
		AllDay:      allDay,
		Recurrence:  parseRecurrence(event.Recurrence),
		Status:      event.Status,
		Visibility:  event.Visibility,
		ColorID:     event.ColorId,
		Created:     created,
		Updated:     updated,
		HTMLLink:    event.HtmlLink,
	}

	// Convert attendees
	if len(event.Attendees) > 0 {
		domainEvent.Attendees = make([]*calendar.Attendee, len(event.Attendees))
		for i, a := range event.Attendees {
			domainEvent.Attendees[i] = &calendar.Attendee{
				Email:          a.Email,
				DisplayName:    a.DisplayName,
				ResponseStatus: a.ResponseStatus,
				Optional:       a.Optional,
				Organizer:      a.Organizer,
				Self:           a.Self,
			}
		}
	}

	// Convert organizer
	if event.Organizer != nil {
		domainEvent.Organizer = &calendar.Attendee{
			Email:       event.Organizer.Email,
			DisplayName: event.Organizer.DisplayName,
		}
	}

	// Convert reminders
	if event.Reminders != nil && !event.Reminders.UseDefault && len(event.Reminders.Overrides) > 0 {
		domainEvent.Reminders = make([]*calendar.Reminder, len(event.Reminders.Overrides))
		for i, r := range event.Reminders.Overrides {
			domainEvent.Reminders[i] = &calendar.Reminder{
				Method:  r.Method,
				Minutes: int(r.Minutes),
			}
		}
	}

	// Convert conference data
	if event.ConferenceData != nil && len(event.ConferenceData.EntryPoints) > 0 {
		confType := ""
		if event.ConferenceData.ConferenceSolution != nil && event.ConferenceData.ConferenceSolution.Key != nil {
			confType = event.ConferenceData.ConferenceSolution.Key.Type
		}
		// Find video entry point
		var confURI string
		for _, ep := range event.ConferenceData.EntryPoints {
			if ep.EntryPointType == "video" {
				confURI = ep.Uri
				break
			}
		}
		if confURI != "" {
			domainEvent.ConferenceData = &calendar.ConferenceData{
				Type: confType,
				URI:  confURI,
			}
		}
	}

	return domainEvent
}

// domainEventToGcal converts a domain Event to a Google Calendar Event.
func domainEventToGcal(event *calendar.Event) *gcal.Event {
	if event == nil {
		return nil
	}

	gcalEvent := &gcal.Event{
		Id:          event.ID,
		Summary:     event.Title,
		Description: event.Description,
		Location:    event.Location,
		Status:      event.Status,
		Visibility:  event.Visibility,
		ColorId:     event.ColorID,
		Recurrence:  event.Recurrence,
	}

	// Set start/end times
	if event.AllDay {
		gcalEvent.Start = &gcal.EventDateTime{
			Date: event.Start.Format("2006-01-02"),
		}
		gcalEvent.End = &gcal.EventDateTime{
			Date: event.End.Format("2006-01-02"),
		}
	} else {
		gcalEvent.Start = &gcal.EventDateTime{
			DateTime: event.Start.Format(time.RFC3339),
		}
		gcalEvent.End = &gcal.EventDateTime{
			DateTime: event.End.Format(time.RFC3339),
		}
	}

	// Convert attendees
	if len(event.Attendees) > 0 {
		gcalEvent.Attendees = make([]*gcal.EventAttendee, len(event.Attendees))
		for i, a := range event.Attendees {
			gcalEvent.Attendees[i] = &gcal.EventAttendee{
				Email:          a.Email,
				DisplayName:    a.DisplayName,
				ResponseStatus: a.ResponseStatus,
				Optional:       a.Optional,
				Organizer:      a.Organizer,
			}
		}
	}

	// Convert reminders
	if len(event.Reminders) > 0 {
		gcalEvent.Reminders = &gcal.EventReminders{
			UseDefault: false,
			Overrides:  make([]*gcal.EventReminder, len(event.Reminders)),
		}
		for i, r := range event.Reminders {
			gcalEvent.Reminders.Overrides[i] = &gcal.EventReminder{
				Method:  r.Method,
				Minutes: int64(r.Minutes),
			}
		}
	}

	return gcalEvent
}

// gcalCalendarToDomain converts a Google Calendar CalendarListEntry to a domain Calendar.
func gcalCalendarToDomain(cal *gcal.CalendarListEntry) *calendar.Calendar {
	if cal == nil {
		return nil
	}

	return &calendar.Calendar{
		ID:          cal.Id,
		Title:       cal.Summary,
		Description: cal.Description,
		TimeZone:    cal.TimeZone,
		ColorID:     cal.ColorId,
		Primary:     cal.Primary,
		Selected:    cal.Selected,
		AccessRole:  cal.AccessRole,
	}
}

// domainCalendarToGcal converts a domain Calendar to a Google Calendar Calendar.
func domainCalendarToGcal(cal *calendar.Calendar) *gcal.Calendar {
	if cal == nil {
		return nil
	}

	return &gcal.Calendar{
		Id:          cal.ID,
		Summary:     cal.Title,
		Description: cal.Description,
		TimeZone:    cal.TimeZone,
	}
}

// gcalACLToDomain converts a Google Calendar AclRule to a domain ACLRule.
func gcalACLToDomain(rule *gcal.AclRule) *calendar.ACLRule {
	if rule == nil {
		return nil
	}

	domainRule := &calendar.ACLRule{
		ID:   rule.Id,
		Role: rule.Role,
	}

	if rule.Scope != nil {
		domainRule.Scope = &calendar.ACLScope{
			Type:  rule.Scope.Type,
			Value: rule.Scope.Value,
		}
	}

	return domainRule
}

// domainACLToGcal converts a domain ACLRule to a Google Calendar AclRule.
func domainACLToGcal(rule *calendar.ACLRule) *gcal.AclRule {
	if rule == nil {
		return nil
	}

	gcalRule := &gcal.AclRule{
		Id:   rule.ID,
		Role: rule.Role,
	}

	if rule.Scope != nil {
		gcalRule.Scope = &gcal.AclRuleScope{
			Type:  rule.Scope.Type,
			Value: rule.Scope.Value,
		}
	}

	return gcalRule
}

// mapAPIError maps Google Calendar API errors to domain errors.
func mapAPIError(err error, resource string) error {
	if err == nil {
		return nil
	}

	var apiErr *googleapi.Error
	if errors.As(err, &apiErr) {
		switch apiErr.Code {
		case http.StatusNotFound:
			switch resource {
			case "event":
				return calendar.ErrEventNotFound
			case "calendar":
				return calendar.ErrCalendarNotFound
			case "acl":
				return calendar.ErrACLNotFound
			default:
				return fmt.Errorf("%s not found", resource)
			}
		case http.StatusBadRequest:
			// Check for specific validation errors
			if containsTimeRangeError(apiErr.Message) {
				return calendar.ErrInvalidTimeRange
			}
			return fmt.Errorf("%w: %s", ErrBadRequest, apiErr.Message)
		case http.StatusTooManyRequests:
			return fmt.Errorf("%w: %s", ErrRateLimited, apiErr.Message)
		case http.StatusInternalServerError,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout:
			return fmt.Errorf("%w: %s", ErrTemporary, apiErr.Message)
		}
	}

	return err
}

// containsTimeRangeError checks if the error message indicates a time range validation error.
func containsTimeRangeError(msg string) bool {
	// Common time range error messages from the Calendar API
	return msg == "Invalid time range" ||
		msg == "The requested time range is invalid" ||
		msg == "Start time must be before end time"
}
