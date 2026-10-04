package realtime_test

import (
	"testing"
	"time"

	"backend-alusi-go/pkg/realtime"
)

func TestSSEHub_LifecycleAndBroadcast(t *testing.T) {
	hub := realtime.NewSSEHub()
	hub.Start()
	defer hub.Stop()

	client := &realtime.Client{
		ID:       "client-1",
		SendChan: make(chan []byte, 10),
	}

	hub.RegisterClient(client)
	time.Sleep(50 * time.Millisecond)

	if hub.ActiveClientsCount() != 1 {
		t.Fatalf("Expected 1 active client, got %d", hub.ActiveClientsCount())
	}

	// Broadcast an event
	event := realtime.StatusEvent{
		Type:           "status_change",
		AppID:          1,
		AppSlug:        "simbatik",
		Nama:           "SIMBATIK",
		StatusLayanan:  "kendala",
		PreviousStatus: "operasional",
		ResponseTimeMs: 4500,
		Timestamp:      time.Now(),
	}

	hub.BroadcastStatusEvent(event)

	select {
	case msg := <-client.SendChan:
		msgStr := string(msg)
		if len(msgStr) == 0 {
			t.Fatal("Expected non-empty SSE message")
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Timed out waiting for SSE broadcast")
	}

	// Unregister
	hub.UnregisterClient(client)
	time.Sleep(50 * time.Millisecond)

	if hub.ActiveClientsCount() != 0 {
		t.Fatalf("Expected 0 active clients after unregister, got %d", hub.ActiveClientsCount())
	}
}
