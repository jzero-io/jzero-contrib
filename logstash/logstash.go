package logstash

// Logstash interface defines the contract for log sending
type Logstash interface {
	SendLog(index, data string) error
}
