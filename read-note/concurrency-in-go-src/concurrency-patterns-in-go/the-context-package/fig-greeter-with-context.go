package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background()) // <1>
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := printGreeting1(ctx); err != nil {
			fmt.Printf("cannot print greeting: %v\n", err)
			cancel() // <2>
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printFarewell1(ctx); err != nil {
			fmt.Printf("cannot print farewell: %v\n", err)
		}
	}()

	wg.Wait()
}

func printGreeting1(ctx context.Context) error {
	greeting, err := genGreeting1(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("%s world!\n", greeting)
	return nil
}

func printFarewell1(ctx context.Context) error {
	farewell, err := genFarewell1(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("%s world!\n", farewell)
	return nil
}

func genGreeting1(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second) // <3>
	defer cancel()

	switch locale, err := locale1(ctx); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "hello", nil
	}
	return "", fmt.Errorf("unsupported locale")
}

func genFarewell1(ctx context.Context) (string, error) {
	switch locale, err := locale1(ctx); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "goodbye", nil
	}
	return "", fmt.Errorf("unsupported locale")
}

func locale1(ctx context.Context) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err() // <4>
	case <-time.After(1 * time.Minute):
	}
	return "EN/US", nil
}
