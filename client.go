package gojourney

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/coocood/freecache"

	"github.com/kiwioneone/gojourney/rest"

	"github.com/kiwioneone/gojourney/discord"

	"github.com/kiwioneone/gojourney/gateway"
	"github.com/kiwioneone/gojourney/gateway/event"
)

type Client struct {
	UserToken            string
	GuildID              string
	ChannelID            string
	ws                   *gateway.Session
	login                bool
	observerLock         sync.RWMutex
	interactionObservers map[InteractionObserver]struct{}
	messageObservers     map[MessageObserver]struct{}
	assistBot            *rest.Client
	cache                *freecache.Cache
}

func NewClient(userToken string, guildID string, channelID string, options ...Option) (*Client, error) {

	ws := gateway.NewSession(&gateway.Options{
		Token:   userToken,
		Intents: gateway.IntentsNonPrivileged,
	})

	client := &Client{
		UserToken:            userToken,
		GuildID:              guildID,
		ChannelID:            channelID,
		ws:                   ws,
		interactionObservers: make(map[InteractionObserver]struct{}),
		messageObservers:     make(map[MessageObserver]struct{}),
	}

	err := ws.On(event.EventReady, func() {
		fmt.Println("Logged in as " + ws.Me().Tag())
		client.login = true
	})

	if err != nil {
		return nil, err
	}

	err = ws.On(event.EventInteractionCreate, client.OnInteractionCreate)

	if err != nil {
		return nil, err
	}

	err = ws.On(event.EventMessageCreate, client.OnMessageCreate)
	if err != nil {
		return nil, err
	}

	err = ws.On(event.EventMessageUpdate, client.OnMessageCreate)
	if err != nil {
		return nil, err
	}

	err = ws.On(event.EventMessageDelete, client.OnMessageDelete)
	if err != nil {
		return nil, err
	}

	err = ws.Login()
	if err != nil {
		return nil, err
	}

	client.waitLogin()
	if !client.login {
		client.Close()
		if ws.Error != nil {
			return nil, ws.Error
		}
		return nil, ErrorNotLogin
	}

	for _, opt := range options {
		opt(client)
	}

	client.cache = freecache.NewCache(10 * 1024 * 1024)

	return client, nil
}

func (c *Client) Close() {
	if c.ws != nil {
		c.ws.Close()
	}
}

type CommandResult struct {
	MessageID string
	ImageURL  string
	JobID     string
}

func (c *Client) waitLogin() {
	start := time.Now()
	for !c.login {
		if c.ws.Error != nil {
			return
		}
		time.Sleep(time.Millisecond * 100)
		if time.Since(start) > time.Second*10 {
			break
		}
	}
}

func (c *Client) sendHttpInteraction(ctx context.Context, payload *Payload) error {

	payloadData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, DiscordInteractionURL, bytes.NewReader(payloadData))
	if err != nil {
		return err
	}

	req.Header.Set("authorization", c.UserToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("status %s\n response %s", resp.Status, string(responseBody))
	}

	return nil
}

func (c *Client) RegisterInteractionObserver(ob InteractionObserver) {
	c.observerLock.Lock()
	defer c.observerLock.Unlock()
	if c.interactionObservers == nil {
		c.interactionObservers = make(map[InteractionObserver]struct{})
	}
	c.interactionObservers[ob] = struct{}{}
}

func (c *Client) UnregisterInteractionObserver(ob InteractionObserver) {
	c.observerLock.Lock()
	defer c.observerLock.Unlock()
	if c.interactionObservers == nil {
		return
	}
	delete(c.interactionObservers, ob)
}

func (c *Client) RegisterMessageObserver(ob MessageObserver) {
	c.observerLock.Lock()
	defer c.observerLock.Unlock()
	if c.messageObservers == nil {
		c.messageObservers = make(map[MessageObserver]struct{})
	}
	c.messageObservers[ob] = struct{}{}
}

func (c *Client) UnregisterMessageObserver(ob MessageObserver) {
	c.observerLock.Lock()
	defer c.observerLock.Unlock()
	if c.messageObservers == nil {
		return
	}
	delete(c.messageObservers, ob)
}

func (c *Client) OnInteractionCreate(interaction *discord.Interaction) {
	data, _ := json.Marshal(interaction)
	fmt.Printf("interaction %s\n", string(data))
	for k := range c.interactionObservers {
		k.Observe(interaction)
	}
}

func (c *Client) OnMessageCreate(msg *discord.Message) {
	data, _ := json.Marshal(msg)
	fmt.Printf("msg %s\n", string(data))

	c.AppealBlockAction(msg)
	//err := m.messageStorage.Set([]byte(MessageID(msg)), data)
	//if err != nil {
	//	fmt.Printf("set %s failed, error %s", MessageID(msg), err.Error())
	//}

	for k := range c.messageObservers {
		k.Observe(msg)
	}
}

func (c *Client) OnMessageDelete(msg *event.MessageDeleteData) {
	data, _ := json.Marshal(msg)
	fmt.Printf("msg delete %s\n", string(data))

	c.cache.Set([]byte(msg.Id), data, int(time.Hour.Seconds()))
}

func (c *Client) SessionID() string {
	if c.ws == nil {
		return ""
	}
	return c.ws.SessionID()
}

func (c *Client) MessageDeleted(msgId string) bool {
	got, err := c.cache.Get([]byte(msgId))
	if err != nil {
		return false
	}
	return len(got) > 0
}
