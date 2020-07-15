package fuse

import (
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
)

func fuse() {
	hystrix.ConfigureCommand("my_command", hystrix.CommandConfig{
		Timeout:               1000,
		MaxConcurrentRequests: 100,
		ErrorPercentThreshold: 25,
	})

	// async without waiting result
	hystrix.Go("my_command", func() error {
		// talk to other services
		fmt.Println("123")
		return nil
	}, func(err error) error {
		// do this when services are down
		fmt.Printf("err = %s", err.Error())
		return nil
	})

	// async with waiting result
	output := make(chan bool, 1)
	errors := hystrix.Go("my_command", func() error {
		// talk to other services
		output <- true
		return nil
	}, nil)

	select {
	case out := <-output:
		fmt.Printf("out == %t", out)
		// success
	case err := <-errors:
		fmt.Printf("err = %s", err.Error())
		// failure
	}

	// sync
	if err := hystrix.Do("my_command", func() error {
		// talk to other services
		return nil
	}, nil); err != nil {
		fmt.Printf("err = %s", err.Error())
	}
}
