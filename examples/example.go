// Package examples
// Created by zc on 2021/9/18.
package main

import (
	"context"

	"github.com/zc2638/aide"
)

func main() {
	first := aide.NewStage("first").AddSteps(
		check().Step("check"),
		check().Step("check"),
		check().Step("check"),
		check().Step("check"),
	)
	second := aide.NewStage("second").AddSteps(
		install().Step("install"),
	)
	third := aide.NewStage("third").AddSteps(
		health().Step("health check"),
	)

	a := aide.New()
	a.AddStages(first, second, third)
	_ = a.Run(context.Background())
}

func check() aide.StepFunc {
	return func(ctx *aide.StepContext) {
		ctx.WriteString("check Port 31181 OK.")
	}
}

func install() aide.StepFunc {
	return func(ctx *aide.StepContext) {
		ctx.WriteString("Install Component Successful.")
	}
}

func health() aide.StepFunc {
	return func(ctx *aide.StepContext) {
		ctx.Exit(1)
		ctx.WriteString("component unhealthy.")
	}
}
