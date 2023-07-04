package logic

import (
	"testing"

	"github.com/stretchr/testify/require"

	"git.lothric.net/examples/go/gogin/internal/pkg/logger"
)

func TestGistsLogicCreation(t *testing.T) {
	for scenario, fn := range map[string]func(
		t *testing.T,
		l logger.Log,
		m *reporterMock,
	){
		"fails to create if no logger provided":   testFailsIfNoLogger,
		"fails to create if no reporter provided": testFailsIfNoReporter,
	} {
		t.Run(scenario, func(t *testing.T) {
			log, _ := logger.NewNullLogger()
			reporter := &reporterMock{}

			fn(t, log, reporter)
		})
	}
}

type reporterMock struct {
}

func testFailsIfNoLogger(
	t *testing.T,
	l logger.Log,
	m *reporterMock,
) {

	gists, err := NewGistsLogic(
		l,
		nil,
	)

	require.Nil(t, gists)
	require.Equal(t, ErrNoReporterProvided, err)
}

func testFailsIfNoReporter(
	t *testing.T,
	l logger.Log,
	m *reporterMock,
) {

	gists, err := NewGistsLogic(
		nil,
		m,
	)

	require.Nil(t, gists)
	require.Equal(t, ErrNoLoggerProvided, err)
}
