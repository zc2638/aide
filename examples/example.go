// Package examples
// Created by zc on 2021/9/18.
package main

import (
	"context"

	"github.com/zc2638/aide"
)

func main() {
	first := aide.NewStage("first").AddSteps(
		aide.StepFunc(check).Step("check"),
		aide.StepFunc(check).Step("check"),
		aide.StepFunc(check).Step("check"),
	).AddStepFuncs(check)
	second := aide.NewStage("second").AddSteps(
		install().Step("install"),
	)
	third := aide.NewStage("third").AddSteps(
		tip().Step("tip"),
		health().Step("health check"),
		unreachable().Step("unreachable"),
	)
	fourth := aide.NewStage("fourth").AddStepFuncs(
		unreachableStage,
	)

	a := aide.New()
	a.AddStages(first, second, third, fourth)
	_ = a.Run(context.Background())
}

func check(sc *aide.StepContext) {
	sc.WriteString("check Port 31181 OK.")
}

func install() aide.StepFunc {
	return func(sc *aide.StepContext) {
		sc.WriteString("Install Component Successful.")
	}
}

func tip() aide.StepFunc {
	return func(sc *aide.StepContext) {
		sc.WithLevel(aide.WarnLevel).WriteString("There is an exception.")
	}
}

func health() aide.StepFunc {
	return func(sc *aide.StepContext) {
		sc.Exit(1).WriteString("component unhealthy.")
	}
}

func unreachable() aide.StepFunc {
	return func(sc *aide.StepContext) {
		sc.WriteString("unreachable.")
	}
}

func unreachableStage(sc *aide.StepContext) {
	sc.WriteString("unreachable stage.")
}
