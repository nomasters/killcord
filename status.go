//go:generate ./scripts/contract-gen.sh
//go:generate ./scripts/build.sh

package killcord

import (
	"fmt"
)

func (s *Session) GetStatus() error {
	if err := s.getProjectStatus(); err != nil {
		return err
	}
	if err := s.getContractStatus(); err != nil {
		return err
	}
	if err := s.getPayloadStatus(); err != nil {
		return err
	}
	if err := s.getPublisherStatus(); err != nil {
		return err
	}
	return nil
}

func (s *Session) getProjectStatus() error {
	fmt.Println("")
	fmt.Printf("Project Status:\t\t\t%v\n", s.Config.Status)
	fmt.Println("--------------")
	return nil
}

func (s *Session) getContractStatus() error {
	fmt.Println("")
	fmt.Printf("Contract Status:\t\t%v\n", s.Config.Contract.Status)
	fmt.Println("---------------")
	var registeredOwnerAddress, registeredPublisherAddress string
	if s.Config.Contract.ID != "" {
		fmt.Printf("\tContract ID:\t\t0x%v\n", s.Config.Contract.ID)

		// get last checkin
		checkin, err := GetLastCheckIn(s.Config.Contract.ID)
		if err != nil {
			return err
		}
		fmt.Printf("\tLast Checkin:\t\t%v\n", checkin)

		// get endpoint info
		endpoint, err := GetPayloadEndpoint(s.Config.Contract.ID)
		if err != nil {
			return err
		}
		if endpoint == "" {
			endpoint = "not configured"
		}
		fmt.Printf("\tRegistered Endpoint:\t%v\n", endpoint)

		// get Key info
		key, err := GetKey(s.Config.Contract.ID)
		if err != nil {
			return err
		}
		// if key is written to the contract and it does not
		// exist in the config, write it to the config
		if key != "" && s.Config.Payload.Secret == "" {
			s.Config.Payload.Secret = key
			fmt.Println("payload secret discovered, writing to config")
			fmt.Println("to decrypt the payload run: killcord decrypt")
		}
		if key == "" {
			key = "not published"
		}
		fmt.Printf("\tDecryption Key:\t\t%v\n", key)

		registeredOwnerAddress, err = GetOwner(s.Config.Contract.ID)
		if err != nil {
			return err
		}
		registeredPublisherAddress, err = GetPublisher(s.Config.Contract.ID)
		if err != nil {
			return err
		}
	}

	if s.Config.Contract.Owner.Address != "" {
		OwnBal := getBalance(s.Config.Contract.Owner.Address)
		fmt.Printf("\tOwner Balance:\t\t%v\n", OwnBal)
	}
	if s.Config.Contract.Owner.Address != "" {
		fmt.Printf("\tOwner (local):\t\t%v\n", "0x"+s.Config.Contract.Owner.Address)
	}
	if registeredOwnerAddress != "" {
		fmt.Printf("\tOwner (contract):\t%v\n", registeredOwnerAddress)
	}
	if s.Config.Contract.Publisher.Address != "" {
		pubBal := getBalance(s.Config.Contract.Publisher.Address)
		fmt.Printf("\tPublisher Balance:\t%v\n", pubBal)
	}
	if s.Config.Contract.Publisher.Address != "" {
		fmt.Printf("\tPublisher (local):\t%v\n", "0x"+s.Config.Contract.Publisher.Address)
	}
	if registeredPublisherAddress != "" {
		fmt.Printf("\tPublisher (contract):\t%v\n", registeredPublisherAddress)
	}
	return nil
}

func (s *Session) getPayloadStatus() error {
	fmt.Println("")
	fmt.Printf("Payload Status:\t\t\t%v\n", s.Config.Payload.Status)
	fmt.Println("--------------")
	return nil
}

func (s *Session) getPublisherStatus() error {
	fmt.Println("")
	fmt.Printf("Publisher Status:\t\t%v\n", s.Config.Publisher.Status)
	fmt.Println("----------------")
	return nil
}
