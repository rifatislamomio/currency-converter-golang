/**
 * @author Rifat I.
 * @project: currency-converter
 * @date: 11-11-23
 */
package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"lambda-currency-converter-golang/handler"
)

func main() {
	lambda.Start(handler.Handler)
}
