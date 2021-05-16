package main

import (
	"fmt"
	"time"
)

var CachedException = fmt.Errorf("CachedException")

type CircuitBreakerService func(timeout time.Duration, args ...interface{}) (interface{}, error)

type State int

const (
	CLOSED State = iota
	OPENED
	HALF_OPENED
)

var stringStates = []string{"CLOSED", "OPENED", "HALF_OPENED"}

type CircuitBreaker struct {
	service          CircuitBreakerService
	timeout          time.Duration
	retryTimeout     time.Duration
	failureThreshold int

	state               State
	lastFailureTime     time.Time
	lastFailureResponse error
	failureCount        int
}

func NewCircuitBreaker(service CircuitBreakerService,
	timeout time.Duration,
	retryTimeout time.Duration,
	failureThreshold int,
) *CircuitBreaker {
	return &CircuitBreaker{
		service:             service,
		timeout:             timeout,
		retryTimeout:        retryTimeout,
		failureThreshold:    failureThreshold,
		state:               CLOSED,
		lastFailureTime:     time.Time{},
		lastFailureResponse: nil,
		failureCount:        0,
	}
}

func (self *CircuitBreaker) AttemptRequest(args ...interface{}) (interface{}, error) {
	self.EvaluateState()
	if self.state == OPENED {
		return nil, CachedException
	}
	response, err := self.service(self.timeout, args...)
	if err != nil {
		self.responseFailure(err)
	} else {
		self.responseSuccess()
	}
	return response, err
}

func (self *CircuitBreaker) EvaluateState() {
	if self.failureCount >= self.failureThreshold {
		if time.Now().Sub(self.lastFailureTime) > self.retryTimeout {
			self.state = HALF_OPENED
			self.failureCount = self.failureCount / 2
		} else {
			self.state = OPENED
		}
	} else {
		self.state = CLOSED
	}
}

func (self *CircuitBreaker) responseSuccess() {
	self.state = CLOSED
	self.failureCount = 0
	self.lastFailureTime = time.Time{}
}

func (self *CircuitBreaker) responseFailure(err error) {
	self.state = OPENED
	self.failureCount += 1
	self.lastFailureTime = time.Now()
	self.lastFailureResponse = err
}

func (self *CircuitBreaker) GetState() State {
	return self.state
}

func (self *CircuitBreaker) GetStrState() string {
	return stringStates[self.state]
}

func (self *CircuitBreaker) SetState(state State) {
	self.state = state
	switch state {
	case OPENED:
		self.failureCount = self.failureThreshold
		self.lastFailureTime = time.Now()
	case CLOSED:
		self.failureCount = 0
		self.lastFailureTime = time.Time{}
	case HALF_OPENED:
		self.failureCount = self.failureThreshold / 2
		self.lastFailureTime = time.Now().Add(-self.retryTimeout)
	}
}

func testService(timeout time.Duration, args ...interface{}) (interface{}, error) {
	fmt.Println("testService timeout", timeout)
	time.Sleep(timeout)
	return "Success Response from service", nil
}

func testErrorService(timeout time.Duration, args ...interface{}) (interface{}, error) {
	fmt.Println("testErrorService timeout", timeout)
	time.Sleep(timeout)
	return "Error Response from service", fmt.Errorf("testErrorService Error")
}

var errorThresholdTime time.Time

func testNotStableService(timeout time.Duration, args ...interface{}) (interface{}, error) {
	if errorThresholdTime.After(time.Now()) {
		fmt.Println("testNotStableService timeout", timeout)
		time.Sleep(timeout)
		if errorThresholdTime.Add(time.Duration(4) * time.Second).Before(time.Now()) {
			errorThresholdTime = errorThresholdTime.Add(time.Duration(30) * time.Minute)
		}
		return "Error Response from service", fmt.Errorf("testNotStableService Error")
	}
	return "Success Response from service", nil
}

func main() {
	timeout := time.Duration(1) * time.Second
	retryTimeout := time.Duration(2) * time.Second
	failureThreshold := 4
	cb := NewCircuitBreaker(
		testService,
		timeout,
		retryTimeout,
		failureThreshold,
	)
	fmt.Println(cb.GetStrState())
	for i := 0; i < 5; i++ {
		r, err := cb.AttemptRequest("test request")
		fmt.Println(r, err, cb.GetStrState())
	}

	fmt.Println("ERROR SERVICE")
	cb = NewCircuitBreaker(
		testErrorService,
		timeout,
		retryTimeout,
		failureThreshold,
	)
	fmt.Println(cb.GetStrState())
	for i := 0; i < 20; i++ {
		r, err := cb.AttemptRequest("test request")
		fmt.Println(r, err, cb.GetStrState())
		if i == 10 {
			fmt.Println("Wait CircuitBreaker state change")
			time.Sleep(time.Duration(3) * time.Second)
		}
	}

	fmt.Println("NOT STABLE SERVICE")
	cb = NewCircuitBreaker(
		testNotStableService,
		timeout,
		retryTimeout,
		failureThreshold,
	)
	errorThresholdTime = time.Now()
	fmt.Println(cb.GetStrState())
	for i := 0; i < 30; i++ {
		r, err := cb.AttemptRequest("test request")
		fmt.Println(r, err, cb.GetStrState())

		switch i {
		case 5:
			errorThresholdTime = time.Now().Add(time.Duration(15) * time.Second)
		case 10:
			fmt.Println("Wait CircuitBreaker state change")
			time.Sleep(time.Duration(4) * time.Second)
		case 20:
			fmt.Println("Wait CircuitBreaker state change")
			time.Sleep(time.Duration(5) * time.Second)
		}
	}
}
