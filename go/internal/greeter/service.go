package greeter

import (
	"context"
	"fmt"

	helloworldv1 "github.com/ahmedalhulaibi/hello-world-api/internal/gen/helloworld/v1"
)

type greeter struct {
	greetingFormat string
	helloworldv1.UnimplementedGreeterServiceServer
}

func NewGreeter(greetingFormat string) *greeter {
	return &greeter{
		greetingFormat: greetingFormat,
	}
}

var _ helloworldv1.GreeterServiceServer = (*greeter)(nil)

func (g *greeter) SayHello(ctx context.Context, r *helloworldv1.SayHelloRequest) (*helloworldv1.SayHelloResponse, error) {
	return &helloworldv1.SayHelloResponse{
		Message: fmt.Sprintf(g.greetingFormat, r.Name),
	}, nil
}
