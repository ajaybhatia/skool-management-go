package shared

import (
	"errors"
	"sync"
	"time"
)

// CircuitBreakerState represents the state of the circuit breaker
type CircuitBreakerState int

const (
	StateClosed CircuitBreakerState = iota
	StateHalfOpen
	StateOpen
)

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	name            string
	maxFailures     int
	resetTimeout    time.Duration
	failureCount    int
	lastFailureTime time.Time
	state           CircuitBreakerState
	mutex           sync.Mutex
}

// CircuitBreakerConfig holds configuration for circuit breaker
type CircuitBreakerConfig struct {
	Name         string
	MaxFailures  int
	ResetTimeout time.Duration
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(config CircuitBreakerConfig) *CircuitBreaker {
	return &CircuitBreaker{
		name:         config.Name,
		maxFailures:  config.MaxFailures,
		resetTimeout: config.ResetTimeout,
		state:        StateClosed,
	}
}

// Execute runs the given function with circuit breaker protection
func (cb *CircuitBreaker) Execute(fn func() error) error {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	// Check if circuit breaker should be reset
	if cb.state == StateOpen && time.Since(cb.lastFailureTime) > cb.resetTimeout {
		cb.state = StateHalfOpen
		cb.failureCount = 0
		LogInfo("CIRCUIT_BREAKER", cb.name+" circuit breaker moved to HALF_OPEN state")
	}

	// If circuit is open, return error immediately
	if cb.state == StateOpen {
		return errors.New("circuit breaker is OPEN")
	}

	// Execute the function
	err := fn()

	if err != nil {
		cb.onFailure()
		return err
	}

	cb.onSuccess()
	return nil
}

// onFailure handles failure cases
func (cb *CircuitBreaker) onFailure() {
	cb.failureCount++
	cb.lastFailureTime = time.Now()

	if cb.failureCount >= cb.maxFailures {
		cb.state = StateOpen
		LogError("CIRCUIT_BREAKER", cb.name+" circuit breaker OPENED",
			errors.New("max failures reached"))
	}
}

// onSuccess handles success cases
func (cb *CircuitBreaker) onSuccess() {
	if cb.state == StateHalfOpen {
		cb.state = StateClosed
		cb.failureCount = 0
		LogInfo("CIRCUIT_BREAKER", cb.name+" circuit breaker moved to CLOSED state")
	} else if cb.state == StateClosed {
		cb.failureCount = 0
	}
}

// GetState returns the current state of the circuit breaker
func (cb *CircuitBreaker) GetState() CircuitBreakerState {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	return cb.state
}

// GetFailureCount returns the current failure count
func (cb *CircuitBreaker) GetFailureCount() int {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	return cb.failureCount
}

// StateString returns string representation of the state
func (state CircuitBreakerState) String() string {
	switch state {
	case StateClosed:
		return "CLOSED"
	case StateHalfOpen:
		return "HALF_OPEN"
	case StateOpen:
		return "OPEN"
	default:
		return "UNKNOWN"
	}
}
