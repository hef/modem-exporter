package client

import (
	"github.com/antchfx/htmlquery"
	"go.uber.org/zap"
	"os"
	"testing"
)

func TestParseStatus(t *testing.T) {
	f, err := os.Open("testdata/cmconnectionstatus.html")
	if err != nil {
		t.Errorf("failed to open file: %s", err)
	}

	doc, err := htmlquery.Parse(f)
	if err != nil {
		t.Errorf("failed to parse file: %s", err)
	}

	actualDownstreamChannels, actualUpstreamChannels, err := parseStatusPage(zap.NewNop(), doc)
	if err != nil {
		t.Errorf("unexpected error parsing status page: %s", err)
	}

	expectedDownstreamChannels := []DownstreamBoundedChannel{
		{5, "Locked", "QAM256", 507000000, 3.8, 33.8, 2968242, 205686},
		{1, "Locked", "QAM256", 477000000, 4.6, 42.2, 68, 122},
		{2, "Locked", "QAM256", 483000000, 4.9, 42.5, 78, 124},
		{3, "Locked", "QAM256", 489000000, 4.3, 42.2, 876, 124},
		{4, "Locked", "QAM256", 495000000, 4.5, 34.0, 84357, 20188},
		{6, "Locked", "QAM256", 513000000, 4.3, 34.2, 10026325, 39921},
		{7, "Locked", "QAM256", 519000000, 4.4, 29.9, 154966327, 33439427},
		{8, "Locked", "QAM256", 525000000, 3.8, 35.0, 73961, 3939},
		{9, "Locked", "QAM256", 531000000, 4.1, 33.4, 523975, 49935},
		{10, "Locked", "QAM256", 537000000, 3.2, 37.8, 2912, 77},
		{11, "Locked", "QAM256", 543000000, 3.3, 40.0, 86, 82},
		{12, "Locked", "QAM256", 549000000, 3.2, 28.7, 2082614, 84381},
		{13, "Locked", "QAM256", 555000000, 3.3, 33.4, 30, 75},
		{14, "Locked", "QAM256", 561000000, 3.4, 40.5, 23, 80},
		{15, "Locked", "QAM256", 567000000, 3.3, 41.4, 57, 47},
		{16, "Locked", "QAM256", 573000000, 2.9, 35.1, 35, 63},
		{17, "Locked", "QAM256", 579000000, 2.8, 39.5, 51963, 63},
		{18, "Locked", "QAM256", 585000000, 2.8, 36.8, 5615, 472},
		{19, "Locked", "QAM256", 591000000, 3.0, 34.4, 428632, 27451},
		{31, "Locked", "Other", 722000000, 4.2, 39.8, 1651639586, 18},
		{33, "Not Locked", "QAM256", 0, 0.0, 40.8, 53, 63},
		{34, "Locked", "QAM256", 453000000, 4.7, 40.8, 60, 35},
		{35, "Locked", "QAM256", 459000000, 4.4, 41.8, 27, 66},
		{36, "Locked", "QAM256", 465000000, 4.1, 41.9, 59, 38},
		{38, "Locked", "QAM256", 471000000, 4.5, 41.5, 54, 34},
		{39, "Locked", "QAM256", 429000000, 3.8, 41.9, 52, 41},
		{40, "Locked", "QAM256", 435000000, 4.0, 42.0, 58, 27},
		{41, "Locked", "QAM256", 441000000, 4.4, 42.2, 39, 28},
		{42, "Locked", "QAM256", 447000000, 4.6, 41.0, 85, 32},
		{43, "Locked", "QAM256", 405000000, 4.3, 41.0, 82, 31},
		{44, "Locked", "QAM256", 411000000, 4.3, 41.4, 68, 21},
		{45, "Locked", "QAM256", 417000000, 4.3, 41.5, 62, 16},
	}

	if len(expectedDownstreamChannels) != len(actualDownstreamChannels) {
		t.Errorf("mismatch in row count, expected %d, got %d", len(expectedDownstreamChannels), len(actualDownstreamChannels))
	}

	for rowNumber, actualRow := range actualDownstreamChannels {

		if rowNumber < len(expectedDownstreamChannels) {
			expectedRow := expectedDownstreamChannels[rowNumber]
			if expectedRow.ChannelId != actualRow.ChannelId {
				t.Errorf("unexpected ChannelId in row %d, expected %d, got %d", rowNumber, expectedRow.ChannelId, actualRow.ChannelId)
			}
			if expectedRow.LockStatus != actualRow.LockStatus {
				t.Errorf("unexpected Lockstatus in row %d, expected %s, got %s", rowNumber, expectedRow.LockStatus, actualRow.LockStatus)
			}
			if expectedRow.Modulation != actualRow.Modulation {
				t.Errorf("unexpected Lockstatus in row %d, expected %s, got %s", rowNumber, expectedRow.Modulation, actualRow.Modulation)
			}
			if expectedRow.Frequency != actualRow.Frequency {
				t.Errorf("unexpected Frequency in row %d, expected %d, got %d", rowNumber, expectedRow.Frequency, actualRow.Frequency)
			}
			if expectedRow.Power != actualRow.Power {
				t.Errorf("unexpected Power in row %d, expected %f, got %f", rowNumber, expectedRow.Power, actualRow.Power)
			}
			if expectedRow.SnrSmr != actualRow.SnrSmr {
				t.Errorf("unexpected SnrSmr in row %d, expected %f, got %f", rowNumber, expectedRow.SnrSmr, actualRow.SnrSmr)
			}
			if expectedRow.Corrected != actualRow.Corrected {
				t.Errorf("unexpected Corrected in row %d, expected %d, got %d", rowNumber, expectedRow.Corrected, actualRow.Corrected)
			}
			if expectedRow.Uncorrectables != actualRow.Uncorrectables {
				t.Errorf("unexpected Uncorrectables in row %d, expected %d, got %d", rowNumber, expectedRow.Uncorrectables, actualRow.Uncorrectables)
			}
		}
	}

	expectedUpstreamChannels := []UpstreamBoundedChannel{
		{1, 1, "Locked", "SC-QAM Upstream", 35600000, 6400000, 45.0},
		{2, 2, "Locked", "SC-QAM Upstream", 29200000, 6400000, 43.0},
		{3, 3, "Locked", "SC-QAM Upstream", 22800000, 6400000, 45.0},
		{4, 4, "Locked", "SC-QAM Upstream", 16400000, 6400000, 45.0},
		{5, 5, "Locked", "SC-QAM Upstream", 40400000, 3200000, 42.0},
	}
	if len(expectedUpstreamChannels) != len(actualUpstreamChannels) {
		t.Errorf("mismatch in row count, expected %d, got %d", len(expectedUpstreamChannels), len(actualUpstreamChannels))
	}
	for rowNumber, actualRow := range actualUpstreamChannels {

		if rowNumber < len(expectedUpstreamChannels) {
			expectedRow := expectedUpstreamChannels[rowNumber]
			if expectedRow.Channel != actualRow.Channel {
				t.Errorf("unexpected Channel in row %d, expected %d, got %d", rowNumber, expectedRow.Channel, actualRow.Channel)
			}
			if expectedRow.ChannelId != actualRow.ChannelId {
				t.Errorf("unexpected ChannelId in row %d, expected %d, got %d", rowNumber, expectedRow.ChannelId, actualRow.ChannelId)
			}
			if expectedRow.LockStatus != actualRow.LockStatus {
				t.Errorf("unexpected Lockstatus in row %d, expected %s, got %s", rowNumber, expectedRow.LockStatus, actualRow.LockStatus)
			}
			if expectedRow.UsChannelType != actualRow.UsChannelType {
				t.Errorf("unexpected US Channel Type in row %d, expected %s, got %s", rowNumber, expectedRow.UsChannelType, actualRow.UsChannelType)
			}
			if expectedRow.Frequency != actualRow.Frequency {
				t.Errorf("unexpected Frequency in row %d, expected %d, got %d", rowNumber, expectedRow.Frequency, actualRow.Frequency)
			}
			if expectedRow.Width != actualRow.Width {
				t.Errorf("unexpected Width in row %d, expected %d, got %d", rowNumber, expectedRow.Width, actualRow.Width)
			}
			if expectedRow.Power != actualRow.Power {
				t.Errorf("unexpected Power in row %d, expected %f, got %f", rowNumber, expectedRow.Power, actualRow.Power)
			}
		}
	}
}
