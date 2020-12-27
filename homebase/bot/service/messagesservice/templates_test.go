package messagesservice

import (
	"strings"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
)

func TestTripCompletedMessage(t *testing.T) {
	clock := clockwork.NewFakeClock()
	cases := []struct {
		name   string
		input  tripCompletedTemplateInput
		output string
	}{
		{
			name: "recent trip with no tags",
			input: tripCompletedTemplateInput{
				ReturnedAt: clock.Now().Add(-3 * time.Minute),
				Duration:   (12 * time.Minute) + (15 * time.Second),
				Tags:       []string{},
			},
			output: `You just returned from a trip that lasted *12 minutes*\.`,
		},
		{
			name: "old trip with no tags",
			input: tripCompletedTemplateInput{
				ReturnedAt: clock.Now().Add(-10 * time.Minute),
				Duration:   3 * time.Hour,
			},
			output: `You returned 10 minutes ago from a trip that lasted *3 hours*\.`,
		},
		{
			name: "trip with tags",
			input: tripCompletedTemplateInput{
				ReturnedAt: clock.Now().Add(-3 * time.Minute),
				Duration:   (12 * time.Minute) + (15 * time.Second),
				Tags:       []string{"dog walk", "cold weather"},
			},
			output: `You just returned from a trip that lasted *12 minutes*\.

🏷 dog walk, cold weather`,
		},
	}

	temps, err := parseTemplates(clock)
	if !assert.NoError(t, err) {
		return
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var output strings.Builder
			if assert.NoError(t, temps.ExecuteTemplate(&output, tripCompletedTemplate, &c.input)) {
				assert.Equal(t, c.output, output.String())
			}
		})
	}
}
