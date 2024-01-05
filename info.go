package gojourney

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func NewInfoCommandPayload(guildID, channelID, sessionID string) *Payload {
	return NewApplicationCommandPayload(
		guildID,
		channelID,
		sessionID,
		"info",
		"",
		"",
		ApplicationCommandIdVersion{
			Version: "1166847114203123799",
			ID:      "972289487818334209",
		})
}

type TimeRemaining struct {
	Remaining float64 `json:"remaining,omitempty"`
	Total     float64 `json:"total,omitempty"`
}

type Usage struct {
	Images int64   `json:"images,omitempty"`
	Hours  float64 `json:"hours,omitempty"`
}

type Subscription struct {
	Plan                string `json:"plan,omitempty"`
	PlanStatus          string `json:"plan_status,omitempty"`
	RenewsNextTimestamp int64  `json:"renews_next_timestamp,omitempty"`
}

type UserInfo struct {
	UserID            string         `json:"user_id,omitempty"`
	Subscription      *Subscription  `json:"subscription,omitempty"`
	VisibilityMode    string         `json:"visibility_mode,omitempty"`
	FastTimeRemaining *TimeRemaining `json:"fast_time_remaining,omitempty"`
	LifetimeUsage     *Usage         `json:"lifetime_usage,omitempty"`
	RelaxedUsage      *Usage         `json:"relaxed_usage,omitempty"`
	QueuedJobsFast    int            `json:"queued_jobs_fast,omitempty"`
	QueuedJobsRelax   int            `json:"queued_jobs_relax,omitempty"`
}

func (c *Client) Info(ctx context.Context) (*UserInfo, error) {

	payload := NewInfoCommandPayload(c.GuildID, c.ChannelID, c.ws.SessionID())
	messageFinder := &MessageFinder{Nonce: payload.Nonce, NoJob: true}

	ob := NewCommonMessageObserver(messageFinder.FilterMessage)
	c.RegisterMessageObserver(ob)
	defer c.UnregisterMessageObserver(ob)

	err := c.sendHttpInteraction(ctx, payload)
	if err != nil {
		return nil, err
	}

	msg, err := ob.WaitMsg(time.Minute * 5)
	if err != nil {
		return nil, err
	}

	if len(msg.Embeds) == 0 {
		return nil, ErrorEmptyEmbeds
	}

	return getUserInfo(msg.Embeds[0].Description), nil
}

func getUserInfo(info string) *UserInfo {
	var result UserInfo
	lines := strings.Split(info, "\n")
	for _, l := range lines {
		index := strings.Index(l, ":")
		if index == -1 {
			continue
		}

		k := strings.Trim(l[:index], "*")
		v := strings.Trim(l[index+1:], " ")
		switch k {
		case "User ID":
			result.UserID = v
		case "Visibility Mode":
			result.VisibilityMode = v
		case "Subscription":
			result.Subscription = ParseSubscription(v)
		case "Fast Time Remaining":
			result.FastTimeRemaining = ParseTimeRemaining(v)
		case "Lifetime Usage":
			result.LifetimeUsage = ParseUsage(v)
		case "Relaxed Usage":
			result.RelaxedUsage = ParseUsage(v)
		case "Queued Jobs (fast)":
			jobs, err := strconv.ParseInt(v, 10, 64)
			if err == nil {
				result.QueuedJobsFast = int(jobs)
			}
		case "Queued Jobs (relax)":
			jobs, err := strconv.ParseInt(v, 10, 64)
			if err == nil {
				result.QueuedJobsRelax = int(jobs)
			}
		}
	}

	return &result
}

func ParseTimeRemaining(s string) *TimeRemaining {
	re := regexp.MustCompile(`(\d+\.\d+)\/(\d+\.\d+)`)
	result := re.FindStringSubmatch(s)
	if len(result) == 3 {
		remaining, _ := strconv.ParseFloat(result[1], 64)
		total, _ := strconv.ParseFloat(result[2], 64)
		return &TimeRemaining{
			Remaining: remaining,
			Total:     total,
		}
	}
	return nil
}

func ParseUsage(str string) *Usage {
	re := regexp.MustCompile(`(\d+) images \(([\d.]+) hours\)`)
	result := re.FindStringSubmatch(str)
	if len(result) == 3 {
		images, _ := strconv.ParseInt(result[1], 10, 64)
		hours, _ := strconv.ParseFloat(result[2], 64)
		return &Usage{
			Images: images,
			Hours:  hours,
		}
	}
	return nil
}

func ParseSubscription(str string) *Subscription {
	re := regexp.MustCompile(`^(\w+) \(([^,]+),.*<t:(\d+)>`)
	submatch := re.FindStringSubmatch(str)

	if len(submatch) == 4 {
		result := &Subscription{}
		result.Plan = submatch[1]
		result.PlanStatus = submatch[2]
		result.RenewsNextTimestamp, _ = strconv.ParseInt(submatch[3], 10, 64)
		return result
	}

	return nil
}
