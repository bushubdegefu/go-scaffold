
package messages

//  request object def for the consumer
// only messages with this struct will be proccessed correctly
// when sent to specfic queue, ( morel like interface for communtication between publisher and consumer )
type RequestObject struct {
	Host     string
	Endpoint string
	Method   string
	Body     string
	Tp       map[string][]string
}
