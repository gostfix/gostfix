package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/gostfix/gostfix/internal/global"
	"github.com/gostfix/gostfix/internal/util"

	"github.com/alecthomas/kong"
)

func main() {
	// Initialize
	syscall.Umask(077)

	if os.Getenv(global.CONF_ENV_VERB) != "" {
		util.MsgVerbose = 1
	}

	// Don't die when a process goes away unexpectedly
	signal.Ignore(syscall.SIGPIPE)

	// Strip and save the process name for diagnostics etc.
	global.VarProcname = filepath.Base(os.Args[0])

	/*
	 * When running a child process, don't leak any open files that were
	 * leaked to us by our own (privileged) parent process. Descriptors 0-2
	 * are taken care of after we have initialized error logging.
	 *
	 * Some systems such as AIX have a huge per-process open file limit. In
	 * those cases, limit the search for potential file descriptor leaks to
	 * just the first couple hundred.
	 *
	 * The Debian post-installation script passes an open file descriptor into
	 * the master process and waits forever for someone to close it. Because
	 * of this we have to close descriptors > 2, and pray that doing so does
	 * not break things.
	 */
	// TODO(alf): postfix closes all file descriptors from [3, 500)
	// closefrom(3)

	// Initialize logging and exit handler.
	global.MailLogClientInit(global.MailTask(global.VarProcname),
		global.MAILLOG_CLIENT_FLAG_LOGWRITER_FALLBACK)

	/*
	 * The mail system must be run by the superuser so it can revoke
	 * privileges for selected operations. That's right - it takes privileges
	 * to toss privileges.
	 */
	// if os.Getuid() != 0 {
	//     log.Fatalf("the master command is reserved for the superuser")
	// }
	// if util.Unsafe() {
	//     log.Fatalf("the master command must not run as a set-uid process")
	// }

	var cli struct {
		ConfigDir string `kong:"short='c',placeholder='DIR',type='existingdir',env='CONF_ENV_PATH',help='configuration directory'"`
		Debug     bool   `kong:"short='D',negatable,env='CONF_ENV_DEBUG',help=''"`
		Detach    bool   `kong:"short='d',default='true',negatable,help='detach process from the shell'"`
		ExitAfter int    `kong:"short='e',placeholder='N',help='exit process after N seconds'"`
		Test      bool   `kong:"short='t',negatable,xor='Wait,Test',help=''"`
		Verbose   int    `kong:"short='v',placeholder='LEVEL',env='CONF_ENV_VERB',type='counter',help='enable more verbose logging'"`
		Wait      bool   `kong:"short='w',negatable,xor='Test,Wait',help=''"`
		Init      bool   `kong:"short='i',hidden,negatable,help=''"`
		Stdout    bool   `kong:"short='s',hidden,negatable,help=''"`
	}
	ctx := kong.Parse(&cli)
	// Sanity check
	if cli.Init && (cli.Debug || !cli.Detach || cli.Wait) {
		ctx.Fatalf("--init can't be used with --debug, --no-detach, or --wait")
	}

	fmt.Printf("cli: %v\n", cli)
	var err error = nil
	ctx.FatalIfErrorf(err)
}
