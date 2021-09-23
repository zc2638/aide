// Package examples
// Created by zc on 2021/9/18.
package main

import (
	"context"
	"fmt"

	"github.com/zc2638/aide"
)

func main() {
	first := aide.NewStage("first").
		AddSteps(
			aide.StepFunc(check).Step("check1"),
			aide.StepFunc(check).Step("check2"),
			aide.StepFunc(check).Step("check3"),
		).
		AddStepFuncs(check)
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
	sc.Return("check Port 31181 OK.")
}

func install() aide.StepFunc {
	return func(sc *aide.StepContext) {
		sc.Log("ceshi log")
		sc.Return("Install Component Successful.")
	}
}

func tip() aide.StepFunc {
	return func(sc *aide.StepContext) {
		sc.Context().WithValue("test", "context test")
		sc.Return("There is an exception.")
	}
}

func health() aide.StepFunc {
	return func(sc *aide.StepContext) {
		fmt.Println(sc.Context().Value("test"))
		sc.Break("component unhealthy.")
	}
}

func unreachable() aide.StepFunc {
	return func(sc *aide.StepContext) {
		sc.Message("unreachable.")
	}
}

func unreachableStage(sc *aide.StepContext) {
	sc.Message("unreachable stage.")
}
