package middleware

import (
	"fmt"

	"github.com/perlou/go-gateway-demo/proxy/public"
)

func FlowCountMiddleWare(counter *public.FlowCountService) func(c *SliceRouterContext) {
	return func(c *SliceRouterContext) {
		counter.Increase()
		fmt.Println("QPS:", counter.QPS)
		fmt.Println("TotalCount:", counter.TotalCount)
		c.Next()
	}
}
