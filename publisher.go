package killcord

import (
	"errors"
	"fmt"
	"log"
	"time"
)

const (
	defaultWarningThreshold    = 86400  // 24 hours
	defaultPublishThreshold    = 172800 // 48 hours
	defaultPublisherMaxRetries = 5
)

func (s *Session) PublishKey() error {
	switch s.Config.Type {
	case "owner":
		if err := SetKey(s.Config.Contract.Owner, s.Config.Contract.ID, s.Config.Payload.Secret); err != nil {
			return err
		}
	case "publisher":
		if err := SetKey(s.Config.Contract.Publisher, s.Config.Contract.ID, s.Config.Payload.Secret); err != nil {
			return err
		}
	default:
		return fmt.Errorf("project type must be `owner` or `publisher`, received: %v\n", s.Config.Type)
	}
	return nil
}

func (s *Session) RunPublisher() error {
	s.configurePublisherSession()
	if err := s.validatePublisherSession(); err != nil {
		return err
	}
	id := s.Config.Contract.ID
	key, err := GetKey(id)
	if err != nil {
		return err
	}
	if key != "" {
		return errors.New("secret key already published, skipping")
	}
	checkin, err := getLastCheckinWithRetries(defaultPublisherMaxRetries, id)
	if err != nil {
		return err
	}
	fmt.Printf("Last Checkin:\t\t%s\n", checkin)
	if publishThresholdReached(checkin, s.Config.Publisher.PublishThreshold) {
		fmt.Println("publish threshold breached, publishing secret key")
		if err := s.PublishKey(); err != nil {
			return err
		}
		return nil
	}
	if warningThresholdReached(checkin, s.Config.Publisher.WarningThreshold) {
		// log warning, add additional features later
		fmt.Println("warning: checkin has broken the warning threshold")
		return nil
	}
	return nil
}

// Configure defaults and waterfall overrides from options and Environmental variables.
// config -> defaults -> options -> ENV

func (s *Session) configurePublisherSession() {
	if s.Config.Publisher.WarningThreshold == 0 {
		s.Config.Publisher.WarningThreshold = defaultWarningThreshold
	}
	if s.Config.Publisher.PublishThreshold == 0 {
		s.Config.Publisher.PublishThreshold = defaultPublishThreshold
	}
	if s.Config.Type == "" {
		s.Config.Type = "publisher"
	}
	// check opts and override for session if values exist
	if s.Options.Publisher.WarningThreshold != 0 {
		s.Config.Publisher.WarningThreshold = s.Options.Publisher.WarningThreshold
	}
	if s.Options.Publisher.PublishThreshold != 0 {
		s.Config.Publisher.PublishThreshold = s.Options.Publisher.PublishThreshold
	}
	if s.Options.Contract.ID != "" {
		s.Config.Contract.ID = s.Options.Contract.ID
	}
	if s.Options.Publisher.Address != "" {
		s.Config.Contract.Publisher.Address = s.Options.Publisher.Address
	}
	if s.Options.Publisher.Password != "" {
		s.Config.Contract.Publisher.Password = s.Options.Publisher.Password
	}
	if s.Options.Publisher.KeyStore != "" {
		s.Config.Contract.Publisher.KeyStore = s.Options.Publisher.KeyStore
	}
	if s.Options.Payload.Secret != "" {
		s.Config.Payload.Secret = s.Options.Payload.Secret
	}
}

func (s *Session) validatePublisherSession() error {
	if s.Config.Payload.Secret == "" {
		return errors.New("payload secret not configured, exiting")
	}
	if s.Config.Contract.ID == "" {
		return errors.New("contract id not configured, exiting")
	}
	if s.Config.Contract.Publisher.Address == "" {
		return errors.New("publisher address not configured, exiting")
	}
	if s.Config.Contract.Publisher.Password == "" {
		return errors.New("publisher password not configured, exiting")
	}

	if s.Config.Contract.Publisher.KeyStore == "" {
		return errors.New("publisher keystore not configured, exiting")
	}

	if s.Config.Publisher.WarningThreshold == 0 {
		return errors.New("warning threshold not configured, exiting")
	}

	if s.Config.Publisher.PublishThreshold == 0 {
		return errors.New("publish threshold not configured, exiting")
	}
	return nil
}

func warningThresholdReached(checkin time.Time, warningThreshold int64) bool {
	threshold := time.Now().Unix() - warningThreshold
	fmt.Printf("Warning Threshold:\t%s\n", time.Unix(threshold, 0))
	if checkin.Unix() <= threshold {
		return true
	}
	return false
}

func publishThresholdReached(checkin time.Time, publishThreshold int64) bool {
	threshold := time.Now().Unix() - publishThreshold
	fmt.Printf("Publish Threshold:\t%s\n", time.Unix(threshold, 0))
	if checkin.Unix() <= threshold {
		return true
	}
	return false
}

func getLastCheckinWithRetries(maxRetries int, id string) (time.Time, error) {
	for retry := 1; retry <= maxRetries; retry++ {
		checkin, err := GetLastCheckIn(id)
		if err != nil {
			fmt.Println("checkin failed, retrying: ", err)
			time.Sleep(1 * time.Second)
			continue
		}
		return checkin, nil
	}
	return time.Time{}, fmt.Errorf("checkin max retry of %v reached, skipping\n", maxRetries)
}

func (s *Session) PublisherLambdaHandler() {
	if err := s.RunPublisher(); err != nil {
		log.Println(err)
	}
}
