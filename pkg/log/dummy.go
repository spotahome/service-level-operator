package log

// Dummy is a dummy logger
var Dummy = dummyLogger{}

type dummyLogger struct{}

func (l dummyLogger) Debug(...interface{})                           {}
func (l dummyLogger) Debugln(...interface{})                         {}
func (l dummyLogger) Debugf(string, ...interface{})                  {}
func (l dummyLogger) Info(...interface{})                            {}
func (l dummyLogger) Infoln(...interface{})                          {}
func (l dummyLogger) Infof(string, ...interface{})                   {}
func (l dummyLogger) Warn(...interface{})                            {}
func (l dummyLogger) Warnln(...interface{})                          {}
func (l dummyLogger) Warnf(string, ...interface{})                   {}
func (l dummyLogger) Warningf(format string, args ...interface{})    {}
func (l dummyLogger) Error(...interface{})                           {}
func (l dummyLogger) Errorln(...interface{})                         {}
func (l dummyLogger) Errorf(string, ...interface{})                  {}
func (l dummyLogger) Fatal(...interface{})                           {}
func (l dummyLogger) Fatalln(...interface{})                         {}
func (l dummyLogger) Fatalf(string, ...interface{})                  {}
func (l dummyLogger) Panic(...interface{})                           {}
func (l dummyLogger) Panicln(...interface{})                         {}
func (l dummyLogger) Panicf(string, ...interface{})                  {}
func (l dummyLogger) With(key string, value interface{}) Logger      { return l }
func (l dummyLogger) WithField(key string, value interface{}) Logger { return l }
func (l dummyLogger) Set(level Level) error                          { return nil }
