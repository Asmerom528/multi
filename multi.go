// Package multi provides concurrent groups for exploring native OS threading
// and multiprocessing in Go.
//
// It offers three implementations of concurrent groups with a similar API:
//   - goro.Group: runs functions in goroutines locked to OS threads.
//   - pthread.Group: runs functions in separate OS threads using POSIX threads.
//   - proc.Group: runs functions in separate OS processes.
//
// This is a research project. pthread and proc implementations bypass Go
// runtime safety and are not suitable for production use.
package multi
