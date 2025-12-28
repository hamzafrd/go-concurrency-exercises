//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"sync"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
	mu        sync.Mutex
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	u.mu.Lock()
	// 1. Check if user is allowed to proceed
	// If NOT premium AND time is already used up, reject immediately
	if !u.IsPremium && u.TimeUsed >= 10 {
		u.mu.Unlock()
		return false
	}
	u.mu.Unlock()

	// 2. Execute the process
	startTime := time.Now()
	process()
	elapsedTime := time.Since(startTime).Seconds()

	// 3. Update the state
	u.mu.Lock()
	defer u.mu.Unlock()

	u.TimeUsed += int64(elapsedTime)

	// 4. Return whether they are still within limits after this run
	// or if they are premium
	if u.IsPremium {
		return true
	}

	return u.TimeUsed < 10
}

func main() {
	RunMockServer()
}
