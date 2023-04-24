package printer

var globalPrint = NewPrint(LevelDebug)

func Printf(p string, a ...any) {
	globalPrint.WriteToStdf(p, a...)
}

func Print(s string) {
	globalPrint.WriteToStd([]byte(s))
}

func PrintError(err error) {
	if err == nil {
		globalPrint.WriteToError([]byte("<nil>"))
	} else {
		globalPrint.WriteToError([]byte(err.Error()))
	}
}

func PrintErrorS(err string) {
	globalPrint.WriteToError([]byte(err))
}

func PrintErrorSf(err string, args ...any) {
	globalPrint.WriteToErrf(err, args...)
}

func Error(format string, a ...interface{}) {
	globalPrint.Error(format, a...)
}

func Warn(format string, a ...interface{}) {
	globalPrint.Warn(format, a...)
}

func Info(format string, a ...interface{}) {
	globalPrint.Info(format, a...)
}

func Debug(format string, a ...interface{}) {
	globalPrint.Debug(format, a...)
}
